package main

import (
	"testing"

	"gorm.io/gorm"
)

func TestLineupPlayerSelectOptionsPrefersPlayablePositions(t *testing.T) {
	players := []Player{
		{Model: gorm.Model{ID: 1}, Name: "Charlie", PlayablePositions: "GK"},
		{Model: gorm.Model{ID: 2}, Name: "Alex", PlayablePositions: "Striker"},
		{Model: gorm.Model{ID: 3}, Name: "Bailey"},
	}

	options := lineupPlayerSelectOptions(players, nil, []FormationPosition{
		{IndexNumber: 9, PositionName: "Striker"},
	}, 9, "Striker", 0, nil)

	if len(options) != 3 {
		t.Fatalf("expected 3 options, got %d", len(options))
	}
	if options[0].Player.ID != 2 {
		t.Fatalf("expected striker match first, got %s", options[0].Player.Name)
	}
}

func TestPlayerPitchLabelUsesNumberAndShortName(t *testing.T) {
	player := Player{
		Name:   "Taylor Morgan",
		Number: "8",
	}

	if got := playerPitchLabel(player); got != "#8 T. Morgan" {
		t.Fatalf("expected numbered short label, got %q", got)
	}
	if got := playerPitchShirtNumber(player, 9); got != "8" {
		t.Fatalf("expected shirt number, got %q", got)
	}
}
