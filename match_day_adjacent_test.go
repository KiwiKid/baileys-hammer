package main

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMatchDayAdjacentLineupsFindsPreviousAndNextMatchByStartTime(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &Lineup{}, &Formation{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}

	teamID := uint(7)
	formation := Formation{TeamID: teamID, Name: "Live", Status: formationStatusLive}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	base := time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	prevStart := base.Add(-48 * time.Hour)
	currentStart := base.Add(-24 * time.Hour)
	nextStart := base.Add(24 * time.Hour)
	otherTeamStart := base.Add(-12 * time.Hour)

	prevMatch := Match{TeamID: teamID, Opponent: "Previous", StartTime: &prevStart}
	currentMatch := Match{TeamID: teamID, Opponent: "Current", StartTime: &currentStart}
	nextMatch := Match{TeamID: teamID, Opponent: "Next", StartTime: &nextStart}
	otherTeamMatch := Match{TeamID: 99, Opponent: "Other", StartTime: &otherTeamStart}
	for _, match := range []*Match{&prevMatch, &currentMatch, &nextMatch, &otherTeamMatch} {
		if err := db.Create(match).Error; err != nil {
			t.Fatalf("create match: %v", err)
		}
	}

	prevLineup := Lineup{TeamID: teamID, FormationID: formation.ID, Name: "Previous"}
	currentLineup := Lineup{TeamID: teamID, FormationID: formation.ID, Name: "Current"}
	nextLineup := Lineup{TeamID: teamID, FormationID: formation.ID, Name: "Next"}
	otherTeamLineup := Lineup{TeamID: 99, FormationID: formation.ID, Name: "Other"}
	for _, lineup := range []*Lineup{&prevLineup, &currentLineup, &nextLineup, &otherTeamLineup} {
		if err := db.Create(lineup).Error; err != nil {
			t.Fatalf("create lineup: %v", err)
		}
	}
	for _, pair := range []struct {
		match  *Match
		lineup *Lineup
	}{
		{&prevMatch, &prevLineup},
		{&currentMatch, &currentLineup},
		{&nextMatch, &nextLineup},
		{&otherTeamMatch, &otherTeamLineup},
	} {
		pair.match.LineupID = pair.lineup.ID
		if err := db.Model(&Match{}).Where("id = ?", pair.match.ID).Update("lineup_id", pair.lineup.ID).Error; err != nil {
			t.Fatalf("attach lineup: %v", err)
		}
		pair.match.LineupID = pair.lineup.ID
	}
	currentLineup.Match = currentMatch

	adjacent := matchDayAdjacentLineups(db, LineupActor{TeamID: teamID, IsAdmin: true}, &currentLineup)
	if adjacent.Previous == nil || adjacent.Previous.ID != prevLineup.ID {
		t.Fatalf("expected previous lineup %d, got %#v", prevLineup.ID, adjacent.Previous)
	}
	if adjacent.Next == nil || adjacent.Next.ID != nextLineup.ID {
		t.Fatalf("expected next lineup %d, got %#v", nextLineup.ID, adjacent.Next)
	}
}

func TestMatchDayAdjacentLineupsRequiresSelectedMatchDate(t *testing.T) {
	adjacent := matchDayAdjacentLineups(nil, LineupActor{TeamID: 7}, &Lineup{Match: Match{Model: gorm.Model{ID: 1}}})
	if adjacent.Previous != nil || adjacent.Next != nil {
		t.Fatalf("expected no adjacent lineups without selected match date, got %#v", adjacent)
	}
}

func TestMatchDayIsHistorical(t *testing.T) {
	now := time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	oldStart := now.Add(-time.Hour)
	futureStart := now.Add(time.Hour)

	if !matchDayIsHistorical(&Lineup{Match: Match{Model: gorm.Model{ID: 1}, StartTime: &oldStart}}, now) {
		t.Fatal("expected old match to be historical")
	}
	if matchDayIsHistorical(&Lineup{Match: Match{Model: gorm.Model{ID: 1}, StartTime: &futureStart}}, now) {
		t.Fatal("expected future match not to be historical")
	}
	if matchDayIsHistorical(&Lineup{}, now) {
		t.Fatal("expected lineup without match not to be historical")
	}
}
