package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const (
	lineupUserCookieName     = "lineup-user"
	matchDayLineupCookieName = "match-day-lineup"

	lineupStatusDraft     = "draft"
	lineupStatusSubmitted = "submitted"
	lineupStatusProposed  = "proposed"
	lineupStatusSelected  = "selected"
	lineupStatusLive      = "live"
	lineupStatusUsed      = "used"

	formationStatusDraft     = "draft"
	formationStatusSubmitted = "submitted"
	formationStatusLive      = "live"

	sharedLineupLiveFormationError = "shared line-ups can only use live formations"

	matchDayEventStarted        = "start-match"
	matchDayEventInit           = "init"
	matchDayEventHalfTime       = "half-time"
	matchDayEventFinished       = "match-finished"
	matchDayEventExtraTimeStart = "match-extra-time-start"
	matchDayEventResumed        = "match-resumed"
	matchDayRegularHalfMinute   = 45
	matchDayRegularEndMinute    = 90
	matchDayExtraHalfMinute     = 105
	matchDayExtraEndMinute      = 120
	matchDayOverrunWarnMinute   = 135
)

var lineupStatuses = []string{lineupStatusDraft, lineupStatusSubmitted, lineupStatusProposed, lineupStatusSelected, lineupStatusLive, lineupStatusUsed}
var formationStatuses = []string{formationStatusDraft, formationStatusSubmitted, formationStatusLive}

type LineupActor struct {
	TeamID  uint
	UserID  uint
	IsAdmin bool
}

type LineupPitchPosition struct {
	Position FormationPosition
	Players  []LineupPlayer
}

type MatchDayPlayerTime struct {
	Player              Player
	Minutes             int
	LiveMinutes         int
	Current             bool
	CurrentIndexNumber  int
	CurrentPositionName string
	Positions           []string
}

type MatchDaySuggestedSub struct {
	PositionIndex     int
	PositionName      string
	DueMinute         int
	CurrentPlayer     Player
	ReplacementPlayer Player
}

type MatchDayAdjacentLineups struct {
	Previous *Lineup
	Next     *Lineup
}

type LineupProblem struct {
	Level string
	Msg   string
}

type matchDayCurrentAction struct {
	Minute    int
	CreatedAt time.Time
	Kind      string
	Index     int
	PlayerID  uint
	OtherID   uint
}

type matchDayClockState struct {
	Started       bool
	Paused        bool
	Finished      bool
	InExtraTime   bool
	BaseTime      *time.Time
	BaseMinute    int
	CurrentMinute int
	HalfTimeCount int
	PauseEvent    string
}

func makeLineupUserToken(teamID uint, userID uint, secret string) string {
	body := fmt.Sprintf("%d:%d", teamID, userID)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	return body + ":" + hex.EncodeToString(mac.Sum(nil))
}

func parseLineupUserToken(token string, secret string) (uint, uint, bool) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return 0, 0, false
	}
	body := parts[0] + ":" + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return 0, 0, false
	}
	teamID64, teamErr := strconv.ParseUint(parts[0], 10, 64)
	userID64, userErr := strconv.ParseUint(parts[1], 10, 64)
	if teamErr != nil || userErr != nil || teamID64 == 0 || userID64 == 0 {
		return 0, 0, false
	}
	return uint(teamID64), uint(userID64), true
}

func getLineupActor(r *http.Request, db *gorm.DB) LineupActor {
	ctx := GetContext(r, db)
	actor := LineupActor{TeamID: getTeamId(ctx)}
	if c, err := r.Cookie(adminTokenCookieName); err == nil {
		if tid, ok := parseAdminToken(c.Value, adminTokenSecret()); ok && tid > 0 {
			actor.TeamID = tid
			actor.IsAdmin = true
		}
	}
	if c, err := r.Cookie(lineupUserCookieName); err == nil {
		if tid, uid, ok := parseLineupUserToken(c.Value, adminTokenSecret()); ok {
			if actor.TeamID == 0 {
				actor.TeamID = tid
			}
			if actor.TeamID == tid {
				actor.UserID = uid
			}
		}
	}
	return actor
}

func actorCanEditLineup(actor LineupActor, lineup Lineup) bool {
	if actor.TeamID != lineup.TeamID {
		return false
	}
	if lineup.Status == lineupStatusUsed {
		return false
	}
	return actor.IsAdmin || actor.UserID > 0 && actor.UserID == lineup.CreatedByUserID && lineup.Status == lineupStatusDraft
}

func actorCanCopyLineup(actor LineupActor, lineup Lineup) bool {
	if actor.TeamID != lineup.TeamID {
		return false
	}
	if lineup.Status == lineupStatusDraft {
		return actorCanSeeDraft(actor, lineup.TeamID, lineup.CreatedByUserID)
	}
	return true
}

func actorCanSeeDraft(actor LineupActor, teamID uint, createdByUserID uint) bool {
	return actor.IsAdmin && actor.TeamID == teamID || actor.UserID > 0 && actor.UserID == createdByUserID && actor.TeamID == teamID
}

func actorCanEditFormation(actor LineupActor, formation Formation) bool {
	if actor.TeamID != formation.TeamID {
		return false
	}
	return actor.IsAdmin || actor.UserID > 0 && actor.UserID == formation.CreatedByUserID && formation.Status == formationStatusDraft
}

func lineupStatusesForActor(actor LineupActor) []string {
	if actor.IsAdmin {
		return lineupStatuses
	}
	return []string{lineupStatusDraft, lineupStatusSubmitted}
}

func formationStatusesForActor(actor LineupActor) []string {
	if actor.IsAdmin {
		return formationStatuses
	}
	return []string{formationStatusDraft, formationStatusSubmitted}
}

func isValidLineupStatus(status string) bool {
	for _, s := range lineupStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func isValidFormationStatus(status string) bool {
	for _, s := range formationStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func getTeamByKeyAndMemberPassword(db *gorm.DB, teamKey string, memberPassword string) (*Team, error) {
	var team Team
	result := db.Where("team_key = ? AND team_member_pass = ?", teamKey, memberPassword).First(&team)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func findOrCreateLineupUser(db *gorm.DB, teamID uint, displayName string) (*LineupUser, error) {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return nil, errors.New("display name is required")
	}
	var user LineupUser
	err := db.Where("team_id = ? AND display_name = ?", teamID, displayName).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	user = LineupUser{TeamID: teamID, DisplayName: displayName}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func lineupLoginPageForRequest(r *http.Request, db *gorm.DB, msg string) templ.Component {
	return lineupLoginPage(msg, strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")), getTeamId(GetContext(r, db)))
}

func setLineupUserCookie(w http.ResponseWriter, teamID uint, userID uint) {
	http.SetCookie(w, &http.Cookie{
		Name:     lineupUserCookieName,
		Value:    makeLineupUserToken(teamID, userID, adminTokenSecret()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 365,
	})
}

func loginLineupUserFromGoogle(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if !googleAuthEnabled() {
		lineupLoginPageForRequest(r, db, "Google line-up access is not configured.").Render(GetContext(r, db), w)
		return
	}
	info, err := googleCredentialVerifier(r.FormValue("credential"))
	if err != nil {
		lineupLoginPageForRequest(r, db, "Google sign-in failed.").Render(GetContext(r, db), w)
		return
	}
	teamID := parseAdminTeamID(r, db)
	if teamID == 0 {
		lineupLoginPageForRequest(r, db, "Choose a team before requesting line-up access.").Render(GetContext(r, db), w)
		return
	}
	user, _, err := SaveAdminUserFromGoogle(db, *info, true)
	if err != nil {
		lineupLoginPageForRequest(r, db, "Could not save Google user.").Render(GetContext(r, db), w)
		return
	}
	if canGoogleUserAccessLineups(db, user.ID, teamID) {
		displayName := strings.TrimSpace(user.DisplayName)
		if displayName == "" {
			displayName = user.Email
		}
		lineupUser, err := findOrCreateLineupUser(db, teamID, displayName)
		if err != nil {
			lineupLoginPageForRequest(r, db, err.Error()).Render(GetContext(r, db), w)
			return
		}
		setLineupUserCookie(w, teamID, lineupUser.ID)
		w.Header().Set("HX-Redirect", "/lineups")
		http.Redirect(w, r, "/lineups", http.StatusSeeOther)
		return
	}
	if err := RequestAdminAccess(db, user.ID, teamID, adminRoleLineupAccess); err != nil {
		lineupLoginPageForRequest(r, db, "Could not request line-up access.").Render(GetContext(r, db), w)
		return
	}
	lineupLoginPageForRequest(r, db, "Line-up access request sent. A team admin or super admin can approve it.").Render(GetContext(r, db), w)
}

func lineupActorName(db *gorm.DB, actor LineupActor) string {
	if actor.UserID == 0 {
		if actor.IsAdmin {
			return "admin"
		}
		return "user"
	}
	var user LineupUser
	if err := db.First(&user, actor.UserID).Error; err != nil || strings.TrimSpace(user.DisplayName) == "" {
		return "user"
	}
	return user.DisplayName
}

func lineupLockActorName(r *http.Request, db *gorm.DB, actor LineupActor) string {
	if actor.IsAdmin {
		if adminUser, _, ok := currentAdminUser(r, db); ok {
			if strings.TrimSpace(adminUser.DisplayName) != "" {
				return adminUser.DisplayName
			}
			if strings.TrimSpace(adminUser.Email) != "" {
				return adminUser.Email
			}
		}
	}
	return lineupActorName(db, actor)
}

func copyName(baseName string, actorName string) string {
	baseName = strings.TrimSpace(baseName)
	if baseName == "" {
		baseName = "Untitled"
	}
	actorName = strings.TrimSpace(actorName)
	if actorName == "" {
		actorName = "user"
	}
	return fmt.Sprintf("%s - [%s - copy]", baseName, actorName)
}

func visibleLineups(db *gorm.DB, actor LineupActor) ([]Lineup, error) {
	if actor.TeamID == 0 {
		return []Lineup{}, nil
	}
	q := db.Preload("Formation").Preload("Match").Preload("Creator").Where("team_id = ?", actor.TeamID)
	if !actor.IsAdmin {
		q = q.Where("status != ? OR created_by_user_id = ?", lineupStatusDraft, actor.UserID)
	}
	var lineups []Lineup
	if err := q.Order("updated_at DESC").Find(&lineups).Error; err != nil {
		return nil, err
	}
	return lineups, nil
}

func liveLineupsForTeam(db *gorm.DB, teamID uint) ([]Lineup, error) {
	if teamID == 0 {
		return []Lineup{}, nil
	}
	var lineups []Lineup
	if err := db.Preload("Formation").Where("team_id = ? AND status = ?", teamID, lineupStatusLive).Order("updated_at DESC").Find(&lineups).Error; err != nil {
		return nil, err
	}
	return lineups, nil
}

func visibleFormations(db *gorm.DB, actor LineupActor) ([]Formation, error) {
	if actor.TeamID == 0 {
		return []Formation{}, nil
	}
	q := db.Preload("Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC")
	}).Where("team_id = ?", actor.TeamID)
	if !actor.IsAdmin {
		q = q.Where("status != ? OR created_by_user_id = ?", formationStatusDraft, actor.UserID)
	}
	var formations []Formation
	if err := q.Order("updated_at DESC").Find(&formations).Error; err != nil {
		return nil, err
	}
	return formations, nil
}

func eligibleFormationsForLineup(db *gorm.DB, actor LineupActor, lineupStatus string) ([]Formation, error) {
	if actor.TeamID == 0 {
		return []Formation{}, nil
	}
	q := db.Preload("Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC")
	}).Where("team_id = ?", actor.TeamID)
	if lineupStatus == lineupStatusDraft {
		if actor.IsAdmin {
			q = q.Where("status IN ?", []string{formationStatusLive, formationStatusSubmitted, formationStatusDraft})
		} else {
			q = q.Where("status = ? OR status = ? OR (status = ? AND created_by_user_id = ?)", formationStatusLive, formationStatusSubmitted, formationStatusDraft, actor.UserID)
		}
	} else {
		q = q.Where("status = ?", formationStatusLive)
	}
	var formations []Formation
	if err := q.Order("status ASC, name ASC").Find(&formations).Error; err != nil {
		return nil, err
	}
	return formations, nil
}

func getLineup(db *gorm.DB, id uint) (*Lineup, error) {
	var lineup Lineup
	err := db.Preload("Players", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC, slot_order ASC")
	}).Preload("Players.Player").Preload("Formation.Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC")
	}).Preload("Match").Preload("Match.Events").Preload("Creator").First(&lineup, id).Error
	if err != nil {
		return nil, err
	}
	return &lineup, nil
}

func lineupMatchID(lineup *Lineup) uint {
	if lineup == nil {
		return 0
	}
	return lineup.Match.ID
}

func getFormation(db *gorm.DB, id uint) (*Formation, error) {
	var formation Formation
	err := db.Preload("Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC")
	}).First(&formation, id).Error
	if err != nil {
		return nil, err
	}
	return &formation, nil
}

func lineupPitchPlayers(lineup *Lineup) []LineupPitchPosition {
	playersByIndex := map[int][]LineupPlayer{}
	for i := range lineup.Players {
		lp := lineup.Players[i]
		if lp.PlayerID > 0 {
			playersByIndex[lp.IndexNumber] = append(playersByIndex[lp.IndexNumber], lp)
		}
	}
	out := []LineupPitchPosition{}
	for _, pos := range lineup.Formation.Positions {
		out = append(out, LineupPitchPosition{
			Position: pos,
			Players:  playersByIndex[pos.IndexNumber],
		})
	}
	return out
}

