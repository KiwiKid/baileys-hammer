package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMatchDaySuggestedSubsShowsDuePlannedSub(t *testing.T) {
	startedAt := time.Now().Add(-35 * time.Minute)
	lineup := suggestedSubTestLineup(startedAt, nil)

	suggestions := matchDaySuggestedSubs(lineup, map[uint]bool{})

	require.Len(t, suggestions, 1)
	require.Equal(t, 1, suggestions[0].PositionIndex)
	require.Equal(t, 30, suggestions[0].DueMinute)
	require.Equal(t, uint(1), suggestions[0].CurrentPlayer.ID)
	require.Equal(t, uint(2), suggestions[0].ReplacementPlayer.ID)
	require.False(t, matchDayPlayerActuallyCurrentlyOn(lineup, 2))
}

func TestMatchDaySuggestedSubsSkipsLoggedSub(t *testing.T) {
	startedAt := time.Now().Add(-35 * time.Minute)
	eventTime := startedAt.Add(30 * time.Minute)
	lineup := suggestedSubTestLineup(startedAt, []MatchEvent{
		{
			Model:       gorm.Model{ID: 2, CreatedAt: eventTime},
			MatchId:     1,
			EventName:   "Starter subbed off",
			EventType:   "subbed-off",
			EventTime:   &eventTime,
			EventMinute: 30,
			PlayerId:    1,
		},
		{
			Model:       gorm.Model{ID: 3, CreatedAt: eventTime.Add(time.Second)},
			MatchId:     1,
			EventName:   "Replacement subbed on",
			EventType:   "subbed-on",
			EventTime:   &eventTime,
			EventMinute: 30,
			PlayerId:    2,
		},
	})

	suggestions := matchDaySuggestedSubs(lineup, map[uint]bool{})

	require.Empty(t, suggestions)
	require.True(t, matchDayPlayerActuallyCurrentlyOn(lineup, 2))
}

func suggestedSubTestLineup(startedAt time.Time, events []MatchEvent) *Lineup {
	events = append([]MatchEvent{
		{
			Model:       gorm.Model{ID: 1, CreatedAt: startedAt},
			MatchId:     1,
			EventName:   "Match started",
			EventType:   matchDayEventStarted,
			EventTime:   &startedAt,
			EventMinute: 0,
		},
	}, events...)
	return &Lineup{
		Model: gorm.Model{ID: 1},
		Name:  "Suggested sub test",
		Formation: Formation{
			Positions: []FormationPosition{
				{IndexNumber: 1, PositionName: "Striker"},
			},
		},
		Players: []LineupPlayer{
			{
				IndexNumber: 1,
				SlotOrder:   0,
				PlayerID:    1,
				Player:      Player{Model: gorm.Model{ID: 1}, Name: "Starter"},
			},
			{
				IndexNumber: 1,
				SlotOrder:   1,
				SubMinute:   30,
				PlayerID:    2,
				Player:      Player{Model: gorm.Model{ID: 2}, Name: "Replacement"},
			},
		},
		Match: Match{
			Model:  gorm.Model{ID: 1},
			Events: events,
		},
	}
}
