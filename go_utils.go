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
}

var config = &Config{
	Title: "🔨 Baileys Hammer 🔨🔨",
	UseRoles: false,
	UseMatchEventTracker: false,
}

// Use a custom type for keys to avoid conflicts in context values.
type contextKey string

const (
	titleKey                 contextKey = "Title"
	useRolesKey              contextKey = "UseRoles"
	useMatchEventTrackerKey  contextKey = "UseMatchEventTracker"
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

func GetContext(r *http.Request) context.Context {
	ctx := context.WithValue(r.Context(), titleKey, config.Title)
	ctx = context.WithValue(ctx, useRolesKey, config.UseRoles)
	ctx = context.WithValue(ctx, useMatchEventTrackerKey, config.UseMatchEventTracker)
	return ctx
}
