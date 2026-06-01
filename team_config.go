package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultLineupPlayerCount = 11
	minLineupPlayerCount     = 1
	maxLineupPlayerCount     = 15
)

func teamLineupPlayerCount(team Team) int {
	if team.LineupPlayerCount == 0 {
		return defaultLineupPlayerCount
	}
	return normalizeLineupPlayerCount(team.LineupPlayerCount)
}

func normalizeLineupPlayerCount(value int) int {
	if value < minLineupPlayerCount {
		return minLineupPlayerCount
	}
	if value > maxLineupPlayerCount {
		return maxLineupPlayerCount
	}
	return value
}

func parseLineupPlayerCount(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultLineupPlayerCount, nil
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if count < minLineupPlayerCount || count > maxLineupPlayerCount {
		return 0, fmt.Errorf("line-up player count must be between %d and %d", minLineupPlayerCount, maxLineupPlayerCount)
	}
	return count, nil
}

func sideLabel(playerCount int) string {
	playerCount = normalizeLineupPlayerCount(playerCount)
	return fmt.Sprintf("%d-a-side", playerCount)
}
