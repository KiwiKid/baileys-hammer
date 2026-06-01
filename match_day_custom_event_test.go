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

func TestMatchDayCustomEventCreatesMatchEvent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	start := time.Now().Add(-10 * time.Minute)
	match := Match{TeamID: team.ID, StartTime: &start, Opponent: "Opposition"}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	player := Player{Name: "Player", Active: true}
	if err := db.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, Status: lineupStatusUsed, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}
	match.LineupID = lineup.ID
	if err := db.Model(&Match{}).Where("id = ?", match.ID).Update("lineup_id", lineup.ID).Error; err != nil {
		t.Fatalf("attach line-up: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-day/{lineupId}/action", matchDayActionHandler(db))

	form := url.Values{}
	form.Set("action", "custom-event")
	form.Set("eventType", "yellow-card")
	form.Set("eventName", "Late yellow card")
	form.Set("minute", "72")
	form.Set("playerId", S(player.ID))
	request := httptest.NewRequest(http.MethodPost, "/match-day/"+S(lineup.ID)+"/action", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", response.Code)
	}
	var event MatchEvent
	if err := db.Where("match_id = ?", match.ID).First(&event).Error; err != nil {
		t.Fatalf("load event: %v", err)
	}
	if event.EventType != "yellow-card" || event.EventName != "Late yellow card" || event.EventMinute != 72 || event.PlayerId != player.ID {
		t.Fatalf("unexpected event: %#v", event)
	}
}

func TestMatchDayResetRequiresPreviewConfirmationAndDeletesEvents(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, Status: lineupStatusUsed, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}
	start := time.Now().Add(-10 * time.Minute)
	match := Match{TeamID: team.ID, LineupID: lineup.ID, StartTime: &start, Opponent: "Opposition"}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	for _, eventType := range []string{"start-match", "goal", "assist"} {
		if err := db.Create(&MatchEvent{MatchId: uint64(match.ID), EventType: eventType, EventName: eventType}).Error; err != nil {
			t.Fatalf("create event: %v", err)
		}
	}

	router := chi.NewRouter()
	router.HandleFunc("/match-day/{lineupId}/action", matchDayActionHandler(db))

	previewForm := url.Values{}
	previewForm.Set("action", "preview-reset-match")
	previewRequest := httptest.NewRequest(http.MethodPost, "/match-day/"+S(lineup.ID)+"/action", strings.NewReader(previewForm.Encode()))
	previewRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	previewRequest.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	previewResponse := httptest.NewRecorder()
	router.ServeHTTP(previewResponse, previewRequest)
	if previewResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected preview redirect, got %d", previewResponse.Code)
	}
	if location := previewResponse.Header().Get("Location"); !strings.Contains(location, "resetPreview=true") {
		t.Fatalf("expected reset preview redirect, got %q", location)
	}

	unconfirmedForm := url.Values{}
	unconfirmedForm.Set("action", "reset-match")
	unconfirmedRequest := httptest.NewRequest(http.MethodPost, "/match-day/"+S(lineup.ID)+"/action", strings.NewReader(unconfirmedForm.Encode()))
	unconfirmedRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	unconfirmedRequest.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	unconfirmedResponse := httptest.NewRecorder()
	router.ServeHTTP(unconfirmedResponse, unconfirmedRequest)
	if unconfirmedResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected unconfirmed redirect, got %d", unconfirmedResponse.Code)
	}
	var count int64
	if err := db.Model(&MatchEvent{}).Where("match_id = ?", match.ID).Count(&count).Error; err != nil {
		t.Fatalf("count events: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected unconfirmed reset to keep 3 events, got %d", count)
	}

	confirmForm := url.Values{}
	confirmForm.Set("action", "reset-match")
	confirmForm.Set("confirmReset", "delete-events")
	confirmRequest := httptest.NewRequest(http.MethodPost, "/match-day/"+S(lineup.ID)+"/action", strings.NewReader(confirmForm.Encode()))
	confirmRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	confirmRequest.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	confirmResponse := httptest.NewRecorder()
	router.ServeHTTP(confirmResponse, confirmRequest)
	if confirmResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected confirmed redirect, got %d", confirmResponse.Code)
	}
	if err := db.Model(&MatchEvent{}).Where("match_id = ?", match.ID).Count(&count).Error; err != nil {
		t.Fatalf("count events after reset: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected reset to delete all events, got %d", count)
	}
}
