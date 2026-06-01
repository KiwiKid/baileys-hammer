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
