package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeleteLineupRemovesLineupAndPlayers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Match{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Delete team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, Name: "Delete me", Status: lineupStatusDraft}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	if err := db.Create(&LineupPlayer{LineupID: lineup.ID, IndexNumber: 1, SlotOrder: 0, PlayerID: 1}).Error; err != nil {
		t.Fatalf("create lineup player: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/lineups/{lineupId}", lineupDetailHandler(db))
	request := httptest.NewRequest(http.MethodDelete, "/lineups/"+S(lineup.ID), nil)
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", response.Code)
	}
	var lineupCount int64
	if err := db.Model(&Lineup{}).Where("id = ?", lineup.ID).Count(&lineupCount).Error; err != nil {
		t.Fatalf("count lineups: %v", err)
	}
	if lineupCount != 0 {
		t.Fatalf("expected lineup to be deleted")
	}
	var playerCount int64
	if err := db.Model(&LineupPlayer{}).Where("lineup_id = ?", lineup.ID).Count(&playerCount).Error; err != nil {
		t.Fatalf("count lineup players: %v", err)
	}
	if playerCount != 0 {
		t.Fatalf("expected lineup players to be deleted")
	}
}

func TestDeleteFormationRemovesFormationAndPositions(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Formation{}, &FormationPosition{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Delete team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Delete shape", Status: formationStatusDraft}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	if err := db.Create(&FormationPosition{FormationID: formation.ID, IndexNumber: 1, PositionName: "GK"}).Error; err != nil {
		t.Fatalf("create position: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/formations/{formationId}/edit", formationEditHandler(db))
	request := httptest.NewRequest(http.MethodDelete, "/formations/"+S(formation.ID)+"/edit", nil)
	request.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect, got %d", response.Code)
	}
	var formationCount int64
	if err := db.Model(&Formation{}).Where("id = ?", formation.ID).Count(&formationCount).Error; err != nil {
		t.Fatalf("count formations: %v", err)
	}
	if formationCount != 0 {
		t.Fatalf("expected formation to be deleted")
	}
	var positionCount int64
	if err := db.Model(&FormationPosition{}).Where("formation_id = ?", formation.ID).Count(&positionCount).Error; err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if positionCount != 0 {
		t.Fatalf("expected formation positions to be deleted")
	}
}
