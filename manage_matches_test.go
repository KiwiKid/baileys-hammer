package main

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetManageMatchesIncludesUnseasonedUpcomingTeamMatch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &Player{}, &PlayerMatchUnavailability{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	upcomingStart := time.Now().Add(48 * time.Hour)
	oldUnseasonedStart := time.Now().Add(-48 * time.Hour)
	activeSeasonMatchStart := time.Now().Add(24 * time.Hour)

	upcoming := Match{TeamID: 7, SeasonId: 0, Opponent: "Unseasoned upcoming", StartTime: &upcomingStart}
	oldUnseasoned := Match{TeamID: 7, SeasonId: 0, Opponent: "Old unseasoned", StartTime: &oldUnseasonedStart}
	activeSeasonMatch := Match{TeamID: 7, SeasonId: 2, Opponent: "Active season", StartTime: &activeSeasonMatchStart}
	otherTeam := Match{TeamID: 99, SeasonId: 2, Opponent: "Other team", StartTime: &activeSeasonMatchStart}
	for _, match := range []*Match{&upcoming, &oldUnseasoned, &activeSeasonMatch, &otherTeam} {
		if err := db.Create(match).Error; err != nil {
			t.Fatalf("create match: %v", err)
		}
	}

	matches, err := GetManageMatches(db, 7, 2, 0, 999)
	if err != nil {
		t.Fatalf("get manage matches: %v", err)
	}
	gotIDs := map[uint]bool{}
	for _, match := range matches {
		gotIDs[match.ID] = true
	}
	if !gotIDs[upcoming.ID] {
		t.Fatalf("expected unseasoned upcoming match in manage matches, got %#v", matches)
	}
	if !gotIDs[activeSeasonMatch.ID] {
		t.Fatalf("expected active-season team match in manage matches, got %#v", matches)
	}
	if gotIDs[oldUnseasoned.ID] {
		t.Fatalf("did not expect old unseasoned match in manage matches, got %#v", matches)
	}
	if gotIDs[otherTeam.ID] {
		t.Fatalf("did not expect other team match in manage matches, got %#v", matches)
	}
}
