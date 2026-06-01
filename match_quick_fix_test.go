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

func TestMatchQuickFixSetsCurrentSeasonAndTeam(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Season{}, &Match{}, &MatchEvent{}, &Player{}, &PlayerMatchUnavailability{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &LineupUser{}, &MLNote{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Current team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	season := Season{Title: "Current season", StartDate: time.Now(), IsActive: true}
	if err := db.Create(&season).Error; err != nil {
		t.Fatalf("create season: %v", err)
	}
	match := Match{Opponent: "Missing owner", SeasonId: 0, TeamID: 0}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/match/{matchId}", matchHandler(db))
	form := url.Values{}
	form.Set("matchId", S(match.ID))
	form.Set("matchAction", "quick-fix-season-team")
	request := httptest.NewRequest(http.MethodPost, "/match/"+S(match.ID), strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected OK, got %d body=%s", response.Code, response.Body.String())
	}
	if err := db.First(&match, match.ID).Error; err != nil {
		t.Fatalf("reload match: %v", err)
	}
	if match.TeamID != team.ID {
		t.Fatalf("expected team %d, got %d", team.ID, match.TeamID)
	}
	if match.SeasonId != uint64(season.ID) {
		t.Fatalf("expected season %d, got %d", season.ID, match.SeasonId)
	}
}
