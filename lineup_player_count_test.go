package main

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLineupProblemsUseTeamPlayerCount(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Formation{}, &FormationPosition{}, &Match{}, &Lineup{}, &LineupPlayer{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Small side", LineupPlayerCount: 7}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Five-ish"}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := db.Create(&FormationPosition{FormationID: formation.ID, IndexNumber: i, PositionName: "P"}).Error; err != nil {
			t.Fatalf("create position: %v", err)
		}
	}
	lineup := Lineup{TeamID: team.ID, FormationID: formation.ID}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	loaded, err := getLineup(db, lineup.ID)
	if err != nil {
		t.Fatalf("load lineup: %v", err)
	}

	problems, err := lineupProblems(db, loaded)
	if err != nil {
		t.Fatalf("lineup problems: %v", err)
	}
	if !hasLineupProblem(problems, "formation has 5 positions; team is configured for 7-a-side") {
		t.Fatalf("expected 7-a-side formation problem, got %#v", problems)
	}
}

func TestSaveFormationPositionsRespectsTeamPlayerCount(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Formation{}, &FormationPosition{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Five side", LineupPlayerCount: 5}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	formation := Formation{TeamID: team.ID, Name: "Full"}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := db.Create(&FormationPosition{FormationID: formation.ID, IndexNumber: i, PositionName: "P"}).Error; err != nil {
			t.Fatalf("create position: %v", err)
		}
	}

	form := url.Values{}
	form.Set("newPositionName", "Extra")
	form.Set("newPositionIndex", "5")
	form.Set("newPositionX", "50")
	form.Set("newPositionY", "50")
	request := httptest.NewRequest("POST", "/formations/1/edit", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := request.ParseForm(); err != nil {
		t.Fatalf("parse form: %v", err)
	}

	err = saveFormationPositions(db, formation.ID, request)
	if err == nil || !strings.Contains(err.Error(), "5-a-side") {
		t.Fatalf("expected 5-a-side limit error, got %v", err)
	}
}

func hasLineupProblem(problems []LineupProblem, msg string) bool {
	for _, problem := range problems {
		if problem.Msg == msg {
			return true
		}
	}
	return false
}