func clampMatchMinute(minute int) int {
	if minute < 0 {
		return 0
	}
	if minute > matchDayExtraEndMinute {
		return matchDayExtraEndMinute
	}
	return minute
}

func appendPlayerPosition(positions []string, position string) []string {
	for _, existing := range positions {
		if existing == position {
			return positions
		}
	}
	return append(positions, position)
}

func matchDayPlayerTimes(lineup *Lineup) []MatchDayPlayerTime {
	if lineup == nil {
		return []MatchDayPlayerTime{}
	}
	currentMinute := matchDayMinute(lineup.Match)
	playerTimes := map[uint]*MatchDayPlayerTime{}
	playerOrder := []uint{}
	for _, pp := range lineupPitchPlayers(lineup) {
		assignedPlayers := pp.Players
		sort.Slice(assignedPlayers, func(i, j int) bool {
			return assignedPlayers[i].SlotOrder < assignedPlayers[j].SlotOrder
		})
		for i, assigned := range assignedPlayers {
			if assigned.PlayerID == 0 {
				continue
			}
			startMinute := 0
			if assigned.SlotOrder > 0 {
				startMinute = clampMatchMinute(assigned.SubMinute)
			}
			endMinute := 90
			if i+1 < len(assignedPlayers) {
				endMinute = clampMatchMinute(assignedPlayers[i+1].SubMinute)
			}
			if endMinute < startMinute {
				endMinute = startMinute
			}
			playerTime, ok := playerTimes[assigned.PlayerID]
			if !ok {
				playerTimes[assigned.PlayerID] = &MatchDayPlayerTime{Player: assigned.Player}
				playerTime = playerTimes[assigned.PlayerID]
				playerOrder = append(playerOrder, assigned.PlayerID)
			}
			playerTime.Minutes += endMinute - startMinute
			liveEndMinute := currentMinute
			if liveEndMinute > endMinute {
				liveEndMinute = endMinute
			}
			if liveEndMinute > startMinute {
				playerTime.LiveMinutes += liveEndMinute - startMinute
			}
			if currentMinute >= startMinute && currentMinute < endMinute || currentMinute >= 90 && i == len(assignedPlayers)-1 {
				playerTime.Current = true
				playerTime.CurrentIndexNumber = pp.Position.IndexNumber
				playerTime.CurrentPositionName = pp.Position.PositionName
			}
			playerTime.Positions = appendPlayerPosition(playerTime.Positions, pp.Position.PositionName)
		}
	}
	out := []MatchDayPlayerTime{}
	for _, playerID := range playerOrder {
		out = append(out, *playerTimes[playerID])
	}
	return out
}

func lineupProblems(db *gorm.DB, lineup *Lineup) ([]LineupProblem, error) {
	problems := []LineupProblem{}
	if lineup == nil {
		return problems, nil
	}
	playerCount := lineupPlayerCountForTeam(db, lineup.TeamID)
	if lineup.FormationID == 0 || len(lineup.Formation.Positions) == 0 {
		problems = append(problems, LineupProblem{Level: "error", Msg: "no formation selected"})
	} else if len(lineup.Formation.Positions) != playerCount {
		problems = append(problems, LineupProblem{Level: "error", Msg: fmt.Sprintf("formation has %d positions; team is configured for %s", len(lineup.Formation.Positions), sideLabel(playerCount))})
	}

	hasStartingGK := false
	for _, position := range lineup.Formation.Positions {
		name := strings.ToLower(strings.TrimSpace(position.PositionName))
		if name != "gk" && !strings.Contains(name, "goalkeeper") && !strings.Contains(name, "goal keeper") {
			continue
		}
		for _, player := range lineup.Players {
			if player.IndexNumber == position.IndexNumber && player.SlotOrder == 0 && player.PlayerID > 0 {
				hasStartingGK = true
				break
			}
		}
		if hasStartingGK {
			break
		}
	}
	if !hasStartingGK {
		problems = append(problems, LineupProblem{Level: "error", Msg: "no starting GK"})
	}
	problems = append(problems, missingStartingPositionProblems(lineup)...)
	problems = append(problems, duplicateStartingPlayerProblems(lineup)...)

	matchID := lineupMatchID(lineup)
	if matchID == 0 {
		problems = append(problems, LineupProblem{Level: "warn", Msg: "no match selected for availability checks"})
		return problems, nil
	}

	players, err := GetPlayers(db, 0, 999)
	if err != nil {
		return problems, err
	}
	unavailablePlayerIDs, err := GetUnavailablePlayerIDsForMatch(db, lineup.TeamID, matchID)
	if err != nil {
		return problems, err
	}
	availablePlayerCount, _ := playerAvailabilityCounts(players, unavailablePlayerIDs)
	if availablePlayerCount < playerCount {
		problems = append(problems, LineupProblem{Level: "error", Msg: fmt.Sprintf("less than %d players available", playerCount)})
	} else if availablePlayerCount < playerCount+4 {
		problems = append(problems, LineupProblem{Level: "warn", Msg: fmt.Sprintf("less than %d players available", playerCount+4)})
	}

	return problems, nil
}

func lineupPlayerCountForTeam(db *gorm.DB, teamID uint) int {
	team, err := GetTeam(db, teamID)
	if err != nil {
		return defaultLineupPlayerCount
	}
	return teamLineupPlayerCount(*team)
}

func missingStartingPositionProblems(lineup *Lineup) []LineupProblem {
	if lineup == nil {
		return nil
	}
	startersByIndex := map[int]bool{}
	for _, player := range lineup.Players {
		if player.PlayerID > 0 && player.SlotOrder == 0 {
			startersByIndex[player.IndexNumber] = true
		}
	}
	missingPositions := []string{}
	for _, position := range lineup.Formation.Positions {
		if startersByIndex[position.IndexNumber] {
			continue
		}
		positionName := strings.TrimSpace(position.PositionName)
		if positionName == "" {
			positionName = fmt.Sprintf("position %d", position.IndexNumber)
		}
		name := strings.ToLower(positionName)
		if name == "gk" || strings.Contains(name, "goalkeeper") || strings.Contains(name, "goal keeper") {
			continue
		}
		missingPositions = append(missingPositions, positionName)
	}
	if len(missingPositions) == 0 {
		return nil
	}
	return []LineupProblem{{
		Level: "warn",
		Msg:   fmt.Sprintf("missing starters: %s", strings.Join(missingPositions, ", ")),
	}}
}

func duplicateStartingPlayerProblems(lineup *Lineup) []LineupProblem {
	if lineup == nil {
		return nil
	}
	positionNames := map[int]string{}
	for _, position := range lineup.Formation.Positions {
		name := strings.TrimSpace(position.PositionName)
		if name == "" {
			name = fmt.Sprintf("position %d", position.IndexNumber)
		}
		positionNames[position.IndexNumber] = name
	}
	startersByPlayerID := map[uint][]LineupPlayer{}
	for _, player := range lineup.Players {
		if player.PlayerID == 0 || player.SlotOrder != 0 {
			continue
		}
		startersByPlayerID[player.PlayerID] = append(startersByPlayerID[player.PlayerID], player)
	}
	problems := []LineupProblem{}
	for _, starters := range startersByPlayerID {
		if len(starters) < 2 {
			continue
		}
		playerName := strings.TrimSpace(starters[0].Player.Name)
		if playerName == "" {
			playerName = fmt.Sprintf("player %d", starters[0].PlayerID)
		}
		positions := []string{}
		for _, starter := range starters {
			positionName := positionNames[starter.IndexNumber]
			if positionName == "" {
				positionName = fmt.Sprintf("position %d", starter.IndexNumber)
			}
			positions = append(positions, positionName)
		}
		problems = append(problems, LineupProblem{
			Level: "error",
			Msg:   fmt.Sprintf("duplicate starting player: %s is starting at %s", playerName, strings.Join(positions, ", ")),
		})
	}
	return problems
}

func currentMatchDayPlayerCount(playerTimes []MatchDayPlayerTime) int {
	count := 0
	for _, playerTime := range playerTimes {
		if playerTime.Current {
			count++
		}
	}
	return count
}

func matchDayPlayerCurrentlyOn(playerTimes []MatchDayPlayerTime, playerID uint) bool {
	for _, playerTime := range playerTimes {
		if playerTime.Player.ID == playerID && playerTime.Current {
			return true
		}
	}
	return false
}

func matchDayPlayerActuallyCurrentlyOn(lineup *Lineup, playerID uint) bool {
	for _, currentPlayerID := range matchDayActualCurrentPositions(lineup) {
		if currentPlayerID == playerID {
			return true
		}
	}
	return false
}

func currentMatchDayPlayerForPosition(playerTimes []MatchDayPlayerTime, indexNumber int) MatchDayPlayerTime {
	for _, playerTime := range playerTimes {
		if playerTime.Current && playerTime.CurrentIndexNumber == indexNumber {
			return playerTime
		}
	}
	return MatchDayPlayerTime{}
}

func matchDayPositionNamesByIndex(lineup *Lineup) map[int]string {
	names := map[int]string{}
	if lineup == nil {
		return names
	}
	for _, position := range lineup.Formation.Positions {
		names[position.IndexNumber] = position.PositionName
	}
	return names
}

func matchDayPlayerByID(lineup *Lineup) map[uint]Player {
	players := map[uint]Player{}
	if lineup == nil {
		return players
	}
	for _, lineupPlayer := range lineup.Players {
		if lineupPlayer.PlayerID > 0 {
			players[lineupPlayer.PlayerID] = lineupPlayer.Player
		}
	}
	return players
}

func matchDayCurrentActions(lineup *Lineup, minute int) []matchDayCurrentAction {
	return matchDayCurrentActionsFor(lineup, minute, true)
}

func matchDayActualCurrentActions(lineup *Lineup, minute int) []matchDayCurrentAction {
	return matchDayCurrentActionsFor(lineup, minute, false)
}

func matchDayCurrentActionsFor(lineup *Lineup, minute int, includePlannedSubs bool) []matchDayCurrentAction {
	if lineup == nil {
		return nil
	}
	actions := []matchDayCurrentAction{}
	if includePlannedSubs {
		for _, lineupPlayer := range lineup.Players {
			if lineupPlayer.PlayerID == 0 || lineupPlayer.SlotOrder == 0 {
				continue
			}
			subMinute := clampMatchMinute(lineupPlayer.SubMinute)
			if subMinute > minute {
				continue
			}
			actions = append(actions, matchDayCurrentAction{
				Minute:    subMinute,
				CreatedAt: lineupPlayer.CreatedAt,
				Kind:      "sub",
				Index:     lineupPlayer.IndexNumber,
				PlayerID:  lineupPlayer.PlayerID,
			})
		}
	}
	swapEvents := []MatchEvent{}
	subOffByMinute := map[int][]MatchEvent{}
	subOnByMinute := map[int][]MatchEvent{}
	for _, event := range lineup.Match.Events {
		if event.EventType == "swap" && event.PlayerId > 0 && clampMatchMinute(event.EventMinute) <= minute {
			swapEvents = append(swapEvents, event)
		}
		switch event.EventType {
		case "subbed-off":
			if event.PlayerId > 0 && clampMatchMinute(event.EventMinute) <= minute {
				event.EventMinute = clampMatchMinute(event.EventMinute)
				subOffByMinute[event.EventMinute] = append(subOffByMinute[event.EventMinute], event)
			}
		case "subbed-on":
			if event.PlayerId > 0 && clampMatchMinute(event.EventMinute) <= minute {
				event.EventMinute = clampMatchMinute(event.EventMinute)
				subOnByMinute[event.EventMinute] = append(subOnByMinute[event.EventMinute], event)
			}
		}
	}
	for subMinute, offEvents := range subOffByMinute {
		onEvents := subOnByMinute[subMinute]
		sort.SliceStable(offEvents, func(i, j int) bool {
			return offEvents[i].CreatedAt.Before(offEvents[j].CreatedAt)
		})
		sort.SliceStable(onEvents, func(i, j int) bool {
			return onEvents[i].CreatedAt.Before(onEvents[j].CreatedAt)
		})
		for i := 0; i < len(offEvents) && i < len(onEvents); i++ {
			index := -1
			for _, lineupPlayer := range lineup.Players {
				if lineupPlayer.PlayerID == offEvents[i].PlayerId && lineupPlayer.SlotOrder == 0 {
					index = lineupPlayer.IndexNumber
					break
				}
			}
			if index < 0 {
				for _, action := range actions {
					if action.Kind == "sub" && action.PlayerID == offEvents[i].PlayerId {
						index = action.Index
					}
				}
			}
			if index >= 0 {
				actions = append(actions, matchDayCurrentAction{
					Minute:    subMinute,
					CreatedAt: onEvents[i].CreatedAt,
					Kind:      "sub",
					Index:     index,
					PlayerID:  onEvents[i].PlayerId,
				})
			}
		}
	}
	sort.SliceStable(swapEvents, func(i, j int) bool {
		if swapEvents[i].EventMinute != swapEvents[j].EventMinute {
			return swapEvents[i].EventMinute < swapEvents[j].EventMinute
		}
		return swapEvents[i].CreatedAt.Before(swapEvents[j].CreatedAt)
	})
	for i := 0; i < len(swapEvents); {
		minute := swapEvents[i].EventMinute
		j := i
		for j < len(swapEvents) && swapEvents[j].EventMinute == minute {
			j++
		}
		for pairIndex := i; pairIndex+1 < j; pairIndex += 2 {
			actions = append(actions, matchDayCurrentAction{
				Minute:    clampMatchMinute(swapEvents[pairIndex].EventMinute),
				CreatedAt: swapEvents[pairIndex].CreatedAt,
				Kind:      "swap",
				PlayerID:  swapEvents[pairIndex].PlayerId,
				OtherID:   swapEvents[pairIndex+1].PlayerId,
			})
		}
		i = j
	}
	sort.SliceStable(actions, func(i, j int) bool {
		if actions[i].Minute != actions[j].Minute {
			return actions[i].Minute < actions[j].Minute
		}
		return actions[i].CreatedAt.Before(actions[j].CreatedAt)
	})
	return actions
}

