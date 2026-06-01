package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMatchDaySelectableMatchesIncludesUpcomingOutsideActiveSeason(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	start := time.Now().Add(2 * time.Hour)
	match := Match{TeamID: 7, Opponent: "Upcoming", StartTime: &start, SeasonId: 999}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}

	matches, err := matchDaySelectableMatchesForTeam(db, 7)
	if err != nil {
		t.Fatalf("select matches: %v", err)
	}
	if len(matches) != 1 || matches[0].ID != match.ID {
		t.Fatalf("expected upcoming match outside active season, got %#v", matches)
	}
}

func TestSaveMatchGeneratesStableMatchURLSlug(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}

	matchID, err := SaveMatch(db, &Match{Opponent: "Rovers"})
	if err != nil {
		t.Fatalf("save match: %v", err)
	}
	match, err := GetMatch(db, matchID)
	if err != nil {
		t.Fatalf("get match: %v", err)
	}
	if match.MatchURLSlug == "" {
		t.Fatal("expected match URL slug")
	}
	originalSlug := match.MatchURLSlug
	match.Location = "Home"
	if _, err := SaveMatch(db, match); err != nil {
		t.Fatalf("save match again: %v", err)
	}
	updated, err := GetMatch(db, matchID)
	if err != nil {
		t.Fatalf("get updated match: %v", err)
	}
	if updated.MatchURLSlug != originalSlug {
		t.Fatalf("expected stable slug %q, got %q", originalSlug, updated.MatchURLSlug)
	}
}

func TestPublicMatchURLRendersReadOnlyMatchDayAndFeedback(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}, &MLNote{}, &PlayerMatchUnavailability{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team", EnableLineupsModule: true, EnablePublicFeedbackForm: true}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Formation", Status: formationStatusLive}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	if err := db.Create(&FormationPosition{FormationID: formation.ID, IndexNumber: 1, PositionName: "ST", X: 50, Y: 50}).Error; err != nil {
		t.Fatalf("create position: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, FormationID: formation.ID, Status: lineupStatusSelected, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}
	player := Player{Name: "Sam", Active: true}
	if err := db.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}
	if err := db.Create(&LineupPlayer{LineupID: lineup.ID, IndexNumber: 1, SlotOrder: 0, PlayerID: player.ID}).Error; err != nil {
		t.Fatalf("create lineup player: %v", err)
	}
	start := time.Now().Add(2 * time.Hour)
	match := Match{TeamID: team.ID, LineupID: lineup.ID, Opponent: "Rovers", Location: "Home", StartTime: &start}
	if _, err := SaveMatch(db, &match); err != nil {
		t.Fatalf("save match: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-url/{matchURLSlug}", publicMatchURLHandler(db))
	router.HandleFunc("/match-url/{matchURLSlug}/feedback", publicMatchFeedbackHandler(db))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/match-url/"+match.MatchURLSlug, nil))
	body := response.Body.String()
	if response.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d: %s", response.Code, body)
	}
	if !strings.Contains(body, "Line-up") || !strings.Contains(body, "Rovers") {
		t.Fatalf("expected public match day content, got %s", body)
	}
	if strings.Contains(body, "Start match now") || strings.Contains(body, "Edit line-up") {
		t.Fatalf("expected read-only supporter view, got %s", body)
	}
	if !strings.Contains(body, "/match-url/"+match.MatchURLSlug+"/feedback") {
		t.Fatalf("expected match-scoped feedback link, got %s", body)
	}

	form := url.Values{}
	form.Set("scope", noteTargetMatch)
	form.Set("matchId", S(match.ID))
	form.Set("type", "specific")
	form.Set("creator", "Sideline")
	form.Set("note", "Great pressure")
	feedbackResponse := httptest.NewRecorder()
	feedbackRequest := httptest.NewRequest(http.MethodPost, "/match-url/"+match.MatchURLSlug+"/feedback", strings.NewReader(form.Encode()))
	feedbackRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(feedbackResponse, feedbackRequest)
	if feedbackResponse.Code != http.StatusOK {
		t.Fatalf("expected feedback OK, got %d: %s", feedbackResponse.Code, feedbackResponse.Body.String())
	}
	var feedback MLNote
	if err := db.Where("channel = ? AND target_kind = ? AND target_id = ?", noteChannelFeedback, noteTargetMatch, match.ID).First(&feedback).Error; err != nil {
		t.Fatalf("expected feedback note: %v", err)
	}
	if feedback.TeamID != team.ID || feedback.Creator != "Sideline" || feedback.Note != "Great pressure" {
		t.Fatalf("unexpected feedback: %+v", feedback)
	}
}

