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
	UseDudOfTheDayName string
	InjuryCounterTrackerName string
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
	InjuryCounterTrackerName: "🚑 Mr Glass 🚑",
	UsePlayerOfTheDayName: "Player of the Day",
	UseDudOfTheDayName: "🍆 Dick of the Day 🍆",
}

// Use a custom type for keys to avoid conflicts in context values.
type contextKey string

const (
	titleKey                 contextKey = "Title"
	useRolesKey              contextKey = "UseRoles"
	useMatchEventTrackerKey  contextKey = "UseMatchEventTracker"
	UsePlayerOfTheDayNameKey contextKey = "UsePlayerOfTheDayName"
	UseDudOfTheDayNameKey contextKey = "UseDudOfTheDayName"
	InjuryCounterTrackerNameKey contextKey = "InjuryCounterTrackerName"
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



func GetContext(r *http.Request) context.Context {
	
	ctx := context.WithValue(r.Context(), titleKey, config.Title)
	ctx = context.WithValue(ctx, useRolesKey, config.UseRoles)
	ctx = context.WithValue(ctx, useMatchEventTrackerKey, config.UseMatchEventTracker)
	ctx = context.WithValue(ctx, UsePlayerOfTheDayNameKey, config.UsePlayerOfTheDayName)
	ctx = context.WithValue(ctx, UseDudOfTheDayNameKey, config.UseDudOfTheDayName)
	ctx = context.WithValue(ctx, InjuryCounterTrackerNameKey, config.InjuryCounterTrackerName)
	return ctx
}
