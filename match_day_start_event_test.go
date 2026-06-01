package main

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMatchDayClockUsesStartMatchEvent(t *testing.T) {
	now := time.Now()
	scheduledStart := now.Add(-60 * time.Minute)
	actualStart := now.Add(-12 * time.Minute)
	match := Match{
		StartTime: &scheduledStart,
		Events: []MatchEvent{
			{
				EventType:   matchDayEventStarted,
				EventTime:   &actualStart,
				EventMinute: 0,
			},
		},
	}

	state := matchDayClockStateForMatch(match)

	if !state.Started {
		t.Fatalf("expected clock to be started")
	}
	if state.CurrentMinute < 11 || state.CurrentMinute > 13 {
		t.Fatalf("expected current minute around 12 from start event, got %d", state.CurrentMinute)
	}
}

func TestMatchDayOverrunWarningAfter135Minutes(t *testing.T) {
	start := time.Now().Add(-136 * time.Minute)
	match := Match{
		Events: []MatchEvent{
			{
				EventType:   matchDayEventStarted,
				EventTime:   &start,
				EventMinute: 0,
			},
		},
	}

	if !matchDayCanShowOverrunWarning(match) {
		t.Fatalf("expected overrun warning after 135 minutes")
	}
	if matchDayRawMinute(match) < 136 {
		t.Fatalf("expected raw minute to remain above 135, got %d", matchDayRawMinute(match))
	}
	if matchDayMinute(match) != matchDayExtraEndMinute {
		t.Fatalf("expected normal match minute to remain capped at %d, got %d", matchDayExtraEndMinute, matchDayMinute(match))
	}
}

func TestStartLineupMatchCreatesAndUpdatesStartMatchEvent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	match := Match{}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	lineup := Lineup{Status: lineupStatusSelected}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	match.LineupID = lineup.ID
	if err := db.Model(&Match{}).Where("id = ?", match.ID).Update("lineup_id", lineup.ID).Error; err != nil {
		t.Fatalf("attach line-up: %v", err)
	}
	match.LineupID = lineup.ID
	lineup.Match = match

	firstStart := time.Date(2026, 5, 31, 12, 0, 0, 0, time.UTC)
	if err := startLineupMatch(db, &lineup, firstStart); err != nil {
		t.Fatalf("start match: %v", err)
	}

	var events []MatchEvent
	if err := db.Where("match_id = ? AND event_type = ?", match.ID, matchDayEventStarted).Find(&events).Error; err != nil {
		t.Fatalf("find start events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one start event, got %d", len(events))
	}
	if events[0].EventTime == nil || !events[0].EventTime.Equal(firstStart) {
		t.Fatalf("expected start event time %v, got %v", firstStart, events[0].EventTime)
	}
	if err := db.First(&lineup, lineup.ID).Error; err != nil {
		t.Fatalf("reload lineup: %v", err)
	}
	if lineup.Status != lineupStatusSelected {
		t.Fatalf("expected shared line-up status to stay selected, got %q", lineup.Status)
	}
	var initEvents []MatchEvent
	if err := db.Where("match_id = ? AND event_type = ?", match.ID, matchDayEventInit).Find(&initEvents).Error; err != nil {
		t.Fatalf("find init events: %v", err)
	}
	if len(initEvents) != 1 {
		t.Fatalf("expected one init event, got %d", len(initEvents))
	}
	if initEvents[0].EventData == "" {
		t.Fatalf("expected init event data")
	}

	updatedStart := firstStart.Add(2 * time.Minute)
	if err := startLineupMatch(db, &lineup, updatedStart); err != nil {
		t.Fatalf("adjust start match: %v", err)
	}
	events = []MatchEvent{}
	if err := db.Where("match_id = ? AND event_type = ?", match.ID, matchDayEventStarted).Find(&events).Error; err != nil {
		t.Fatalf("find adjusted start events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one adjusted start event, got %d", len(events))
	}
	if events[0].EventTime == nil || !events[0].EventTime.Equal(updatedStart) {
		t.Fatalf("expected adjusted start event time %v, got %v", updatedStart, events[0].EventTime)
	}
}
