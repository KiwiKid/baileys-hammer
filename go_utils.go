package main

import (
	"context"
	"net/http"
)

// Define a Config struct to hold our configuration.
type Config struct {
	Title                 string
	UseRoles              bool
	UseMatchEventTracker  bool
	DefaultToApproved bool
	UsePlayerOfTheDayName string
	ShowGoalScorerMatchList bool
	ShowGoalAssistMatchList bool
	UseDudOfTheDayName string
	InjuryCounterTrackerName string
	ShowOpponentScore bool
}

var devConfig = &Config{
	Title: "🔨 Baileys Hammer 🔨",
	UseRoles: true,
	UseMatchEventTracker: true,
	InjuryCounterTrackerName: "Injuries (Mr Glass)",
	UsePlayerOfTheDayName: "Player of the Day",
	UseDudOfTheDayName: "Dick of the Day",
}

var config = &Config{
	Title: "🔨 Baileys Hammer 🔨",
	UseRoles: true,
	UseMatchEventTracker: true,
	/**

	If true, new fines will be approved by default (can be later decline by the finemaster)

	*/
	DefaultToApproved: true,

	ShowGoalScorerMatchList: true,
	ShowGoalAssistMatchList: true,
	InjuryCounterTrackerName: "🚑 Mr Glass 🚑",
	UsePlayerOfTheDayName: "Player of the Day",
	UseDudOfTheDayName: "🍆 Dick of the Day 🍆",
	ShowOpponentScore: true,
}

// Use a custom type for keys to avoid conflicts in context values.
type contextKey string

const (
	titleKey                 contextKey = "Title"
	useRolesKey              contextKey = "UseRoles"
	useMatchEventTrackerKey  contextKey = "UseMatchEventTracker"
	UsePlayerOfTheDayNameKey contextKey = "UsePlayerOfTheDayName"
	ShowGoalScorerMatchListKey contextKey = "ShowGoalScorerMatchList"
	ShowGoalAssistMatchListKey contextKey = "ShowGoalAssistMatchList"
	UseDudOfTheDayNameKey contextKey = "UseDudOfTheDayName"
	InjuryCounterTrackerNameKey contextKey = "InjuryCounterTrackerName"
	ShowOpponentScoreKey contextKey = "ShowOpponentScoreKey"
)

func GetTitle(ctx context.Context) string {
	if title, ok := ctx.Value(titleKey).(string); ok {
		return title
	}
	return ""
}

func UseRoles(ctx context.Context) bool {
	if useRoles, ok := ctx.Value(useRolesKey).(bool); ok {
		return useRoles
	}
	return false
}

func UseMatchEventTracker(ctx context.Context) bool {
	if useMatchEventTracker, ok := ctx.Value(useMatchEventTrackerKey).(bool); ok {
		return useMatchEventTracker
	}
	return false
}

func UsePlayerOfTheDayName(ctx context.Context) string {
	if usePlayerOfTheDay, ok := ctx.Value(UsePlayerOfTheDayNameKey).(string); ok {
		return usePlayerOfTheDay
	}
	return ""
}


func UseDudOfTheDayName(ctx context.Context) string {
	if useDudOfTheDay, ok := ctx.Value(UseDudOfTheDayNameKey).(string); ok {
		return useDudOfTheDay
	}
	return ""
}

func UseInjuryCounterTrackerName(ctx context.Context) string {
	if useInjuryCounterTrackerName, ok := ctx.Value(InjuryCounterTrackerNameKey).(string); ok {
		return useInjuryCounterTrackerName
	}
	return ""
}

func UseShowGoalScorerMatchList(ctx context.Context) bool {
	if ShowGoalScorerMatchList, ok := ctx.Value(ShowGoalScorerMatchListKey).(bool); ok {
		return ShowGoalScorerMatchList
	}
	return false
}


func UseShowGoalAssister(ctx context.Context) bool {
	if ShowGoalAssistMatchList, ok := ctx.Value(ShowGoalAssistMatchListKey).(bool); ok {
		return ShowGoalAssistMatchList
	}
	return false
}

func UseShowOpponentScore(ctx context.Context) bool {
	if ShowOpponentScore, ok := ctx.Value(ShowOpponentScoreKey).(bool); ok {
		return ShowOpponentScore
	}
	return false
}





func GetContext(r *http.Request) context.Context {
	
	ctx := context.WithValue(r.Context(), titleKey, config.Title)
	ctx = context.WithValue(ctx, useRolesKey, config.UseRoles)
	ctx = context.WithValue(ctx, useMatchEventTrackerKey, config.UseMatchEventTracker)
	ctx = context.WithValue(ctx, UsePlayerOfTheDayNameKey, config.UsePlayerOfTheDayName)
	ctx = context.WithValue(ctx, UseDudOfTheDayNameKey, config.UseDudOfTheDayName)
	ctx = context.WithValue(ctx, InjuryCounterTrackerNameKey, config.InjuryCounterTrackerName)
	ctx = context.WithValue(ctx, ShowGoalScorerMatchListKey, config.ShowGoalScorerMatchList)
	ctx = context.WithValue(ctx, ShowGoalAssistMatchListKey, config.ShowGoalAssistMatchList)
	ctx = context.WithValue(ctx, ShowOpponentScoreKey, config.ShowOpponentScore)

	return ctx
}