func matchDayCurrentPositions(lineup *Lineup) map[int]uint {
	return matchDayCurrentPositionsFor(lineup, true)
}

func matchDayActualCurrentPositions(lineup *Lineup) map[int]uint {
	return matchDayCurrentPositionsFor(lineup, false)
}

func matchDayCurrentPositionsFor(lineup *Lineup, includePlannedSubs bool) map[int]uint {
	positions := map[int]uint{}
	if lineup == nil {
		return positions
	}
	for _, lineupPlayer := range lineup.Players {
		if lineupPlayer.PlayerID > 0 && lineupPlayer.SlotOrder == 0 {
			positions[lineupPlayer.IndexNumber] = lineupPlayer.PlayerID
		}
	}
	for _, action := range matchDayCurrentActionsFor(lineup, matchDayMinute(lineup.Match), includePlannedSubs) {
		switch action.Kind {
		case "sub":
			positions[action.Index] = action.PlayerID
		case "swap":
			playerIndex := -1
			otherIndex := -1
			for index, playerID := range positions {
				if playerID == action.PlayerID {
					playerIndex = index
				}
				if playerID == action.OtherID {
					otherIndex = index
				}
			}
			if playerIndex >= 0 && otherIndex >= 0 {
				positions[playerIndex], positions[otherIndex] = positions[otherIndex], positions[playerIndex]
			}
		}
	}
	return positions
}

func matchDaySuggestedSubs(lineup *Lineup, unavailablePlayerIDs map[uint]bool) []MatchDaySuggestedSub {
	if lineup == nil || lineup.Locked || lineupMatchID(lineup) == 0 || !matchDayCanLogLiveEvent(lineup.Match) {
		return nil
	}
	currentMinute := matchDayMinute(lineup.Match)
	actualPositions := matchDayActualCurrentPositions(lineup)
	actualPlayerOn := map[uint]bool{}
	for _, playerID := range actualPositions {
		if playerID > 0 {
			actualPlayerOn[playerID] = true
		}
	}
	playerByID := matchDayPlayerByID(lineup)
	positionNames := matchDayPositionNamesByIndex(lineup)
	suggestions := []MatchDaySuggestedSub{}
	for _, pitchPosition := range lineupPitchPlayers(lineup) {
		indexNumber := pitchPosition.Position.IndexNumber
		currentPlayerID := actualPositions[indexNumber]
		if currentPlayerID == 0 {
			continue
		}
		assignedPlayers := append([]LineupPlayer(nil), pitchPosition.Players...)
		sort.SliceStable(assignedPlayers, func(i, j int) bool {
			return assignedPlayers[i].SlotOrder < assignedPlayers[j].SlotOrder
		})
		for i := 1; i < len(assignedPlayers); i++ {
			replacement := assignedPlayers[i]
			if replacement.PlayerID == 0 || replacement.SubMinute <= 0 || clampMatchMinute(replacement.SubMinute) > currentMinute {
				continue
			}
			if playerUnavailableForMatch(unavailablePlayerIDs, replacement.PlayerID) || actualPlayerOn[replacement.PlayerID] {
				continue
			}
			if assignedPlayers[i-1].PlayerID != currentPlayerID {
				continue
			}
			positionName := positionNames[indexNumber]
			if strings.TrimSpace(positionName) == "" {
				positionName = fmt.Sprintf("Position %d", indexNumber)
			}
			currentPlayer := playerByID[currentPlayerID]
			if currentPlayer.ID == 0 {
				currentPlayer = Player{Model: gorm.Model{ID: currentPlayerID}, Name: fmt.Sprintf("player %d", currentPlayerID)}
			}
			replacementPlayer := playerByID[replacement.PlayerID]
			if replacementPlayer.ID == 0 {
				replacementPlayer = replacement.Player
			}
			if replacementPlayer.ID == 0 {
				replacementPlayer = Player{Model: gorm.Model{ID: replacement.PlayerID}, Name: fmt.Sprintf("player %d", replacement.PlayerID)}
			}
			suggestions = append(suggestions, MatchDaySuggestedSub{
				PositionIndex:     indexNumber,
				PositionName:      positionName,
				DueMinute:         clampMatchMinute(replacement.SubMinute),
				CurrentPlayer:     currentPlayer,
				ReplacementPlayer: replacementPlayer,
			})
			break
		}
	}
	sort.SliceStable(suggestions, func(i, j int) bool {
		if suggestions[i].DueMinute != suggestions[j].DueMinute {
			return suggestions[i].DueMinute < suggestions[j].DueMinute
		}
		return suggestions[i].PositionIndex < suggestions[j].PositionIndex
	})
	return suggestions
}

func matchDayPlayerTimesWithEvents(lineup *Lineup) []MatchDayPlayerTime {
	playerTimes := matchDayPlayerTimes(lineup)
	if lineup == nil {
		return playerTimes
	}
	for i := range playerTimes {
		playerTimes[i].Current = false
		playerTimes[i].CurrentIndexNumber = 0
		playerTimes[i].CurrentPositionName = ""
	}
	playerTimeByID := map[uint]int{}
	for i := range playerTimes {
		playerTimeByID[playerTimes[i].Player.ID] = i
	}
	positionNames := matchDayPositionNamesByIndex(lineup)
	players := matchDayPlayerByID(lineup)
	currentMinute := matchDayMinute(lineup.Match)
	currentActions := matchDayCurrentActions(lineup, currentMinute)
	for index, playerID := range matchDayCurrentPositions(lineup) {
		if playerID == 0 {
			continue
		}
		i, ok := playerTimeByID[playerID]
		if !ok {
			player := players[playerID]
			if player.ID == 0 {
				player = Player{Model: gorm.Model{ID: playerID}, Name: fmt.Sprintf("player %d", playerID)}
			}
			playerTimes = append(playerTimes, MatchDayPlayerTime{Player: player})
			i = len(playerTimes) - 1
			playerTimeByID[playerID] = i
		}
		playerTimes[i].Current = true
		playerTimes[i].CurrentIndexNumber = index
		playerTimes[i].CurrentPositionName = positionNames[index]
		for _, action := range currentActions {
			if action.Kind == "sub" && action.PlayerID == playerID {
				playerTimes[i].LiveMinutes = currentMinute - action.Minute
				playerTimes[i].Minutes = matchDayRegularEndMinute - action.Minute
			}
		}
	}
	return playerTimes
}

func matchDayErrors(db *gorm.DB, lineup *Lineup, playerTimes []MatchDayPlayerTime, unavailablePlayerIDs map[uint]bool) ([]LineupProblem, error) {
	problems, err := lineupProblems(db, lineup)
	if err != nil {
		return problems, err
	}
	if lineup == nil {
		return problems, nil
	}
	if lineupMatchID(lineup) == 0 {
		filtered := []LineupProblem{}
		for _, problem := range problems {
			if problem.Msg == "no match selected for availability checks" {
				continue
			}
			filtered = append(filtered, problem)
		}
		problems = filtered
		problems = append(problems, LineupProblem{Level: "error", Msg: "no match selected"})
		return problems, nil
	}
	if matchDayStartTime(lineup.Match) == nil {
		problems = append(problems, LineupProblem{Level: "warn", Msg: "match has no start time set"})
	}

	currentPlayerCount := currentMatchDayPlayerCount(playerTimes)
	if currentPlayerCount < 11 {
		problems = append(problems, LineupProblem{Level: "error", Msg: fmt.Sprintf("only %d players on the pitch", currentPlayerCount)})
	} else if currentPlayerCount > 11 {
		problems = append(problems, LineupProblem{Level: "error", Msg: fmt.Sprintf("%d players on the pitch", currentPlayerCount)})
	}

	unavailableOnPitch := []string{}
	for _, playerTime := range playerTimes {
		if !playerTime.Current || !unavailablePlayerIDs[playerTime.Player.ID] {
			continue
		}
		playerName := strings.TrimSpace(playerTime.Player.Name)
		if playerName == "" {
			playerName = fmt.Sprintf("player %d", playerTime.Player.ID)
		}
		unavailableOnPitch = append(unavailableOnPitch, playerName)
	}
	if len(unavailableOnPitch) > 0 {
		problems = append(problems, LineupProblem{
			Level: "error",
			Msg:   fmt.Sprintf("unavailable players on the pitch: %s", strings.Join(unavailableOnPitch, ", ")),
		})
	}

	return problems, nil
}

func matchDayClockEventTime(event MatchEvent) time.Time {
	if event.EventTime != nil {
		return *event.EventTime
	}
	return event.CreatedAt
}

