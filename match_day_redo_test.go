package main

import (
	"testing"
	"time"
)

func TestMatchDayRedoEventRoundTrip(t *testing.T) {
	eventTime := time.Date(2026, 5, 31, 12, 34, 0, 0, time.UTC)
	event := MatchEvent{
		EventName:   "Sam - goal",
		EventType:   "goal",
		EventTime:   &eventTime,
		EventMinute: 17,
		PlayerId:    42,
	}

	values := matchDayRedoEventValues("Match event undone.", event)
	redoEvent := matchDayRedoEventFromQuery(values)

	if redoEvent == nil {
		t.Fatalf("expected redo event")
	}
	if redoEvent.EventName != event.EventName {
		t.Fatalf("expected event name %q, got %q", event.EventName, redoEvent.EventName)
	}
	if redoEvent.EventType != event.EventType {
		t.Fatalf("expected event type %q, got %q", event.EventType, redoEvent.EventType)
	}
	if redoEvent.EventMinute != event.EventMinute {
		t.Fatalf("expected minute %d, got %d", event.EventMinute, redoEvent.EventMinute)
	}
	if redoEvent.PlayerId != event.PlayerId {
		t.Fatalf("expected player ID %d, got %d", event.PlayerId, redoEvent.PlayerId)
	}
	if redoEvent.EventTime == nil || !redoEvent.EventTime.Equal(eventTime) {
		t.Fatalf("expected event time %v, got %v", eventTime, redoEvent.EventTime)
	}
}
