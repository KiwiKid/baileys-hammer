package main

import (
	"testing"
	"time"
)

func TestMatchDayClockPausesAndResumes(t *testing.T) {
	now := time.Now()
	start := now.Add(-60 * time.Minute)
	halfTime := now.Add(-10 * time.Minute)
	match := Match{
		StartTime: &start,
		Events: []MatchEvent{
			{
				EventType:   matchDayEventHalfTime,
				EventTime:   &halfTime,
				EventMinute: matchDayRegularHalfMinute,
			},
		},
	}

	paused := matchDayClockStateForMatch(match)
	if !paused.Paused {
		t.Fatalf("expected clock to be paused")
	}
	if paused.CurrentMinute != matchDayRegularHalfMinute {
		t.Fatalf("expected paused minute %d, got %d", matchDayRegularHalfMinute, paused.CurrentMinute)
	}
	if !matchDayCanShowResume(match) {
		t.Fatalf("expected resume controls to be available")
	}
	if matchDayClockStatusLabel(match) != "Paused for half-time" {
		t.Fatalf("expected half-time pause reason, got %q", matchDayClockStatusLabel(match))
	}

	resumedAt := now.Add(-5 * time.Minute)
	match.Events = append(match.Events, MatchEvent{
		EventType:   matchDayEventResumed,
		EventTime:   &resumedAt,
		EventMinute: matchDayRegularHalfMinute,
	})
	resumed := matchDayClockStateForMatch(match)
	if resumed.Paused {
		t.Fatalf("expected resumed clock not to be paused")
	}
	if resumed.CurrentMinute < 49 || resumed.CurrentMinute > 51 {
		t.Fatalf("expected resumed minute around 50, got %d", resumed.CurrentMinute)
	}
}

func TestMatchDayClockExtraTimeHalfTimeVisibility(t *testing.T) {
	now := time.Now()
	start := now.Add(-130 * time.Minute)
	extraTimeStart := now.Add(-16 * time.Minute)
	match := Match{
		StartTime: &start,
		Events: []MatchEvent{
			{
				EventType:   matchDayEventExtraTimeStart,
				EventTime:   &extraTimeStart,
				EventMinute: matchDayRegularEndMinute,
			},
		},
	}

	state := matchDayClockStateForMatch(match)
	if !state.InExtraTime {
		t.Fatalf("expected extra-time state")
	}
	if state.CurrentMinute < matchDayExtraHalfMinute {
		t.Fatalf("expected clock to be at least extra-time half-time, got %d", state.CurrentMinute)
	}
	if !matchDayCanShowHalfTime(match) {
		t.Fatalf("expected second half-time control to be visible")
	}
}

func TestMatchDayClockFinished(t *testing.T) {
	now := time.Now()
	start := now.Add(-100 * time.Minute)
	finishedAt := now.Add(-1 * time.Minute)
	match := Match{
		StartTime: &start,
		Events: []MatchEvent{
			{
				EventType:   matchDayEventFinished,
				EventTime:   &finishedAt,
				EventMinute: 92,
			},
		},
	}

	state := matchDayClockStateForMatch(match)
	if !state.Finished {
		t.Fatalf("expected finished clock state")
	}
	if state.CurrentMinute != 92 {
		t.Fatalf("expected finished minute 92, got %d", state.CurrentMinute)
	}
	if matchDayCanShowMatchFinished(match) {
		t.Fatalf("expected finish control to be hidden after finish event")
	}
	if matchDayClockStatusLabel(match) != "Game ended" {
		t.Fatalf("expected game ended status, got %q", matchDayClockStatusLabel(match))
	}
}

func TestMatchDayPauseReasonFallback(t *testing.T) {
	if got := matchDayPauseReasonLabel("match-paused"); got != "Game pause" {
		t.Fatalf("expected fallback game pause label, got %q", got)
	}
}
