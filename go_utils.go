package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

// Define a Config struct to hold our configuration.
type Config struct {
	Title                    string
	UsePreviewPassword       bool
	UseRoles                 bool
	UseMatchEventTracker     bool
	DefaultToApproved        bool
	UsePlayerOfTheDayName    string
	ShowGoalScorerMatchList  bool
	ShowGoalAssistMatchList  bool
	UseDudOfTheDayName       string
	InjuryCounterTrackerName string
	ShowOpponentScore        bool
}

var devConfig = &Config{
	Title:                    "🔨 [YOUR-TEAMS]'s fines 🔨",
	UseRoles:                 true,
	UsePreviewPassword:       os.Getenv("PREVIEW_ENV") == "true",
	UseMatchEventTracker:     true,
	InjuryCounterTrackerName: "Injuries (Mr Glass)",
	UsePlayerOfTheDayName:    "Player of the Day",
	UseDudOfTheDayName:       "Dick of the Day",
}

var config = &Config{
	Title:                "🔨 Baileys Hammer 🔨",
	UseRoles:             false,
	UseMatchEventTracker: true,
	UsePreviewPassword:   os.Getenv("PREVIEW_ENV") == "true",
	/**

	If true, new fines will be approved by default (can be later decline by the finemaster)

	*/
	DefaultToApproved: true,

	ShowGoalScorerMatchList:  true,
	ShowGoalAssistMatchList:  true,
	InjuryCounterTrackerName: "🚑 Mr Glass 🚑",
	UsePlayerOfTheDayName:    "Player of the Day",
	UseDudOfTheDayName:       "🍆 Dick of the Day 🍆",
	ShowOpponentScore:        true,
}

// Use a custom type for keys to avoid conflicts in context values.
type contextKey string

const (
	titleKey                    contextKey = "Title"
	useRolesKey                 contextKey = "UseRoles"
	useMatchEventTrackerKey     contextKey = "UseMatchEventTracker"
	UsePlayerOfTheDayNameKey    contextKey = "UsePlayerOfTheDayName"
	ShowGoalScorerMatchListKey  contextKey = "ShowGoalScorerMatchList"
	ShowGoalAssistMatchListKey  contextKey = "ShowGoalAssistMatchList"
	UseDudOfTheDayNameKey       contextKey = "UseDudOfTheDayName"
	InjuryCounterTrackerNameKey contextKey = "InjuryCounterTrackerName"
	ShowOpponentScoreKey        contextKey = "ShowOpponentScoreKey"
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

func saveTeamToSession(r *http.Request, team Team, session *sessions.Session) {
	session.Values["TeamID"] = team.ID
	session.Values["team"] = team
	session.Values["TeamName"] = team.TeamName
	session.Values["ShowFineAddOnHomePage"] = team.ShowFineAddOnHomePage
	session.Values["ShowCourtTotals"] = team.ShowCourtTotals
	session.Save(r, nil)
}

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func getTeamId(ctx context.Context) uint {
	teamId, ok := ctx.Value("team_id").(uint)
	if !ok {
		// Handle the case where the title is not found or not a string
		teamId = 0
	}
	return teamId
}

func GetContext(r *http.Request, db *gorm.DB) context.Context {

	session, _ := store.Get(r, "session-name")

	// Get the selected team ID from the session
	teamID, teamExists := session.Values["team_id"].(uint)

	var team Team
	ctx := r.Context()

	// Check if team data exists in session; if not, fetch from the database
	if teamExists && teamID > 0 {
		if t, ok := session.Values["team"].(Team); ok {
			team = t
		} else {
			team, err := GetTeam(db, teamID)
			if err != nil {
				log.Printf("Error fetching team data %d from the database: %+v", teamID, err)
			} else {
				saveTeamToSession(r, *team, session)
			}

		}
	} else {
		// Query the database only if no team is selected and no session data exists

		teams, err := GetTeams(db, 1, 0)
		if err != nil {
			log.Printf("Error fetching team data from the database: %+v", err)
		}
		if len(teams) == 1 {
			team = teams[0]
			saveTeamToSession(r, team, session)
		} else {
			// Handle the case when no or multiple teams exist
			// You may want to redirect to a team selection page
		}
	}

	title := os.Getenv("TITLE")
	if title == "" {
		title = "🔨 Baileys Hammer 🔨"
	}
	ctx = context.WithValue(r.Context(), titleKey, title)

	if os.Getenv("DEV_ENV") == "true" {
		ctx = context.WithValue(ctx, useRolesKey, devConfig.UseRoles)
		ctx = context.WithValue(ctx, useMatchEventTrackerKey, devConfig.UseMatchEventTracker)
		ctx = context.WithValue(ctx, UsePlayerOfTheDayNameKey, devConfig.UsePlayerOfTheDayName)
		ctx = context.WithValue(ctx, UseDudOfTheDayNameKey, devConfig.UseDudOfTheDayName)
		ctx = context.WithValue(ctx, InjuryCounterTrackerNameKey, devConfig.InjuryCounterTrackerName)
		ctx = context.WithValue(ctx, ShowGoalScorerMatchListKey, devConfig.ShowGoalScorerMatchList)
		ctx = context.WithValue(ctx, ShowGoalAssistMatchListKey, devConfig.ShowGoalAssistMatchList)
		ctx = context.WithValue(ctx, ShowOpponentScoreKey, devConfig.ShowOpponentScore)
	} else {
		ctx := context.WithValue(r.Context(), titleKey, config.Title)
		ctx = context.WithValue(ctx, useRolesKey, config.UseRoles)
		ctx = context.WithValue(ctx, useMatchEventTrackerKey, config.UseMatchEventTracker)
		ctx = context.WithValue(ctx, UsePlayerOfTheDayNameKey, config.UsePlayerOfTheDayName)
		ctx = context.WithValue(ctx, UseDudOfTheDayNameKey, config.UseDudOfTheDayName)
		ctx = context.WithValue(ctx, InjuryCounterTrackerNameKey, config.InjuryCounterTrackerName)
		ctx = context.WithValue(ctx, ShowGoalScorerMatchListKey, config.ShowGoalScorerMatchList)
		ctx = context.WithValue(ctx, ShowGoalAssistMatchListKey, config.ShowGoalAssistMatchList)
		ctx = context.WithValue(ctx, ShowOpponentScoreKey, config.ShowOpponentScore)
	}

	return ctx
}
