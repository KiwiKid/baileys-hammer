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
	UsePlayerOfTheDayName string
	UseDudOfTheDayName string
}

var config = &Config{
	Title: "🔨 Baileys WOAH Hammer 🔨",
	UseRoles: true,
	UseMatchEventTracker: false,
	UsePlayerOfTheDayName: "",
	UseDudOfTheDayName: "",
	// UsePlayerOfTheDayName defines the label for recognizing the best player of the day (empty for off).
	//UsePlayerOfTheDayName: "Player of the Day",
	// UseDudOfTheDayName defines the label for pointing out the least effective player of the day (empty for off).
	//UseDudOfTheDayName: "Dick of the Day",
}

// Use a custom type for keys to avoid conflicts in context values.
type contextKey string

const (
	titleKey                 contextKey = "Title"
	useRolesKey              contextKey = "UseRoles"
	useMatchEventTrackerKey  contextKey = "UseMatchEventTracker"
	UsePlayerOfTheDayNameKey contextKey = "UsePlayerOfTheDayName"
	UseDudOfTheDayNameKey contextKey = "UseDudOfTheDayName"
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



func GetContext(r *http.Request) context.Context {
	ctx := context.WithValue(r.Context(), titleKey, config.Title)
	ctx = context.WithValue(ctx, useRolesKey, config.UseRoles)
	ctx = context.WithValue(ctx, useMatchEventTrackerKey, config.UseMatchEventTracker)
	ctx = context.WithValue(ctx, UsePlayerOfTheDayNameKey, config.UsePlayerOfTheDayName)
	ctx = context.WithValue(ctx, UseDudOfTheDayNameKey, config.UseDudOfTheDayName)
	return ctx
}
