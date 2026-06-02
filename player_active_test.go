package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPlayerHandlerCanActivateAndDeactivatePlayer(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Season{}, &Player{}, &Fine{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	player := Player{Name: "Morgan", Active: true}
	if err := db.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	postPlayerForm(t, db, url.Values{
		"ID":           {S(player.ID)},
		"Name":         {"Morgan"},
		"Active":       {"true"},
		"playerAction": {"deactivate"},
	})
	var updated Player
	if err := db.First(&updated, player.ID).Error; err != nil {
		t.Fatalf("fetch deactivated player: %v", err)
	}
	if updated.Active {
		t.Fatalf("expected player to be deactivated")
	}

	postPlayerForm(t, db, url.Values{
		"ID":           {S(player.ID)},
		"Name":         {"Morgan"},
		"Active":       {"false"},
		"playerAction": {"activate"},
	})
	if err := db.First(&updated, player.ID).Error; err != nil {
		t.Fatalf("fetch activated player: %v", err)
	}
	if !updated.Active {
		t.Fatalf("expected player to be activated")
	}
}

func TestPlayerHandlerPreservesActiveStateOnNormalUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Season{}, &Player{}, &Fine{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	player := Player{Name: "Taylor", Active: true}
	if err := db.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	postPlayerForm(t, db, url.Values{
		"ID":                {S(player.ID)},
		"Name":              {"Taylor Updated"},
		"number":            {"8"},
		"playablePositions": {"CM, ST"},
	})
	var updated Player
	if err := db.First(&updated, player.ID).Error; err != nil {
		t.Fatalf("fetch updated player: %v", err)
	}
	if !updated.Active {
		t.Fatalf("expected normal update to preserve active state")
	}
	if updated.Name != "Taylor Updated" {
		t.Fatalf("expected player name to update, got %q", updated.Name)
	}
	if updated.Number != "8" {
		t.Fatalf("expected player number to update, got %q", updated.Number)
	}
	if updated.PlayablePositions != "CM, ST" {
		t.Fatalf("expected playable positions to update, got %q", updated.PlayablePositions)
	}
}

func postPlayerForm(t *testing.T, db *gorm.DB, form url.Values) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/players", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	playerHandler(db).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("post player returned %d: %s", response.Code, response.Body.String())
	}
}