func matchDayClockEvents(match Match) []MatchEvent {
	events := []MatchEvent{}
	for _, event := range match.Events {
		switch event.EventType {
		case matchDayEventStarted, matchDayEventHalfTime, matchDayEventFinished, matchDayEventExtraTimeStart, matchDayEventResumed:
			events = append(events, event)
		}
	}
	sort.SliceStable(events, func(i, j int) bool {
		iTime := matchDayClockEventTime(events[i])
		jTime := matchDayClockEventTime(events[j])
		if !iTime.Equal(jTime) {
			return iTime.Before(jTime)
		}
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	return events
}

func matchDayStartTime(match Match) *time.Time {
	var startTime *time.Time
	for _, event := range matchDayClockEvents(match) {
		if event.EventType != matchDayEventStarted {
			continue
		}
		eventTime := matchDayClockEventTime(event)
		if startTime == nil || eventTime.Before(*startTime) {
			startTime = &eventTime
		}
	}
	if startTime != nil {
		return startTime
	}
	return match.StartTime
}

func matchDayClockStateForMatch(match Match) matchDayClockState {
	state := matchDayClockState{}
	startTime := matchDayStartTime(match)
	if startTime == nil {
		return state
	}
	state.Started = !time.Now().Before(*startTime)
	state.BaseTime = startTime
	for _, event := range matchDayClockEvents(match) {
		eventTime := matchDayClockEventTime(event)
		eventMinute := clampMatchMinute(event.EventMinute)
		switch event.EventType {
		case matchDayEventStarted:
			state.Started = !time.Now().Before(eventTime)
			state.Paused = false
			state.Finished = false
			state.BaseTime = &eventTime
			state.BaseMinute = 0
			state.PauseEvent = ""
		case matchDayEventHalfTime:
			state.Started = true
			state.Paused = true
			state.Finished = false
			state.BaseTime = &eventTime
			state.BaseMinute = eventMinute
			state.HalfTimeCount++
			state.PauseEvent = event.EventType
		case matchDayEventFinished:
			state.Started = true
			state.Paused = false
			state.Finished = true
			state.BaseTime = &eventTime
			state.BaseMinute = eventMinute
			state.PauseEvent = event.EventType
		case matchDayEventExtraTimeStart:
			state.Started = true
			state.Paused = false
			state.Finished = false
			state.InExtraTime = true
			state.BaseTime = &eventTime
			state.BaseMinute = eventMinute
			state.PauseEvent = ""
		case matchDayEventResumed:
			state.Started = true
			state.Paused = false
			state.Finished = false
			state.BaseTime = &eventTime
			state.BaseMinute = eventMinute
			state.PauseEvent = ""
		}
	}
	state.CurrentMinute = state.BaseMinute
	if state.Started && !state.Paused && !state.Finished && state.BaseTime != nil {
		state.CurrentMinute = state.BaseMinute + int(time.Since(*state.BaseTime).Minutes())
	}
	state.CurrentMinute = clampMatchMinute(state.CurrentMinute)
	return state
}

func matchDayMinute(match Match) int {
	return matchDayClockStateForMatch(match).CurrentMinute
}

func matchDayRawMinute(match Match) int {
	state := matchDayClockStateForMatch(match)
	if state.Started && !state.Paused && !state.Finished && state.BaseTime != nil {
		return state.BaseMinute + int(time.Since(*state.BaseTime).Minutes())
	}
	return state.CurrentMinute
}

func matchDayClockBaseDataValue(match Match) string {
	state := matchDayClockStateForMatch(match)
	if state.BaseTime == nil {
		return ""
	}
	return state.BaseTime.Format("2006-01-02T15:04:05")
}

func matchDayClockDataBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func matchDayCanShowHalfTime(match Match) bool {
	state := matchDayClockStateForMatch(match)
	if !state.Started || state.Paused || state.Finished {
		return false
	}
	if state.InExtraTime {
		return state.HalfTimeCount < 2 && state.CurrentMinute >= matchDayExtraHalfMinute
	}
	return state.HalfTimeCount == 0 && state.CurrentMinute >= matchDayRegularHalfMinute
}

func matchDayCanShowExtraTimeStart(match Match) bool {
	state := matchDayClockStateForMatch(match)
	return state.Started && !state.Paused && !state.Finished && !state.InExtraTime && state.CurrentMinute >= matchDayRegularEndMinute
}

func matchDayCanShowMatchFinished(match Match) bool {
	state := matchDayClockStateForMatch(match)
	if !state.Started || state.Paused || state.Finished {
		return false
	}
	if state.InExtraTime {
		return state.CurrentMinute >= matchDayExtraEndMinute
	}
	return state.CurrentMinute >= matchDayRegularEndMinute
}

func matchDayCanShowOverrunWarning(match Match) bool {
	state := matchDayClockStateForMatch(match)
	return state.Started && !state.Paused && !state.Finished && matchDayRawMinute(match) > matchDayOverrunWarnMinute
}

func matchDayCanShowResume(match Match) bool {
	state := matchDayClockStateForMatch(match)
	return state.Started && state.Paused && !state.Finished
}

func matchDayCanLogLiveEvent(match Match) bool {
	state := matchDayClockStateForMatch(match)
	return state.Started && !state.Paused && !state.Finished
}

func matchDayResumeMinute(match Match) int {
	return matchDayClockStateForMatch(match).BaseMinute
}

func matchDayClockStatusLabel(match Match) string {
	state := matchDayClockStateForMatch(match)
	switch {
	case state.Finished:
		return "Game ended"
	case state.Paused:
		return matchDayPauseReasonLabel(state.PauseEvent)
	case state.Started:
		return "Match running"
	default:
		return "Kick-off"
	}
}

func matchDayPauseReasonLabel(eventType string) string {
	switch eventType {
	case matchDayEventHalfTime:
		return "Paused for half-time"
	case matchDayEventFinished:
		return "Game ended"
	default:
		return "Game pause"
	}
}

func createMatchDayEventAt(db *gorm.DB, matchID uint, eventType string, playerID uint, eventMinute int, eventName string, eventTime time.Time) error {
	if strings.TrimSpace(eventName) == "" {
		eventName = eventType
	}
	return db.Create(&MatchEvent{
		MatchId:     uint64(matchID),
		EventName:   eventName,
		EventType:   eventType,
		EventTime:   &eventTime,
		EventMinute: clampMatchMinute(eventMinute),
		PlayerId:    playerID,
	}).Error
}

func createMatchDayEvent(db *gorm.DB, matchID uint, eventType string, playerID uint, eventMinute int, eventName string) error {
	return createMatchDayEventAt(db, matchID, eventType, playerID, eventMinute, eventName, time.Now())
}

func matchDayRedoEventValues(message string, event MatchEvent) url.Values {
	values := url.Values{}
	values.Set("msg", message)
	values.Set("redoEvent", "1")
	values.Set("eventName", event.EventName)
	values.Set("eventType", event.EventType)
	values.Set("playerId", fmt.Sprintf("%d", event.PlayerId))
	values.Set("minute", fmt.Sprintf("%d", event.EventMinute))
	if event.EventTime != nil {
		values.Set("eventTime", event.EventTime.Format(time.RFC3339Nano))
	}
	return values
}

func matchDayRedoEventFromQuery(values url.Values) *MatchEvent {
	if values.Get("redoEvent") != "1" {
		return nil
	}
	eventType := strings.TrimSpace(values.Get("eventType"))
	if eventType == "" {
		return nil
	}
	minute, err := strconv.Atoi(strings.TrimSpace(values.Get("minute")))
	if err != nil {
		minute = 0
	}
	playerID, _ := strconv.ParseUint(strings.TrimSpace(values.Get("playerId")), 10, 64)
	event := &MatchEvent{
		EventName:   strings.TrimSpace(values.Get("eventName")),
		EventType:   eventType,
		EventMinute: clampMatchMinute(minute),
		PlayerId:    uint(playerID),
	}
	if eventTimeValue := strings.TrimSpace(values.Get("eventTime")); eventTimeValue != "" {
		if eventTime, err := time.Parse(time.RFC3339Nano, eventTimeValue); err == nil {
			event.EventTime = &eventTime
		}
	}
	return event
}

func matchDayExtraPlayerEventLabel(eventType string) (string, bool) {
	labels := map[string]string{
		"yellow-card":       "yellow card",
		"red-card":          "red card",
		"gave-away-penalty": "gave away a pen",
	}
	label, ok := labels[eventType]
	return label, ok
}

type matchDayLineupSnapshot struct {
	LineupID      uint                           `json:"lineupId"`
	FormationID   uint                           `json:"formationId"`
	LineupName    string                         `json:"lineupName"`
	LineupDetails string                         `json:"lineupDetails"`
	FormationName string                         `json:"formationName"`
	Positions     []matchDaySnapshotPosition     `json:"positions"`
	Players       []matchDaySnapshotLineupPlayer `json:"players"`
}

type matchDaySnapshotPosition struct {
	IndexNumber  int     `json:"indexNumber"`
	PositionName string  `json:"positionName"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
}

type matchDaySnapshotLineupPlayer struct {
	IndexNumber int    `json:"indexNumber"`
	SlotOrder   int    `json:"slotOrder"`
	SubMinute   int    `json:"subMinute"`
	PlayerID    uint   `json:"playerId"`
	PlayerName  string `json:"playerName"`
}

func matchDaySnapshotFromLineup(lineup *Lineup) matchDayLineupSnapshot {
	snapshot := matchDayLineupSnapshot{
		LineupID:      lineup.ID,
		FormationID:   lineup.FormationID,
		LineupName:    lineup.Name,
		LineupDetails: lineup.Details,
		FormationName: lineup.Formation.Name,
	}
	for _, position := range lineup.Formation.Positions {
		snapshot.Positions = append(snapshot.Positions, matchDaySnapshotPosition{
			IndexNumber:  position.IndexNumber,
			PositionName: position.PositionName,
			X:            position.X,
			Y:            position.Y,
		})
	}
	for _, player := range lineup.Players {
		if player.PlayerID == 0 {
			continue
		}
		snapshot.Players = append(snapshot.Players, matchDaySnapshotLineupPlayer{
			IndexNumber: player.IndexNumber,
			SlotOrder:   player.SlotOrder,
			SubMinute:   player.SubMinute,
			PlayerID:    player.PlayerID,
			PlayerName:  player.Player.Name,
		})
	}
	return snapshot
}

func applyMatchDayInitSnapshot(lineup *Lineup) {
	if lineup == nil {
		return
	}
	for _, event := range lineup.Match.Events {
		if event.EventType != matchDayEventInit || strings.TrimSpace(event.EventData) == "" {
			continue
		}
		var snapshot matchDayLineupSnapshot
		if err := json.Unmarshal([]byte(event.EventData), &snapshot); err != nil {
			continue
		}
		lineup.Name = snapshot.LineupName
		lineup.Details = snapshot.LineupDetails
		lineup.FormationID = snapshot.FormationID
		lineup.Formation = Formation{
			Model:  gorm.Model{ID: snapshot.FormationID},
			TeamID: lineup.TeamID,
			Name:   snapshot.FormationName,
			Status: formationStatusLive,
		}
		lineup.Formation.Positions = []FormationPosition{}
		for _, position := range snapshot.Positions {
			lineup.Formation.Positions = append(lineup.Formation.Positions, FormationPosition{
				FormationID:  snapshot.FormationID,
				IndexNumber:  position.IndexNumber,
				PositionName: position.PositionName,
				X:            position.X,
				Y:            position.Y,
			})
		}
		lineup.Players = []LineupPlayer{}
		for _, player := range snapshot.Players {
			lineup.Players = append(lineup.Players, LineupPlayer{
				LineupID:    lineup.ID,
				IndexNumber: player.IndexNumber,
				SlotOrder:   player.SlotOrder,
				SubMinute:   player.SubMinute,
				PlayerID:    player.PlayerID,
				Player:      Player{Model: gorm.Model{ID: player.PlayerID}, Name: player.PlayerName},
			})
		}
		return
	}
}

func ensureMatchDayInitEvent(tx *gorm.DB, lineup *Lineup, matchID uint, eventTime time.Time) error {
	fullLineup, err := getLineup(tx, lineup.ID)
	if err != nil {
		return err
	}
	fullLineup.Match = lineup.Match
	data, err := json.Marshal(matchDaySnapshotFromLineup(fullLineup))
	if err != nil {
		return err
	}
	var initEvent MatchEvent
	err = tx.Where("match_id = ? AND event_type = ?", matchID, matchDayEventInit).Order("created_at ASC").First(&initEvent).Error
	if err == nil {
		return tx.Model(&initEvent).Updates(map[string]interface{}{
			"event_name":   "Line-up snapshot",
			"event_time":   eventTime,
			"event_minute": 0,
			"player_id":    0,
			"event_data":   string(data),
		}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(&MatchEvent{
		MatchId:     uint64(matchID),
		EventName:   "Line-up snapshot",
		EventType:   matchDayEventInit,
		EventTime:   &eventTime,
		EventMinute: 0,
		EventData:   string(data),
	}).Error
}

func startLineupMatch(db *gorm.DB, lineup *Lineup, startTime time.Time) error {
	matchID := lineupMatchID(lineup)
	if matchID == 0 {
		return errors.New("line-up has no match")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := ensureMatchDayInitEvent(tx, lineup, matchID, startTime); err != nil {
			return err
		}
		if err := tx.Model(&Match{}).Where("id = ?", matchID).Updates(map[string]interface{}{
			"start_time": startTime,
			"lineup_id":  lineup.ID,
		}).Error; err != nil {
			return err
		}
		var startEvent MatchEvent
		err := tx.Where("match_id = ? AND event_type = ?", matchID, matchDayEventStarted).Order("created_at ASC").First(&startEvent).Error
		if err == nil {
			if err := tx.Model(&startEvent).Updates(map[string]interface{}{
				"event_name":   "Match started",
				"event_time":   startTime,
				"event_minute": 0,
				"player_id":    0,
			}).Error; err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := createMatchDayEventAt(tx, matchID, matchDayEventStarted, 0, 0, "Match started", startTime); err != nil {
				return err
			}
		} else {
			return err
		}
		return nil
	})
}

func lineupForMatch(db *gorm.DB, match Match) (*Lineup, error) {
	var lineup Lineup
	q := db.Preload("Players", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC, slot_order ASC")
	}).Preload("Players.Player").Preload("Formation.Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("index_number ASC")
	}).Preload("Match").Preload("Match.Events").Preload("Creator")
	if match.LineupID > 0 {
		if err := q.First(&lineup, match.LineupID).Error; err != nil {
			return nil, err
		}
		lineup.Match = match
		return &lineup, nil
	}
	if err := q.Joins("JOIN matches ON matches.lineup_id = lineups.id").Where("lineups.team_id = ? AND matches.id = ?", match.TeamID, match.ID).Order("lineups.updated_at DESC").First(&lineup).Error; err != nil {
		return nil, err
	}
	return &lineup, nil
}

func attachLineupToMatch(tx *gorm.DB, teamID uint, lineupID uint, matchID uint) error {
	if lineupID == 0 || matchID == 0 {
		return nil
	}
	return tx.Model(&Match{}).Where("id = ? AND team_id = ?", matchID, teamID).Update("lineup_id", lineupID).Error
}

func syncLineupMatch(tx *gorm.DB, lineup Lineup, matchID uint) error {
	if matchID == 0 {
		return nil
	}
	return attachLineupToMatch(tx, lineup.TeamID, lineup.ID, matchID)
}

func nextLineupSlot(players []LineupPlayer, indexNumber int) int {
	nextSlot := 0
	for _, player := range players {
		if player.IndexNumber == indexNumber && player.SlotOrder >= nextSlot {
			nextSlot = player.SlotOrder + 1
		}
	}
	return nextSlot
}

func validateLineupFormation(db *gorm.DB, actor LineupActor, lineupStatus string, formationID uint) error {
	if formationID == 0 {
		return nil
	}
	formation, err := getFormation(db, formationID)
	if err != nil {
		return err
	}
	if formation.TeamID != actor.TeamID {
		return errors.New("formation is not for this team")
	}
	if lineupStatus == lineupStatusDraft || lineupStatus == lineupStatusSubmitted {
		if formation.Status == formationStatusLive || formation.Status == formationStatusSubmitted || actorCanSeeDraft(actor, formation.TeamID, formation.CreatedByUserID) {
			return nil
		}
		return errors.New("draft formation is not visible")
	}
	if formation.Status != formationStatusLive {
		return errors.New(sharedLineupLiveFormationError)
	}
	return nil
}

func parseUintFormValue(r *http.Request, name string) (uint, error) {
	value := strings.TrimSpace(r.FormValue(name))
	if value == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(value, 10, 64)
	return uint(id), err
}

func lineupLoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
		case "POST":
			if err := r.ParseForm(); err != nil {
				lineupLoginPageForRequest(r, db, "Invalid form data").Render(GetContext(r, db), w)
				return
			}
			if strings.TrimSpace(r.FormValue("credential")) != "" {
				loginLineupUserFromGoogle(db, w, r)
				return
			}
			team, err := getTeamByKeyAndMemberPassword(db, r.FormValue("teamKey"), r.FormValue("teamMemberPass"))
			if err != nil {
				lineupLoginPageForRequest(r, db, "Team key or member password did not match").Render(GetContext(r, db), w)
				return
			}
			user, err := findOrCreateLineupUser(db, team.ID, r.FormValue("displayName"))
			if err != nil {
				lineupLoginPageForRequest(r, db, err.Error()).Render(GetContext(r, db), w)
				return
			}
			setLineupUserCookie(w, team.ID, user.ID)
			w.Header().Set("HX-Redirect", "/lineups")
			http.Redirect(w, r, "/lineups", http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func lineupListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		lineups, err := visibleLineups(db, actor)
		if err != nil {
			warning(fmt.Sprintf("Could not load line-ups: %v", err)).Render(GetContext(r, db), w)
			return
		}
		formations, err := visibleFormations(db, actor)
		if err != nil {
			warning(fmt.Sprintf("Could not load formations: %v", err)).Render(GetContext(r, db), w)
			return
		}
		var selected *Lineup
		hasUpcomingMatch := hasUpcomingMatchForTeam(db, actor.TeamID)
		matchDayLineupID := matchDayLineupIDForTeam(db, actor.TeamID)
		lineupListPage(lineups, formations, selected, actor, hasUpcomingMatch, matchDayLineupID).Render(GetContext(r, db), w)
	}
}

func hasUpcomingMatchForTeam(db *gorm.DB, teamID uint) bool {
	match, err := upcomingMatchForTeam(db, teamID)
	return err == nil && match != nil
}

func upcomingMatchLineupForTeam(db *gorm.DB, teamID uint) (*Match, *Lineup, error) {
	match, err := upcomingMatchForTeam(db, teamID)
	if err != nil {
		return nil, nil, err
	}
	lineup, err := lineupForMatch(db, *match)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return match, nil, nil
	}
	if err != nil {
		return match, nil, err
	}
	if lineup.TeamID != teamID || !matchDayLineupSelectable(*lineup) {
		return match, nil, nil
	}
	return match, lineup, nil
}

func matchDayLineupIDForTeam(db *gorm.DB, teamID uint) uint {
	_, lineup, err := upcomingMatchLineupForTeam(db, teamID)
	if err != nil || lineup == nil {
		return 0
	}
	return lineup.ID
}

func upcomingMatchForTeam(db *gorm.DB, teamID uint) (*Match, error) {
	if teamID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if os.Getenv("ACTIVE_MATCH_OVERRIDE") != "" {
		match, err := GetActiveMatch(db)
		if err != nil || match == nil || match.TeamID != teamID {
			return nil, gorm.ErrRecordNotFound
		}
		return match, nil
	}

	var match Match
	if err := db.Where("team_id = ? AND start_time > ?", teamID, time.Now()).Order("start_time ASC").First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func activeMatchForTeam(db *gorm.DB, teamID uint) (*Match, error) {
	if teamID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if os.Getenv("ACTIVE_MATCH_OVERRIDE") != "" {
		match, err := GetActiveMatch(db)
		if err != nil || match == nil || match.TeamID != teamID {
			return nil, gorm.ErrRecordNotFound
		}
		return match, nil
	}
	var match Match
	if err := db.Where("team_id = ? AND start_time > ?", teamID, time.Now().Add(-72*time.Hour)).Order("start_time ASC").First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func hasActiveMatchLineupForTeam(db *gorm.DB, teamID uint, excludeLineupID uint) bool {
	match, err := activeMatchForTeam(db, teamID)
	if err != nil || match == nil {
		return false
	}
	var count int64
	q := db.Model(&Lineup{}).Joins("JOIN matches ON matches.lineup_id = lineups.id").Where("lineups.team_id = ? AND matches.id = ?", teamID, match.ID)
	if excludeLineupID > 0 {
		q = q.Where("lineups.id != ?", excludeLineupID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func setActiveMatchLineup(db *gorm.DB, actor LineupActor, lineupID uint) error {
	match, err := activeMatchForTeam(db, actor.TeamID)
	if err != nil {
		return errors.New("no active match to set this as the line-up")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Lineup{}).Where("id = ? AND team_id = ?", lineupID, actor.TeamID).Updates(map[string]interface{}{
			"status": lineupStatusSelected,
		}).Error; err != nil {
			return err
		}
		return attachLineupToMatch(tx, actor.TeamID, lineupID, match.ID)
	})
}

func setActiveMatchFormationAsLineup(db *gorm.DB, actor LineupActor, formation *Formation) error {
	if formation == nil || formation.TeamID != actor.TeamID {
		return errors.New("formation not found for this team")
	}
	match, err := activeMatchForTeam(db, actor.TeamID)
	if err != nil {
		return errors.New("no active match to set this as the line-up")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var existing Lineup
		if match.LineupID > 0 {
			err := tx.Where("id = ? AND team_id = ?", match.LineupID, actor.TeamID).First(&existing).Error
			if err != nil {
				return err
			}
			if err := tx.Model(&Lineup{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"formation_id": formation.ID,
				"status":       lineupStatusSelected,
			}).Error; err != nil {
				return err
			}
			return attachLineupToMatch(tx, actor.TeamID, existing.ID, match.ID)
		}
		lineup := Lineup{
			TeamID:          actor.TeamID,
			CreatedByUserID: actor.UserID,
			FormationID:     formation.ID,
			Name:            fmt.Sprintf("%s line-up", formation.Name),
			Status:          lineupStatusSelected,
		}
		if err := tx.Create(&lineup).Error; err != nil {
			return err
		}
		return attachLineupToMatch(tx, actor.TeamID, lineup.ID, match.ID)
	})
}

func preferredMatchDayLineupID(lineups []Lineup) uint {
	for _, status := range []string{lineupStatusLive, lineupStatusSelected, lineupStatusProposed} {
		for _, lineup := range lineups {
			if lineup.Status == status && matchDayLineupSelectable(lineup) {
				return lineup.ID
			}
		}
	}
	for _, lineup := range lineups {
		if matchDayLineupSelectable(lineup) {
			return lineup.ID
		}
	}
	return 0
}

func matchDayLineupSelectable(lineup Lineup) bool {
	return lineup.FormationID > 0 && lineup.Formation.ID > 0
}

func visibleMatchDayLineupID(lineups []Lineup, lineupID uint) uint {
	if lineupID == 0 {
		return 0
	}
	for _, lineup := range lineups {
		if lineup.ID == lineupID && matchDayLineupSelectable(lineup) {
			return lineupID
		}
	}
	return 0
}

func matchDayLineupIDFromCookie(r *http.Request, lineups []Lineup) uint {
	cookie, err := r.Cookie(matchDayLineupCookieName)
	if err != nil {
		return 0
	}
	lineupID, err := strconv.ParseUint(strings.TrimSpace(cookie.Value), 10, 64)
	if err != nil {
		return 0
	}
	return visibleMatchDayLineupID(lineups, uint(lineupID))
}

func matchDayAdjacentLineups(db *gorm.DB, actor LineupActor, selected *Lineup) MatchDayAdjacentLineups {
	if selected == nil || lineupMatchID(selected) == 0 || selected.Match.StartTime == nil {
		return MatchDayAdjacentLineups{}
	}
	return MatchDayAdjacentLineups{
		Previous: matchDayAdjacentLineup(db, actor, *selected, "previous"),
		Next:     matchDayAdjacentLineup(db, actor, *selected, "next"),
	}
}

func matchDayIsHistorical(selected *Lineup, now time.Time) bool {
	if selected == nil || lineupMatchID(selected) == 0 || selected.Match.StartTime == nil {
		return false
	}
	return selected.Match.StartTime.Before(now)
}

func matchDayAdjacentLineup(db *gorm.DB, actor LineupActor, selected Lineup, direction string) *Lineup {
	q := db.Preload("Match").
		Joins("JOIN matches ON matches.lineup_id = lineups.id").
		Where("lineups.team_id = ? AND lineups.id != ? AND matches.lineup_id > 0 AND lineups.formation_id > 0", actor.TeamID, selected.ID).
		Where("matches.start_time IS NOT NULL")
	if !actor.IsAdmin {
		q = q.Where("lineups.status != ? OR lineups.created_by_user_id = ?", lineupStatusDraft, actor.UserID)
	}
	if direction == "previous" {
		q = q.Where("matches.start_time < ?", selected.Match.StartTime).Order("matches.start_time DESC")
	} else {
		q = q.Where("matches.start_time > ?", selected.Match.StartTime).Order("matches.start_time ASC")
	}
	var lineup Lineup
	if err := q.First(&lineup).Error; err != nil {
		return nil
	}
	return &lineup
}

func setMatchDayLineupCookie(w http.ResponseWriter, lineupID uint) {
	http.SetCookie(w, &http.Cookie{
		Name:     matchDayLineupCookieName,
		Value:    fmt.Sprintf("%d", lineupID),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 180,
	})
}

func matchDayHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		lineups, err := visibleLineups(db, actor)
		if err != nil {
			warning(fmt.Sprintf("Could not load line-ups: %v", err)).Render(GetContext(r, db), w)
			return
		}
		activeMatch, activeMatchErr := activeMatchForTeam(db, actor.TeamID)
		if activeMatchErr != nil && !errors.Is(activeMatchErr, gorm.ErrRecordNotFound) {
			activeMatch = nil
		}
		activeMatchNeedsLineup := actor.IsAdmin && activeMatch != nil && activeMatch.LineupID == 0
		var lineupID uint
		lineupIDParam := strings.TrimSpace(chi.URLParam(r, "lineupId"))
		if lineupIDParam != "" {
			parsedID, err := strconv.ParseUint(lineupIDParam, 10, 64)
			if err != nil {
				http.Error(w, "Invalid line-up", http.StatusBadRequest)
				return
			}
			lineupID = uint(parsedID)
		} else if !activeMatchNeedsLineup {
			lineupID = matchDayLineupIDFromCookie(r, lineups)
			if lineupID == 0 {
				lineupID = preferredMatchDayLineupID(lineups)
			}
		}
		var selected *Lineup
		pageErrorMsg := r.URL.Query().Get("errorMsg")
		if lineupID > 0 {
			selected, err = getLineup(db, lineupID)
			if err != nil {
				http.Error(w, "Line-up not found", http.StatusNotFound)
				return
			}
			if selected.TeamID != actor.TeamID || selected.Status == lineupStatusDraft && !actorCanEditLineup(actor, *selected) {
				http.Error(w, "Line-up not found", http.StatusNotFound)
				return
			}
			if !matchDayLineupSelectable(*selected) {
				selected = nil
				pageErrorMsg = "Choose a line-up with a formation attached before using match day."
			} else {
				applyMatchDayInitSnapshot(selected)
				setMatchDayLineupCookie(w, selected.ID)
			}
		}
		players, _ := FetchActivePlayers(db)
		matches, _ := matchDaySelectableMatchesForTeam(db, actor.TeamID)
		unavailablePlayerIDs := map[uint]bool{}
		if matchID := lineupMatchID(selected); matchID > 0 {
			unavailablePlayerIDs, _ = GetUnavailablePlayerIDsForMatch(db, selected.TeamID, matchID)
		}
		playerTimes := matchDayPlayerTimesWithEvents(selected)
		problems, err := matchDayErrors(db, selected, playerTimes, unavailablePlayerIDs)
		if err != nil {
			warning(fmt.Sprintf("Could not check match day: %v", err)).Render(GetContext(r, db), w)
			return
		}
		adjacentLineups := matchDayAdjacentLineups(db, actor, selected)
		isHistoricalMatch := matchDayIsHistorical(selected, time.Now())
		matchDayLineupID := matchDayLineupIDForTeam(db, actor.TeamID)
		focusPointsData := NotesPageData{}
		if selected != nil && selected.Match.ID > 0 {
			focusPointsData, err = matchDayFocusPointsData(db, actor.TeamID, selected.Match, fmt.Sprintf("/match-day/%d", selected.ID), actor.IsAdmin)
			if err != nil {
				warning(fmt.Sprintf("Could not load focus points: %v", err)).Render(GetContext(r, db), w)
				return
			}
		}
		matchDayPage(lineups, selected, playerTimes, players, unavailablePlayerIDs, problems, matches, activeMatch, adjacentLineups, isHistoricalMatch, actor.IsAdmin, r.URL.Query().Get("msg"), pageErrorMsg, activeMatch != nil, matchDayLineupID, matchDayRedoEventFromQuery(r.URL.Query()), focusPointsData, "", r.URL.Query().Get("resetPreview") == "true").Render(GetContext(r, db), w)
	}
}

func publicMatchURLHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		slug := strings.TrimSpace(chi.URLParam(r, "matchURLSlug"))
		if slug == "" {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}
		match, err := GetMatchByURLSlug(db, slug)
		if err != nil {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}
		team, err := GetTeam(db, match.TeamID)
		if err != nil || team == nil || !team.EnableLineupsModule {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}
		selected, err := lineupForMatch(db, *match)
		if err != nil {
			http.Error(w, "Match does not have a public line-up yet", http.StatusNotFound)
			return
		}
		if selected.TeamID != match.TeamID || !matchDayLineupSelectable(*selected) {
			http.Error(w, "Match does not have a public line-up yet", http.StatusNotFound)
			return
		}
		applyMatchDayInitSnapshot(selected)
		unavailablePlayerIDs, _ := GetUnavailablePlayerIDsForMatch(db, selected.TeamID, match.ID)
		playerTimes := matchDayPlayerTimesWithEvents(selected)
		problems, err := matchDayErrors(db, selected, playerTimes, unavailablePlayerIDs)
		if err != nil {
			warning(fmt.Sprintf("Could not check match day: %v", err)).Render(GetContext(r, db), w)
			return
		}
		focusPointsData, err := matchDayFocusPointsData(db, match.TeamID, selected.Match, publicMatchURLPath(*match), false)
		if err != nil {
			warning(fmt.Sprintf("Could not load focus points: %v", err)).Render(GetContext(r, db), w)
			return
		}
		ctx := context.WithValue(GetContext(r, db), "team_id", match.TeamID)
		ctx = context.WithValue(ctx, teamKey, *team)
		matchDayPage(nil, selected, playerTimes, nil, unavailablePlayerIDs, problems, nil, nil, MatchDayAdjacentLineups{}, matchDayIsHistorical(selected, time.Now()), false, r.URL.Query().Get("msg"), r.URL.Query().Get("errorMsg"), false, 0, nil, focusPointsData, publicMatchFeedbackURLPath(*match), false).Render(ctx, w)
	}
}

func matchDaySelectableMatchesForTeam(db *gorm.DB, teamID uint) ([]Match, error) {
	if teamID == 0 {
		return []Match{}, nil
	}
	var matches []Match
	if err := db.Where("team_id = ? AND (start_time IS NULL OR start_time > ?)", teamID, time.Now().Add(-72*time.Hour)).Order("start_time ASC").Limit(999).Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func matchDayActionHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		if lineup.TeamID != actor.TeamID || lineup.Status == lineupStatusDraft && !actorCanEditLineup(actor, *lineup) {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		setMatchDayLineupCookie(w, lineup.ID)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		action := strings.TrimSpace(r.FormValue("action"))
		redirectWithMessage := func(message string) {
			http.Redirect(w, r, fmt.Sprintf("/match-day/%d?msg=%s", lineup.ID, url.QueryEscape(message)), http.StatusSeeOther)
		}
		redirectWithRedoEvent := func(message string, event MatchEvent) {
			http.Redirect(w, r, fmt.Sprintf("/match-day/%d?%s", lineup.ID, matchDayRedoEventValues(message, event).Encode()), http.StatusSeeOther)
		}
		adminMatchDayActions := map[string]bool{
			"set-match":                 true,
			"start-match-now":           true,
			"start-match-at":            true,
			"custom-event":              true,
			"undo-event":                true,
			"redo-event":                true,
			"preview-reset-match":       true,
			"reset-match":               true,
			matchDayEventHalfTime:       true,
			matchDayEventFinished:       true,
			matchDayEventExtraTimeStart: true,
			matchDayEventResumed:        true,
		}
		if adminMatchDayActions[action] && !actor.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if lineup.Locked && (action == "set-match" || action == "sub") {
			http.Error(w, "Line-up is locked", http.StatusForbidden)
			return
		}
		if action == "set-match" {
			matchID, err := parseUintFormValue(r, "matchId")
			if err != nil || matchID == 0 {
				http.Redirect(w, r, fmt.Sprintf("/match-day/%d?errorMsg=%s", lineup.ID, url.QueryEscape("Choose a match to attach to this line-up.")), http.StatusSeeOther)
				return
			}
			match, err := GetMatch(db, matchID)
			if err != nil || match.TeamID != lineup.TeamID {
				http.Redirect(w, r, fmt.Sprintf("/match-day/%d?errorMsg=%s", lineup.ID, url.QueryEscape("Match not found for this team.")), http.StatusSeeOther)
				return
			}
			if err := db.Transaction(func(tx *gorm.DB) error {
				return syncLineupMatch(tx, *lineup, matchID)
			}); err != nil {
				http.Error(w, "Could not set match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match added to line-up.")
			return
		}
		matchID := lineupMatchID(lineup)
		if matchID == 0 {
			http.Redirect(w, r, fmt.Sprintf("/match-day/%d?errorMsg=%s", lineup.ID, url.QueryEscape("Choose a match before logging match-day actions.")), http.StatusSeeOther)
			return
		}
		minute, err := strconv.Atoi(strings.TrimSpace(r.FormValue("minute")))
		if err != nil {
			minute = matchDayMinute(lineup.Match)
		}
		minute = clampMatchMinute(minute)
		playerID, _ := parseUintFormValue(r, "playerId")
		switch action {
		case "start-match-now":
			if err := startLineupMatch(db, lineup, time.Now()); err != nil {
				http.Error(w, "Could not start match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match started now.")
		case "start-match-at":
			startTimeValue := strings.TrimSpace(r.FormValue("startTime"))
			if startTimeValue == "" {
				http.Error(w, "Start time is required", http.StatusBadRequest)
				return
			}
			startTime, err := time.Parse("2006-01-02T15:04", startTimeValue)
			if err != nil {
				http.Error(w, "Invalid start time", http.StatusBadRequest)
				return
			}
			if err := startLineupMatch(db, lineup, startTime); err != nil {
				http.Error(w, "Could not start match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match start time updated.")
		case "goal-against":
			if err := createMatchDayEvent(db, matchID, "conceded-goal", 0, minute, "Goal for them"); err != nil {
				http.Error(w, "Could not save goal", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Goal added for them.")
		case matchDayEventHalfTime:
			if err := createMatchDayEvent(db, matchID, matchDayEventHalfTime, 0, minute, "Half-time"); err != nil {
				http.Error(w, "Could not save half-time", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Half-time recorded.")
		case matchDayEventExtraTimeStart:
			if err := createMatchDayEvent(db, matchID, matchDayEventExtraTimeStart, 0, matchDayRegularEndMinute, "Extra time started"); err != nil {
				http.Error(w, "Could not start extra time", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Extra time started.")
		case matchDayEventResumed:
			minutesAgo, err := strconv.Atoi(strings.TrimSpace(r.FormValue("minutesAgo")))
			if err != nil || minutesAgo < 0 {
				minutesAgo = 0
			}
			if minutesAgo > 5 {
				minutesAgo = 5
			}
			eventTime := time.Now().Add(-time.Duration(minutesAgo) * time.Minute)
			resumeMinute := matchDayResumeMinute(lineup.Match)
			if minute > 0 {
				resumeMinute = minute
			}
			if err := createMatchDayEventAt(db, matchID, matchDayEventResumed, 0, resumeMinute, "Match resumed", eventTime); err != nil {
				http.Error(w, "Could not resume match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match resumed.")
		case matchDayEventFinished:
			if err := createMatchDayEvent(db, matchID, matchDayEventFinished, 0, minute, "Match finished"); err != nil {
				http.Error(w, "Could not finish match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match finished.")
		case "undo-event":
			eventID, err := parseUintFormValue(r, "eventId")
			if err != nil || eventID == 0 {
				http.Redirect(w, r, fmt.Sprintf("/match-day/%d?errorMsg=%s", lineup.ID, url.QueryEscape("Choose an event to undo.")), http.StatusSeeOther)
				return
			}
			var event MatchEvent
			if err := db.Where("id = ? AND match_id = ?", eventID, matchID).First(&event).Error; err != nil {
				http.Redirect(w, r, fmt.Sprintf("/match-day/%d?errorMsg=%s", lineup.ID, url.QueryEscape("Event not found for this match.")), http.StatusSeeOther)
				return
			}
			if err := db.Delete(&event).Error; err != nil {
				http.Error(w, "Could not undo event", http.StatusInternalServerError)
				return
			}
			redirectWithRedoEvent("Match event undone.", event)
		case "preview-reset-match":
			http.Redirect(w, r, fmt.Sprintf("/match-day/%d?resetPreview=true", lineup.ID), http.StatusSeeOther)
		case "reset-match":
			if r.FormValue("confirmReset") != "delete-events" {
				http.Redirect(w, r, fmt.Sprintf("/match-day/%d?resetPreview=true&errorMsg=%s", lineup.ID, url.QueryEscape("Confirm the reset before deleting match events.")), http.StatusSeeOther)
				return
			}
			result := db.Where("match_id = ?", matchID).Delete(&MatchEvent{})
			if result.Error != nil {
				http.Error(w, "Could not reset match", http.StatusInternalServerError)
				return
			}
			redirectWithMessage(fmt.Sprintf("Match reset. Deleted %d events.", result.RowsAffected))
		case "redo-event":
			eventType := strings.TrimSpace(r.FormValue("eventType"))
			if eventType == "" {
				http.Error(w, "Event type is required", http.StatusBadRequest)
				return
			}
			eventTime := time.Now()
			if eventTimeValue := strings.TrimSpace(r.FormValue("eventTime")); eventTimeValue != "" {
				parsedTime, err := time.Parse(time.RFC3339Nano, eventTimeValue)
				if err != nil {
					http.Error(w, "Invalid event time", http.StatusBadRequest)
					return
				}
				eventTime = parsedTime
			}
			eventName := strings.TrimSpace(r.FormValue("eventName"))
			if eventName == "" {
				eventName = eventType
			}
			if err := createMatchDayEventAt(db, matchID, eventType, playerID, minute, eventName, eventTime); err != nil {
				http.Error(w, "Could not re-do event", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Match event re-done.")
		case "custom-event":
			eventType := strings.TrimSpace(r.FormValue("eventType"))
			if eventType == "" {
				http.Error(w, "Event type is required", http.StatusBadRequest)
				return
			}
			eventName := strings.TrimSpace(r.FormValue("eventName"))
			if eventName == "" {
				eventName = eventType
			}
			if err := createMatchDayEvent(db, matchID, eventType, playerID, minute, eventName); err != nil {
				http.Error(w, "Could not save custom event", http.StatusInternalServerError)
				return
			}
			redirectWithMessage("Custom match event added.")
		case "injury", "assist", "yellow-card", "red-card", "gave-away-penalty":
			if playerID == 0 {
				http.Error(w, "Player is required", http.StatusBadRequest)
				return
			}
			player, err := GetPlayerByID(db, playerID)
			if err != nil {
				http.Error(w, "Player not found", http.StatusBadRequest)
				return
			}
			eventLabel := action
			if label, ok := matchDayExtraPlayerEventLabel(action); ok {
				eventLabel = label
			}
			eventName := fmt.Sprintf("%s - %s", player.Name, eventLabel)
			if err := createMatchDayEvent(db, matchID, action, player.ID, minute, eventName); err != nil {
				http.Error(w, "Could not save event", http.StatusInternalServerError)
				return
			}
			redirectWithMessage(fmt.Sprintf("%s marked for %s.", player.Name, eventLabel))
		case "goal":
			if playerID == 0 {
				http.Error(w, "Player is required", http.StatusBadRequest)
				return
			}
			player, err := GetPlayerByID(db, playerID)
			if err != nil {
				http.Error(w, "Player not found", http.StatusBadRequest)
				return
			}
			assisterID, _ := parseUintFormValue(r, "assisterPlayerId")
			var assister *Player
			if assisterID > 0 && assisterID != player.ID {
				assister, err = GetPlayerByID(db, assisterID)
				if err != nil {
					http.Error(w, "Assister not found", http.StatusBadRequest)
					return
				}
			}
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := createMatchDayEvent(tx, matchID, "goal", player.ID, minute, fmt.Sprintf("%s - goal", player.Name)); err != nil {
					return err
				}
				if assister != nil {
					return createMatchDayEvent(tx, matchID, "assist", assister.ID, minute, fmt.Sprintf("%s - assist", assister.Name))
				}
				return nil
			}); err != nil {
				http.Error(w, "Could not save goal", http.StatusInternalServerError)
				return
			}
			if assister != nil {
				redirectWithMessage(fmt.Sprintf("%s goal, assisted by %s.", player.Name, assister.Name))
			} else {
				redirectWithMessage(fmt.Sprintf("%s goal added.", player.Name))
			}
		case "swap":
			if playerID == 0 {
				http.Error(w, "Player is required", http.StatusBadRequest)
				return
			}
			swapPlayerID, err := parseUintFormValue(r, "swapPlayerId")
			if err != nil || swapPlayerID == 0 || swapPlayerID == playerID {
				http.Error(w, "Swap player is required", http.StatusBadRequest)
				return
			}
			player, err := GetPlayerByID(db, playerID)
			if err != nil {
				http.Error(w, "Player not found", http.StatusBadRequest)
				return
			}
			swapPlayer, err := GetPlayerByID(db, swapPlayerID)
			if err != nil {
				http.Error(w, "Swap player not found", http.StatusBadRequest)
				return
			}
			playerTimes := matchDayPlayerTimesWithEvents(lineup)
			if !matchDayPlayerCurrentlyOn(playerTimes, player.ID) || !matchDayPlayerCurrentlyOn(playerTimes, swapPlayer.ID) {
				http.Error(w, "Both players must be on the pitch to swap", http.StatusBadRequest)
				return
			}
			now := time.Now()
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&MatchEvent{
					MatchId:     uint64(matchID),
					EventName:   fmt.Sprintf("%s swapped with %s", player.Name, swapPlayer.Name),
					EventType:   "swap",
					EventTime:   &now,
					EventMinute: minute,
					PlayerId:    player.ID,
				}).Error; err != nil {
					return err
				}
				return tx.Create(&MatchEvent{
					MatchId:     uint64(matchID),
					EventName:   fmt.Sprintf("%s swapped with %s", swapPlayer.Name, player.Name),
					EventType:   "swap",
					EventTime:   &now,
					EventMinute: minute,
					PlayerId:    swapPlayer.ID,
				}).Error
			}); err != nil {
				http.Error(w, "Could not save swap", http.StatusInternalServerError)
				return
			}
			redirectWithMessage(fmt.Sprintf("%s and %s swapped.", player.Name, swapPlayer.Name))
		case "sub":
			if playerID == 0 {
				http.Error(w, "Player is required", http.StatusBadRequest)
				return
			}
			replacementID, err := parseUintFormValue(r, "replacementPlayerId")
			if err != nil || replacementID == 0 {
				http.Error(w, "Replacement player is required", http.StatusBadRequest)
				return
			}
			currentPlayer, err := GetPlayerByID(db, playerID)
			if err != nil {
				http.Error(w, "Player not found", http.StatusBadRequest)
				return
			}
			replacementPlayer, err := GetPlayerByID(db, replacementID)
			if err != nil {
				http.Error(w, "Replacement player not found", http.StatusBadRequest)
				return
			}
			unavailablePlayerIDs, err := GetUnavailablePlayerIDsForMatch(db, lineup.TeamID, matchID)
			if err != nil {
				http.Error(w, "Could not check player availability", http.StatusInternalServerError)
				return
			}
			if playerUnavailableForMatch(unavailablePlayerIDs, replacementPlayer.ID) {
				http.Error(w, "Replacement player is unavailable for this match", http.StatusBadRequest)
				return
			}
			if matchDayPlayerActuallyCurrentlyOn(lineup, replacementPlayer.ID) {
				http.Error(w, "Replacement player is already on the pitch", http.StatusBadRequest)
				return
			}
			markInjured := r.FormValue("markInjured") == "on"
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := createMatchDayEvent(tx, matchID, "subbed-off", currentPlayer.ID, minute, fmt.Sprintf("%s subbed off", currentPlayer.Name)); err != nil {
					return err
				}
				if markInjured {
					if err := createMatchDayEvent(tx, matchID, "injury", currentPlayer.ID, minute, fmt.Sprintf("%s - injury", currentPlayer.Name)); err != nil {
						return err
					}
				}
				return createMatchDayEvent(tx, matchID, "subbed-on", replacementPlayer.ID, minute, fmt.Sprintf("%s subbed on", replacementPlayer.Name))
			}); err != nil {
				http.Error(w, "Could not save substitution", http.StatusInternalServerError)
				return
			}
			if markInjured {
				redirectWithMessage(fmt.Sprintf("%s on for %s, injury marked.", replacementPlayer.Name, currentPlayer.Name))
			} else {
				redirectWithMessage(fmt.Sprintf("%s on for %s.", replacementPlayer.Name, currentPlayer.Name))
			}
		default:
			http.Error(w, "Unknown action", http.StatusBadRequest)
		}
	}
}

func lineupNewHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		renderNewPage := func(msg string) {
			formations, _ := eligibleFormationsForLineup(db, actor, lineupStatusDraft)
			matches, _ := GetMatches(db, activeSeasonID(db), 0, 999)
			var selected *Formation
			hasUpcomingMatch := hasUpcomingMatchForTeam(db, actor.TeamID)
			matchDayLineupID := matchDayLineupIDForTeam(db, actor.TeamID)
			lineupNewPage(formations, matches, selected, actor.IsAdmin, msg, hasUpcomingMatch, matchDayLineupID).Render(GetContext(r, db), w)
		}
		switch r.Method {
		case "GET":
			renderNewPage("")
		case "POST":
			if err := r.ParseForm(); err != nil {
				renderNewPage("Invalid form data")
				return
			}
			formationID, err := parseUintFormValue(r, "formationId")
			if err != nil {
				renderNewPage("Invalid formation")
				return
			}
			matchID, err := parseUintFormValue(r, "matchId")
			if err != nil {
				renderNewPage("Invalid match")
				return
			}
			if matchID > 0 {
				if _, err := GetMatch(db, matchID); err != nil {
					renderNewPage("Match not found")
					return
				}
			}
			if err := validateLineupFormation(db, actor, lineupStatusDraft, formationID); err != nil {
				renderNewPage(err.Error())
				return
			}
			lineup := Lineup{
				TeamID:          actor.TeamID,
				CreatedByUserID: actor.UserID,
				FormationID:     formationID,
				Name:            strings.TrimSpace(r.FormValue("name")),
				Details:         strings.TrimSpace(r.FormValue("details")),
				Status:          lineupStatusDraft,
			}
			if lineup.Name == "" {
				lineup.Name = "Untitled line-up"
			}
			if err := db.Create(&lineup).Error; err != nil {
				renderNewPage(fmt.Sprintf("Could not create line-up: %v", err))
				return
			}
			if matchID > 0 {
				if err := syncLineupMatch(db, lineup, matchID); err != nil {
					renderNewPage(fmt.Sprintf("Could not attach line-up to match: %v", err))
					return
				}
			}
			http.Redirect(w, r, fmt.Sprintf("/lineups/%d", lineup.ID), http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func lineupDetailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		baseCanEdit := actorCanEditLineup(actor, *lineup)
		canEdit := baseCanEdit && !lineup.Locked
		if lineup.Status == lineupStatusDraft && !baseCanEdit {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		switch r.Method {
		case "GET":
			players, _ := GetPlayers(db, 0, 999)
			formations, _ := eligibleFormationsForLineup(db, actor, lineup.Status)
			matches, _ := GetMatches(db, activeSeasonID(db), 0, 999)
			unavailablePlayerIDs, _ := GetUnavailablePlayerIDsForMatch(db, lineup.TeamID, lineupMatchID(lineup))
			activeMatch, _ := activeMatchForTeam(db, actor.TeamID)
			matchDayLineupID := matchDayLineupIDForTeam(db, actor.TeamID)
			lineupEditPage(lineup, players, formations, matches, unavailablePlayerIDs, canEdit, baseCanEdit, actorCanCopyLineup(actor, *lineup), actor.IsAdmin, "", activeMatch, hasActiveMatchLineupForTeam(db, actor.TeamID, lineup.ID), matchDayLineupID).Render(GetContext(r, db), w)
		case "DELETE":
			if !canEdit {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if err := deleteLineup(db, lineup.ID); err != nil {
				http.Error(w, "Could not delete line-up", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Redirect", "/lineups")
			http.Redirect(w, r, "/lineups", http.StatusSeeOther)
		case "POST":
			if err := r.ParseForm(); err != nil {
				warning("Invalid form data").Render(GetContext(r, db), w)
				return
			}
			lineupAction := strings.TrimSpace(r.FormValue("lineupAction"))
			switch lineupAction {
			case "lock":
				if !baseCanEdit {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				if err := db.Model(&Lineup{}).Where("id = ?", lineup.ID).Updates(map[string]interface{}{
					"locked":            true,
					"locked_by_user_id": actor.UserID,
					"locked_by_name":    lineupLockActorName(r, db, actor),
				}).Error; err != nil {
					warning(fmt.Sprintf("Could not lock line-up: %v", err)).Render(GetContext(r, db), w)
					return
				}
				http.Redirect(w, r, fmt.Sprintf("/lineups/%d", lineup.ID), http.StatusSeeOther)
				return
			case "unlock":
				if !baseCanEdit {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				if err := db.Model(&Lineup{}).Where("id = ?", lineup.ID).Updates(map[string]interface{}{
					"locked":            false,
					"locked_by_user_id": 0,
					"locked_by_name":    "",
				}).Error; err != nil {
					warning(fmt.Sprintf("Could not unlock line-up: %v", err)).Render(GetContext(r, db), w)
					return
				}
				http.Redirect(w, r, fmt.Sprintf("/lineups/%d", lineup.ID), http.StatusSeeOther)
				return
			}
			if !canEdit {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			status := r.FormValue("status")
			if status == "" {
				status = lineup.Status
			}
			switch lineupAction {
			case "publish":
				if lineup.Status == lineupStatusDraft {
					if actor.IsAdmin {
						status = lineupStatusSelected
					} else {
						status = lineupStatusSubmitted
					}
				}
			case "unselect":
				if actor.IsAdmin && lineup.Status == lineupStatusSelected {
					status = lineupStatusProposed
				}
			case "set-active-match-lineup":
				if !actor.IsAdmin || lineup.Status != lineupStatusLive {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
			if !isValidLineupStatus(status) || !actor.IsAdmin && status != lineupStatusDraft && status != lineupStatusSubmitted {
				status = lineupStatusDraft
			}
			formationID, err := parseUintFormValue(r, "formationId")
			if err != nil {
				warning("Invalid formation").Render(GetContext(r, db), w)
				return
			}
			if err := validateLineupFormation(db, actor, status, formationID); err != nil {
				players, _ := GetPlayers(db, 0, 999)
				formations, _ := eligibleFormationsForLineup(db, actor, lineup.Status)
				matches, _ := GetMatches(db, activeSeasonID(db), 0, 999)
				unavailablePlayerIDs, _ := GetUnavailablePlayerIDsForMatch(db, lineup.TeamID, lineupMatchID(lineup))
				activeMatch, _ := activeMatchForTeam(db, actor.TeamID)
				matchDayLineupID := matchDayLineupIDForTeam(db, actor.TeamID)
				lineupEditPage(lineup, players, formations, matches, unavailablePlayerIDs, canEdit, baseCanEdit, actorCanCopyLineup(actor, *lineup), actor.IsAdmin, err.Error(), activeMatch, hasActiveMatchLineupForTeam(db, actor.TeamID, lineup.ID), matchDayLineupID).Render(GetContext(r, db), w)
				return
			}
			matchID, err := parseUintFormValue(r, "matchId")
			if err != nil {
				warning("Invalid match").Render(GetContext(r, db), w)
				return
			}
			updates := map[string]interface{}{
				"name":         strings.TrimSpace(r.FormValue("name")),
				"details":      strings.TrimSpace(r.FormValue("details")),
				"status":       status,
				"formation_id": formationID,
			}
			if updates["name"] == "" {
				updates["name"] = "Untitled line-up"
			}
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&Lineup{}).Where("id = ?", lineup.ID).Updates(updates).Error; err != nil {
					return err
				}
				return syncLineupMatch(tx, *lineup, matchID)
			}); err != nil {
				warning(fmt.Sprintf("Could not save line-up: %v", err)).Render(GetContext(r, db), w)
				return
			}
			if err := saveLineupPlayers(db, lineup.ID, r); err != nil {
				warning(fmt.Sprintf("Could not save players: %v", err)).Render(GetContext(r, db), w)
				return
			}
			if lineupAction == "set-active-match-lineup" {
				if err := setActiveMatchLineup(db, actor, lineup.ID); err != nil {
					warning(err.Error()).Render(GetContext(r, db), w)
					return
				}
			}
			http.Redirect(w, r, fmt.Sprintf("/lineups/%d", lineup.ID), http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func deleteLineup(db *gorm.DB, lineupID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Match{}).Where("lineup_id = ?", lineupID).Update("lineup_id", 0).Error; err != nil {
			return err
		}
		if err := tx.Where("lineup_id = ?", lineupID).Delete(&LineupPlayer{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Lineup{}, lineupID).Error
	})
}

func lineupErrorsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		canEdit := actorCanEditLineup(actor, *lineup)
		if lineup.Status == lineupStatusDraft && !canEdit {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		problems, err := lineupProblems(db, lineup)
		if err != nil {
			warning(fmt.Sprintf("Could not check line-up: %v", err)).Render(GetContext(r, db), w)
			return
		}
		lineupProblemsPanel(problems).Render(GetContext(r, db), w)
	}
}

func lineupCopyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		if !actorCanCopyLineup(actor, *lineup) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		if strings.TrimSpace(name) == "" {
			name = lineup.Name
		}
		formationID, err := parseUintFormValue(r, "formationId")
		if err != nil {
			http.Error(w, "Invalid formation", http.StatusBadRequest)
			return
		}
		if formationID == 0 {
			formationID = lineup.FormationID
		}
		if err := validateLineupFormation(db, actor, lineupStatusDraft, formationID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		matchID, err := parseUintFormValue(r, "matchId")
		if err != nil {
			http.Error(w, "Invalid match", http.StatusBadRequest)
			return
		}
		if matchID == 0 {
			matchID = lineupMatchID(lineup)
		}

		copyLineup := Lineup{
			TeamID:          actor.TeamID,
			CreatedByUserID: actor.UserID,
			FormationID:     formationID,
			Name:            copyName(name, lineupActorName(db, actor)),
			Details:         strings.TrimSpace(r.FormValue("details")),
			Status:          lineupStatusDraft,
		}
		if copyLineup.Details == "" {
			copyLineup.Details = lineup.Details
		}
		if err := db.Create(&copyLineup).Error; err != nil {
			http.Error(w, "Could not copy line-up", http.StatusInternalServerError)
			return
		}
		if matchID > 0 {
			if err := syncLineupMatch(db, copyLineup, matchID); err != nil {
				http.Error(w, "Could not attach copied line-up to match", http.StatusInternalServerError)
				return
			}
		}
		if len(r.Form) > 0 {
			if err := saveLineupPlayers(db, copyLineup.ID, r); err != nil {
				http.Error(w, "Could not copy players", http.StatusInternalServerError)
				return
			}
		} else {
			for _, lp := range lineup.Players {
				if lp.PlayerID == 0 {
					continue
				}
				if err := db.Create(&LineupPlayer{LineupID: copyLineup.ID, IndexNumber: lp.IndexNumber, SlotOrder: lp.SlotOrder, SubMinute: lp.SubMinute, PlayerID: lp.PlayerID}).Error; err != nil {
					http.Error(w, "Could not copy players", http.StatusInternalServerError)
					return
				}
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/lineups/%d", copyLineup.ID), http.StatusSeeOther)
	}
}

func saveLineupPlayers(db *gorm.DB, lineupID uint, r *http.Request) error {
	if err := db.Where("lineup_id = ?", lineupID).Delete(&LineupPlayer{}).Error; err != nil {
		return err
	}
	movePlayerIDs := parseLineupMovePlayerIDs(r)
	type lineupPlayerFormEntry struct {
		index     int
		slotOrder int
		playerID  uint
		subMinute int
	}
	entries := []lineupPlayerFormEntry{}
	for key, values := range r.Form {
		if !strings.HasPrefix(key, "player_") || len(values) == 0 || values[0] == "" {
			continue
		}
		fieldParts := strings.Split(strings.TrimPrefix(key, "player_"), "_")
		if len(fieldParts) == 0 {
			continue
		}
		index, err := strconv.Atoi(fieldParts[0])
		if err != nil {
			return err
		}
		slotOrder := 0
		if len(fieldParts) > 1 {
			slotOrder, err = strconv.Atoi(fieldParts[1])
			if err != nil {
				return err
			}
		}
		playerID, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return err
		}
		if playerID == 0 {
			continue
		}
		subMinute := 0
		if minuteStr := r.FormValue(fmt.Sprintf("minute_%d_%d", index, slotOrder)); minuteStr != "" {
			subMinute, err = strconv.Atoi(minuteStr)
			if err != nil {
				return err
			}
		}
		entries = append(entries, lineupPlayerFormEntry{
			index:     index,
			slotOrder: slotOrder,
			playerID:  uint(playerID),
			subMinute: subMinute,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].index == entries[j].index {
			return entries[i].slotOrder < entries[j].slotOrder
		}
		return entries[i].index < entries[j].index
	})
	playerByID := map[uint]LineupPlayer{}
	for _, entry := range entries {
		lineupPlayer := LineupPlayer{LineupID: lineupID, IndexNumber: entry.index, SlotOrder: entry.slotOrder, SubMinute: entry.subMinute, PlayerID: entry.playerID}
		if !movePlayerIDs[entry.playerID] {
			if err := db.Create(&lineupPlayer).Error; err != nil {
				return err
			}
			continue
		}
		playerByID[entry.playerID] = lineupPlayer
	}
	for _, lineupPlayer := range playerByID {
		if err := db.Create(&lineupPlayer).Error; err != nil {
			return err
		}
	}
	return nil
}

func parseLineupMovePlayerIDs(r *http.Request) map[uint]bool {
	playerIDs := map[uint]bool{}
	for _, value := range r.Form["lineupMovePlayerIds"] {
		playerID, err := strconv.ParseUint(value, 10, 64)
		if err == nil && playerID > 0 {
			playerIDs[uint(playerID)] = true
		}
	}
	return playerIDs
}

func lineupSlotHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		if !actorCanEditLineup(actor, *lineup) || lineup.Locked {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		indexNumber, err := strconv.Atoi(r.FormValue("indexNumber"))
		if err != nil {
			http.Error(w, "Invalid position", http.StatusBadRequest)
			return
		}
		if err := db.Where("lineup_id = ? AND index_number = ?", lineup.ID, indexNumber).Delete(&LineupPlayer{}).Error; err != nil {
			http.Error(w, "Could not clear position", http.StatusInternalServerError)
			return
		}
		movePlayerIDs := parseLineupMovePlayerIDs(r)
		seenPlayers := map[uint]bool{}
		for slotOrder := 0; slotOrder < 4; slotOrder++ {
			playerID, err := parseUintFormValue(r, fmt.Sprintf("playerId_%d", slotOrder))
			if err != nil {
				http.Error(w, "Invalid player", http.StatusBadRequest)
				return
			}
			if playerID == 0 || seenPlayers[playerID] {
				continue
			}
			if _, err := GetPlayerByID(db, playerID); err != nil {
				http.Error(w, "Player not found", http.StatusBadRequest)
				return
			}
			if movePlayerIDs[playerID] {
				if err := db.Where("lineup_id = ? AND player_id = ? AND index_number != ?", lineup.ID, playerID, indexNumber).Delete(&LineupPlayer{}).Error; err != nil {
					http.Error(w, "Could not move player", http.StatusInternalServerError)
					return
				}
			}
			subMinute := 0
			if minuteStr := r.FormValue(fmt.Sprintf("subMinute_%d", slotOrder)); minuteStr != "" {
				subMinute, err = strconv.Atoi(minuteStr)
				if err != nil {
					http.Error(w, "Invalid sub minute", http.StatusBadRequest)
					return
				}
			}
			if err := db.Create(&LineupPlayer{LineupID: lineup.ID, IndexNumber: indexNumber, SlotOrder: slotOrder, SubMinute: subMinute, PlayerID: playerID}).Error; err != nil {
				http.Error(w, "Could not save player", http.StatusInternalServerError)
				return
			}
			seenPlayers[playerID] = true
		}
		http.Redirect(w, r, fmt.Sprintf("/lineups/%d", lineup.ID), http.StatusSeeOther)
	}
}

func lineupShareHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "lineupId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid line-up", http.StatusBadRequest)
			return
		}
		lineup, err := getLineup(db, uint(id))
		if err != nil {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		if lineup.Status == lineupStatusDraft && !actorCanEditLineup(actor, *lineup) {
			http.Error(w, "Line-up not found", http.StatusNotFound)
			return
		}
		lineupSharePage(lineup, matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
	}
}

func formationListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		formations, err := visibleFormations(db, actor)
		if err != nil {
			warning(fmt.Sprintf("Could not load formations: %v", err)).Render(GetContext(r, db), w)
			return
		}
		var selected *Formation
		formationListPage(formations, selected, actor, hasUpcomingMatchForTeam(db, actor.TeamID), matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
	}
}

func formationNewHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		switch r.Method {
		case "GET":
			formationNewPage(actor.IsAdmin, "", hasUpcomingMatchForTeam(db, actor.TeamID), matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
		case "POST":
			if err := r.ParseForm(); err != nil {
				formationNewPage(actor.IsAdmin, "Invalid form data", hasUpcomingMatchForTeam(db, actor.TeamID), matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
				return
			}
			name := strings.TrimSpace(r.FormValue("name"))
			formation := Formation{
				TeamID:          actor.TeamID,
				CreatedByUserID: actor.UserID,
				Name:            name,
				Status:          formationStatusDraft,
			}
			if formation.Name == "" {
				formation.Name = "Untitled formation"
			}
			if err := db.Create(&formation).Error; err != nil {
				formationNewPage(actor.IsAdmin, fmt.Sprintf("Could not create formation: %v", err), hasUpcomingMatchForTeam(db, actor.TeamID), matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/formations/%d/edit?mode=edit&drawer=expanded", formation.ID), http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func formationEditHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "formationId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid formation", http.StatusBadRequest)
			return
		}
		formation, err := getFormation(db, uint(id))
		if err != nil {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		canEdit := actorCanEditFormation(actor, *formation)
		if formation.Status == formationStatusDraft && !actorCanSeeDraft(actor, formation.TeamID, formation.CreatedByUserID) {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		renderEdit := func(msg string, editMode bool) {
			notesData, err := contextualNotesData(db, actor.TeamID, noteTargetFormation, formation.ID, formation.Name, fmt.Sprintf("/formations/%d/edit", formation.ID))
			if err != nil {
				msg = strings.TrimSpace(msg + " " + fmt.Sprintf("Could not load model notes: %v", err))
			}
			notesData.ShowAdminNotes = actor.IsAdmin
			activeMatch, _ := activeMatchForTeam(db, actor.TeamID)
			formationEditPage(formation, canEdit, actor.IsAdmin, editMode, msg, notesData, r.URL.Query().Get("drawer") == "expanded", activeMatch, hasActiveMatchLineupForTeam(db, actor.TeamID, 0), matchDayLineupIDForTeam(db, actor.TeamID)).Render(GetContext(r, db), w)
		}
		switch r.Method {
		case "GET":
			renderEdit("", r.URL.Query().Get("mode") == "edit")
		case "DELETE":
			if !canEdit {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if err := deleteFormation(db, formation.ID); err != nil {
				http.Error(w, "Could not delete formation", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Redirect", "/formations")
			http.Redirect(w, r, "/formations", http.StatusSeeOther)
		case "POST":
			if !canEdit {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if err := r.ParseForm(); err != nil {
				renderEdit("Invalid form data", true)
				return
			}
			status := r.FormValue("status")
			formationAction := strings.TrimSpace(r.FormValue("formationAction"))
			switch formationAction {
			case "publish":
				if formation.Status == formationStatusDraft {
					if actor.IsAdmin {
						status = formationStatusLive
					} else {
						status = formationStatusSubmitted
					}
				}
			case "set-active-match-lineup":
				if !actor.IsAdmin || formation.Status != formationStatusLive {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}
			if !isValidFormationStatus(status) || !actor.IsAdmin && status != formationStatusDraft && status != formationStatusSubmitted {
				status = formationStatusDraft
			}
			if err := db.Model(&Formation{}).Where("id = ?", formation.ID).Updates(map[string]interface{}{
				"name":   strings.TrimSpace(r.FormValue("name")),
				"status": status,
			}).Error; err != nil {
				renderEdit(fmt.Sprintf("Could not save formation: %v", err), true)
				return
			}
			if err := saveFormationPositions(db, formation.ID, r); err != nil {
				renderEdit(fmt.Sprintf("Could not save positions: %v", err), true)
				return
			}
			if formationAction == "set-active-match-lineup" {
				formation.Name = strings.TrimSpace(r.FormValue("name"))
				if formation.Name == "" {
					formation.Name = "Untitled formation"
				}
				formation.Status = status
				if err := setActiveMatchFormationAsLineup(db, actor, formation); err != nil {
					renderEdit(err.Error(), true)
					return
				}
			}
			http.Redirect(w, r, fmt.Sprintf("/formations/%d/edit", formation.ID), http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func deleteFormation(db *gorm.DB, formationID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("formation_id = ?", formationID).Delete(&FormationPosition{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Formation{}, formationID).Error
	})
}

func formationLiveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.ParseUint(chi.URLParam(r, "formationId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid formation", http.StatusBadRequest)
			return
		}
		formation, err := getFormation(db, uint(id))
		if err != nil {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		actor := getLineupActor(r, db)
		if formation.TeamID != actor.TeamID || !actor.IsAdmin {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		if err := db.Model(&Formation{}).Where("id = ?", formation.ID).Update("status", formationStatusLive).Error; err != nil {
			http.Error(w, "Could not make formation live", http.StatusInternalServerError)
			return
		}
		returnTo := strings.TrimSpace(r.FormValue("returnTo"))
		if returnTo == "" || !strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "//") {
			returnTo = fmt.Sprintf("/formations/%d/edit", formation.ID)
		}
		http.Redirect(w, r, returnTo, http.StatusSeeOther)
	}
}

func saveFormationPositions(db *gorm.DB, formationID uint, r *http.Request) error {
	formation, err := getFormation(db, formationID)
	if err != nil {
		return err
	}
	for _, idStr := range r.Form["positionIds"] {
		positionID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return err
		}
		prefix := fmt.Sprintf("position_%d_", positionID)
		indexNumber, err := strconv.Atoi(r.FormValue(prefix + "index"))
		if err != nil {
			return err
		}
		x, err := strconv.ParseFloat(r.FormValue(prefix+"x"), 64)
		if err != nil {
			return err
		}
		y, err := strconv.ParseFloat(r.FormValue(prefix+"y"), 64)
		if err != nil {
			return err
		}
		if err := db.Model(&FormationPosition{}).Where("id = ? AND formation_id = ?", positionID, formationID).Updates(map[string]interface{}{
			"index_number":  indexNumber,
			"position_name": strings.TrimSpace(r.FormValue(prefix + "name")),
			"x":             x,
			"y":             y,
		}).Error; err != nil {
			return err
		}
	}
	newName := strings.TrimSpace(r.FormValue("newPositionName"))
	if newName != "" {
		playerCount := lineupPlayerCountForTeam(db, formation.TeamID)
		if len(formation.Positions) >= playerCount {
			return fmt.Errorf("team is configured for %s; remove a position before adding another", sideLabel(playerCount))
		}
		indexNumber, _ := strconv.Atoi(r.FormValue("newPositionIndex"))
		x, _ := strconv.ParseFloat(r.FormValue("newPositionX"), 64)
		y, _ := strconv.ParseFloat(r.FormValue("newPositionY"), 64)
		if x == 0 {
			x = 50
		}
		if y == 0 {
			y = 50
		}
		return db.Create(&FormationPosition{FormationID: formationID, IndexNumber: indexNumber, PositionName: newName, X: x, Y: y}).Error
	}
	return nil
}

func formationCopyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "formationId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid formation", http.StatusBadRequest)
			return
		}
		actor := getLineupActor(r, db)
		if actor.TeamID == 0 || (!actor.IsAdmin && actor.UserID == 0) {
			lineupLoginPageForRequest(r, db, "").Render(GetContext(r, db), w)
			return
		}
		formation, err := getFormation(db, uint(id))
		if err != nil {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		if formation.TeamID != actor.TeamID || formation.Status == formationStatusDraft && !actorCanSeeDraft(actor, formation.TeamID, formation.CreatedByUserID) {
			http.Error(w, "Formation not found", http.StatusNotFound)
			return
		}
		_ = r.ParseForm()
		name := r.FormValue("name")
		if strings.TrimSpace(name) == "" {
			name = formation.Name
		}
		copyFormation := Formation{
			TeamID:          actor.TeamID,
			CreatedByUserID: actor.UserID,
			Name:            copyName(name, lineupActorName(db, actor)),
			Status:          formationStatusDraft,
		}
		if len(r.Form["positionIds"]) > 0 {
			for _, idStr := range r.Form["positionIds"] {
				positionID, err := strconv.ParseUint(idStr, 10, 64)
				if err != nil {
					http.Error(w, "Invalid position", http.StatusBadRequest)
					return
				}
				prefix := fmt.Sprintf("position_%d_", positionID)
				indexNumber, err := strconv.Atoi(r.FormValue(prefix + "index"))
				if err != nil {
					http.Error(w, "Invalid position index", http.StatusBadRequest)
					return
				}
				x, err := strconv.ParseFloat(r.FormValue(prefix+"x"), 64)
				if err != nil {
					http.Error(w, "Invalid x coordinate", http.StatusBadRequest)
					return
				}
				y, err := strconv.ParseFloat(r.FormValue(prefix+"y"), 64)
				if err != nil {
					http.Error(w, "Invalid y coordinate", http.StatusBadRequest)
					return
				}
				copyFormation.Positions = append(copyFormation.Positions, FormationPosition{
					IndexNumber:  indexNumber,
					PositionName: strings.TrimSpace(r.FormValue(prefix + "name")),
					X:            x,
					Y:            y,
				})
			}
		} else {
			for _, pos := range formation.Positions {
				copyFormation.Positions = append(copyFormation.Positions, FormationPosition{
					IndexNumber:  pos.IndexNumber,
					PositionName: pos.PositionName,
					X:            pos.X,
					Y:            pos.Y,
				})
			}
		}
		if err := db.Create(&copyFormation).Error; err != nil {
			warning(fmt.Sprintf("Could not copy formation: %v", err)).Render(GetContext(r, db), w)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/formations/%d/edit", copyFormation.ID), http.StatusSeeOther)
	}
}