func TestActiveMatchForTeamIgnoresOtherTeamNearestMatch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	otherStart := time.Now().Add(30 * time.Minute)
	teamStart := time.Now().Add(2 * time.Hour)
	if err := db.Create(&Match{TeamID: 99, Opponent: "Other", StartTime: &otherStart}).Error; err != nil {
		t.Fatalf("create other match: %v", err)
	}
	teamMatch := Match{TeamID: 7, Opponent: "Team", StartTime: &teamStart}
	if err := db.Create(&teamMatch).Error; err != nil {
		t.Fatalf("create team match: %v", err)
	}

	match, err := activeMatchForTeam(db, 7)
	if err != nil {
		t.Fatalf("active match for team: %v", err)
	}
	if match.ID != teamMatch.ID {
		t.Fatalf("expected team match %d, got %d", teamMatch.ID, match.ID)
	}
}

func TestMatchDayLineupIDForTeamRequiresUpcomingMatchLineup(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	start := time.Now().Add(2 * time.Hour)
	match := Match{TeamID: 7, Opponent: "Upcoming", StartTime: &start}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	if got := matchDayLineupIDForTeam(db, 7); got != 0 {
		t.Fatalf("expected no match-day line-up before one is assigned, got %d", got)
	}
	formation := Formation{TeamID: 7, Name: "Formation", Status: formationStatusLive}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	lineup := Lineup{TeamID: 7, FormationID: formation.ID, Status: lineupStatusSelected, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}
	match.LineupID = lineup.ID
	if err := db.Model(&Match{}).Where("id = ?", match.ID).Update("lineup_id", lineup.ID).Error; err != nil {
		t.Fatalf("attach line-up: %v", err)
	}

	if got := matchDayLineupIDForTeam(db, 7); got != lineup.ID {
		t.Fatalf("expected match-day line-up %d, got %d", lineup.ID, got)
	}
}

func TestMatchDayAdminAssignsLineupWhenUpcomingMatchHasNoLineup(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team", EnableLineupsModule: true}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	start := time.Now().Add(2 * time.Hour)
	match := Match{TeamID: team.ID, Opponent: "Upcoming", Location: "Home", StartTime: &start}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Formation", Status: formationStatusLive}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, FormationID: formation.ID, Status: lineupStatusSelected, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-day", matchDayHandler(db))
	request := httptest.NewRequest(http.MethodGet, "/match-day", nil)
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	body := response.Body.String()
	if response.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d: %s", response.Code, body)
	}
	if !strings.Contains(body, "Assign a line-up to this match") {
		t.Fatalf("expected assignment prompt, got %s", body)
	}
	if !strings.Contains(body, `name="action" value="set-match"`) || !strings.Contains(body, `name="matchId" value="`+S(match.ID)+`"`) {
		t.Fatalf("expected set-match form for upcoming match, got %s", body)
	}
}

func TestMatchDayAdminWithoutUpcomingMatchLinksToManageMatches(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team", EnableLineupsModule: true}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-day", matchDayHandler(db))
	request := httptest.NewRequest(http.MethodGet, "/match-day", nil)
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	body := response.Body.String()
	if response.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d: %s", response.Code, body)
	}
	if !strings.Contains(body, "No upcoming match is ready for match day.") {
		t.Fatalf("expected no upcoming match warning, got %s", body)
	}
	if !strings.Contains(body, "Manage matches") || !strings.Contains(body, "manage-matches") {
		t.Fatalf("expected manage matches link, got %s", body)
	}
}

func TestMatchDayShowsMatchFocusPoints(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}, &MLNote{}, &PlayerMatchUnavailability{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team", EnableLineupsModule: true}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Formation", Status: formationStatusLive}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, FormationID: formation.ID, Status: lineupStatusSelected, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}
	start := time.Now().Add(2 * time.Hour)
	match := Match{TeamID: team.ID, LineupID: lineup.ID, Opponent: "Rovers", Location: "Home", StartTime: &start}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	if err := db.Create(&MLNote{TeamID: team.ID, Channel: noteChannelNote, TargetKind: noteTargetMatch, TargetID: match.ID, TargetLabel: "vs Rovers", Type: noteTypeFocus, Priority: 5, Creator: "admin", Note: "Press early from kick-off"}).Error; err != nil {
		t.Fatalf("create focus note: %v", err)
	}
	if err := db.Create(&MLNote{TeamID: team.ID, Channel: noteChannelNote, TargetKind: noteTargetMatch, TargetID: match.ID, TargetLabel: "vs Rovers", Type: "tactic", Priority: 5, Creator: "admin", Note: "Hidden broader note"}).Error; err != nil {
		t.Fatalf("create tactic note: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-day/{lineupId}", matchDayHandler(db))
	request := httptest.NewRequest(http.MethodGet, "/match-day/"+S(lineup.ID), nil)
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	body := response.Body.String()
	if response.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d: %s", response.Code, body)
	}
	if !strings.Contains(body, "Focus points") || !strings.Contains(body, "Press early from kick-off") {
		t.Fatalf("expected published focus point, got %s", body)
	}
	if strings.Contains(body, "Hidden broader note") {
		t.Fatalf("expected non-focus note to stay hidden, got %s", body)
	}
	if !strings.Contains(body, "Add focus point") || !strings.Contains(body, `name="type" value="focus"`) {
		t.Fatalf("expected admin focus-point form, got %s", body)
	}
}
