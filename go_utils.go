package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

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
	teamKey                     contextKey = "Team"
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

func FooterTeam(ctx context.Context) Team {
	if team, ok := ctx.Value(teamKey).(Team); ok {
		return team
	}
	return Team{
		EnableFinesModule:        true,
		EnablePaymentsModule:     true,
		ShowPitchMatchOnHomePage: true,
		ShowCourtSheetOnHomePage: true,
		EnablePublicFeedbackForm: true,
		EnableLeaderboardModule:  true,
		AllowAdminRegistration:   true,
		LineupPlayerCount:        defaultLineupPlayerCount,
	}
}

func saveTeamToSession(r *http.Request, team Team, session *sessions.Session) {
	setTeamSessionValues(team, session)
	session.Save(r, nil)
}

func setTeamSessionValues(team Team, session *sessions.Session) {
	team.LineupPlayerCount = teamLineupPlayerCount(team)
	session.Values["TeamID"] = team.ID
	// Primary key used throughout the codebase (and expected by GetContext).
	session.Values["team_id"] = team.ID
	session.Values["team"] = team
	session.Values["TeamName"] = team.TeamName
	session.Values["ShowFineAddOnHomePage"] = team.ShowFineAddOnHomePage
	session.Values["ShowPitchMatchOnHomePage"] = team.ShowPitchMatchOnHomePage
	session.Values["ShowCourtSheetOnHomePage"] = team.ShowCourtSheetOnHomePage
	session.Values["EnablePublicFeedbackForm"] = team.EnablePublicFeedbackForm
	session.Values["ShowCourtTotals"] = team.ShowCourtTotals
	session.Values["EnableFinesModule"] = team.EnableFinesModule
	session.Values["EnableLineupsModule"] = team.EnableLineupsModule
	session.Values["EnableMatchesModule"] = team.EnableMatchesModule
	session.Values["EnablePlayersModule"] = team.EnablePlayersModule
	session.Values["EnablePaymentsModule"] = team.EnablePaymentsModule
	session.Values["EnableCourtModule"] = team.EnableCourtModule
	session.Values["EnableLeaderboardModule"] = team.EnableLeaderboardModule
	session.Values["AllowAdminRegistration"] = team.AllowAdminRegistration
	session.Values["LineupPlayerCount"] = teamLineupPlayerCount(team)
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

	ctx := r.Context()

	// Determine the active team for this request.
	//
	// Priority:
	// - google admin-session cookie
	// - admin-token cookie (set by password /admin)
	// - session (team_id / TeamID)
	// - if only one team exists, auto-select it
	var teamID uint
	if googleAuthEnabled() {
		teamID = selectedAdminTeamID(r, db)
	} else {
		if c, err := r.Cookie(adminTokenCookieName); err == nil {
			if tid, ok := parseAdminToken(c.Value, adminTokenSecret()); ok && tid > 0 {
				teamID = tid
			}
		}
	}
	// Backward compatibility with the older finemaster auth flow.
	if teamID == 0 && !googleAuthEnabled() {
		if c, err := r.Cookie("admin-user"); err == nil {
			if tid64, err := strconv.ParseUint(c.Value, 10, 64); err == nil && tid64 > 0 {
				teamID = uint(tid64)
			}
		}
	}
	if teamID == 0 {
		if tid, ok := session.Values["team_id"].(uint); ok && tid > 0 {
			teamID = tid
		} else if tid, ok := session.Values["TeamID"].(uint); ok && tid > 0 {
			teamID = tid
		}
	}

	// Resolve the team record and persist into session for later requests.
	var activeTeam *Team
	if teamID > 0 {
		if team, err := GetTeam(db, teamID); err != nil {
			log.Printf("Error fetching team data %d from the database: %+v", teamID, err)
			teamID = 0
		} else {
			activeTeam = team
			saveTeamToSession(r, *team, session)
		}
	}
	if teamID == 0 {
		teams, err := GetTeams(db, 1, 0)
		if err != nil {
			log.Printf("Error fetching team data from the database: %+v", err)
		} else if len(teams) == 1 {
			teamID = teams[0].ID
			activeTeam = &teams[0]
			saveTeamToSession(r, teams[0], session)
		}
	}

	if teamID > 0 {
		ctx = context.WithValue(ctx, "team_id", teamID)
	}
	if activeTeam != nil {
		ctx = context.WithValue(ctx, teamKey, *activeTeam)
	}

	title := os.Getenv("TITLE")
	if title == "" {
		title = "🔨 Baileys Hammer 🔨"
	}
	ctx = context.WithValue(ctx, titleKey, title)

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
		ctx = context.WithValue(ctx, titleKey, config.Title)
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
