package main

import (
	"strings"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSportySyncResultSkipFixtureRecordsReason(t *testing.T) {
	result := &SportySyncResult{}
	fixture := sportyFixture{
		ID:           123,
		HomeTeamID:   10,
		HomeTeamName: "Home",
		AwayTeamID:   20,
		AwayTeamName: "Away",
	}

	result.skipFixture(fixture, "wrong team")

	if result.Skipped != 1 {
		t.Fatalf("expected skipped count 1, got %d", result.Skipped)
	}
	if len(result.SkippedRecords) != 1 {
		t.Fatalf("expected one skipped record, got %d", len(result.SkippedRecords))
	}
	if result.SkippedRecords[0].Fixture.ID != fixture.ID {
		t.Fatalf("expected fixture ID %d, got %d", fixture.ID, result.SkippedRecords[0].Fixture.ID)
	}
	if result.SkippedRecords[0].Reason != "wrong team" {
		t.Fatalf("expected skip reason, got %q", result.SkippedRecords[0].Reason)
	}
}

func TestSportyFixtureSeasonSkipReasonBeforeSeasonStart(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Season{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	season := Season{
		Title:     "Winter 2026",
		StartDate: time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local),
		IsActive:  true,
	}
	if err := db.Create(&season).Error; err != nil {
		t.Fatalf("create season: %v", err)
	}

	reason, skip, err := sportyFixtureSeasonSkipReason(db, sportyFixture{ID: 123, From: "2026-05-31T12:00:00"})
	if err != nil {
		t.Fatalf("season skip reason: %v", err)
	}
	if !skip {
		t.Fatal("expected fixture before season start to be skipped")
	}
	if !strings.Contains(reason, "before active season") {
		t.Fatalf("expected season-start reason, got %q", reason)
	}
}

func TestSportyFixtureSeasonSkipReasonOnOrAfterSeasonStart(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Season{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	season := Season{
		Title:     "Winter 2026",
		StartDate: time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local),
		IsActive:  true,
	}
	if err := db.Create(&season).Error; err != nil {
		t.Fatalf("create season: %v", err)
	}

	reason, skip, err := sportyFixtureSeasonSkipReason(db, sportyFixture{ID: 123, From: "2026-06-01T00:00:00"})
	if err != nil {
		t.Fatalf("season skip reason: %v", err)
	}
	if skip {
		t.Fatalf("expected fixture on season start not to be skipped, reason %q", reason)
	}
}
