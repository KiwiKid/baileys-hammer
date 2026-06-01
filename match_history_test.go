package main

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetHistoricMatchesOrdersPastMatchesForActiveTeam(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &Player{}, &PlayerMatchUnavailability{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}

	now := time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	oldStart := now.Add(-72 * time.Hour)
	recentStart := now.Add(-2 * time.Hour)
	futureStart := now.Add(2 * time.Hour)
	otherTeamStart := now.Add(-1 * time.Hour)

	oldMatch := Match{TeamID: 7, Opponent: "Old", StartTime: &oldStart}
	recentMatch := Match{TeamID: 7, Opponent: "Recent", StartTime: &recentStart}
	futureMatch := Match{TeamID: 7, Opponent: "Future", StartTime: &futureStart}
	otherTeamMatch := Match{TeamID: 99, Opponent: "Other", StartTime: &otherTeamStart}
	for _, match := range []*Match{&oldMatch, &recentMatch, &futureMatch, &otherTeamMatch} {
		if err := db.Create(match).Error; err != nil {
			t.Fatalf("create match: %v", err)
		}
	}

	matches, total, err := GetHistoricMatches(db, 7, 1, 20, now)
	if err != nil {
		t.Fatalf("historic matches: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 historic matches, got %d", total)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].ID != recentMatch.ID || matches[1].ID != oldMatch.ID {
		t.Fatalf("expected recent then old, got %#v", matches)
	}
}

func TestHistoricMatchTotalPages(t *testing.T) {
	tests := []struct {
		name  string
		total int
		limit int
		want  int
	}{
		{name: "empty", total: 0, limit: 20, want: 1},
		{name: "partial", total: 19, limit: 20, want: 1},
		{name: "exact", total: 20, limit: 20, want: 1},
		{name: "extra", total: 21, limit: 20, want: 2},
		{name: "bad limit", total: 21, limit: 0, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := historicMatchTotalPages(tt.total, tt.limit); got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
