package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type MatchPageData struct {
	Match  Match
	isOpen bool
	Msg    string
}

type MatchForm struct {
	MatchId               uint64   `schema:"matchId" validate:"required"`
	Location              string   `schema:"location"`
	StartTime             string   `schema:"startTime"`
	Opponent              string   `schema:"opponent"`
	Subtitle              string   `schema:"subtitle"`
	PlayerOfTheDay        uint64   `schema:"playerOfTheDay"`
	MatchLng              float64  `schema:"matchLng"`
	MatchLat              float64  `schema:"matchLat"`
	LineupID              string   `schema:"lineupId"`
	DudOfTheDay           uint64   `schema:"dudOfTheDay"`
	EventTypeInjury       []uint64 `schema:"eventTypeInjury"`
	EventTypeGoal         []uint64 `schema:"eventTypeGoal"`
	EventTypeAssist       []uint64 `schema:"eventTypeAssist"`
	EventTypeConcededGoal []string `schema:"eventTypeConceded-Goal"`
}

func publicMatchURLPath(match Match) string {
	if strings.TrimSpace(match.MatchURLSlug) == "" {
		return ""
	}
	return fmt.Sprintf("/match-url/%s", match.MatchURLSlug)
}

func publicMatchFeedbackURLPath(match Match) string {
	if strings.TrimSpace(match.MatchURLSlug) == "" {
		return ""
	}
	return fmt.Sprintf("/match-url/%s/feedback", match.MatchURLSlug)
}

func historicMatchTotalPages(total int, limit int) int {
	if limit < 1 {
		limit = 20
	}
	totalPages := (total + limit - 1) / limit
	if totalPages < 1 {
		return 1
	}
	return totalPages
}

func renderMatchesManage(db *gorm.DB, w http.ResponseWriter, r *http.Request, baseUrl string, msg string) {
	ctx := GetContext(r, db)
	seasonId := activeSeasonID(db)
	teamId := getTeamId(ctx)

	matches, err := GetManageMatches(db, teamId, seasonId, 0, 9999)
	if err != nil {
		errMsg("Could not get matches").Render(ctx, w)
		return
	}

	pwfs, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
	if err != nil {
		http.Error(w, "Could not get matches", http.StatusNotFound)
		return
	}

	var team *Team
	if teamId > 0 {
		team, err = GetTeam(db, teamId)
		if err != nil {
			warning(fmt.Sprintf("Could not load active team: %v", err)).Render(ctx, w)
			return
		}
	}

	matchesManage(baseUrl, true, matches, pwfs, team, msg).Render(ctx, w)
}

func activeSeasonID(db *gorm.DB) uint {
	activeSeason, err := GetActiveSeason(db)
	if err != nil || activeSeason == nil {
		return 0
	}
	return activeSeason.ID
}

func matchListHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var isOpen bool = false

		switch r.Method {
		case "GET":
			isOpen = r.URL.Query().Get("isOpen") == "true"
			logAutoSportySync(db, getTeamId(GetContext(r, db)))

			var season uint
			seasonStr := r.URL.Query().Get("season")
			if len(seasonStr) == 0 {
				season = activeSeasonID(db)
			} else {
				seasonUint, err := strconv.ParseInt(seasonStr, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
					return
				}
				season = uint(seasonUint)
			}

			var page = 0
			pageStr := r.URL.Query().Get("page")
			if len(pageStr) == 0 {
				page = 0
			} else {
				pageIdUint, err := strconv.ParseInt(pageStr, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
					return
				}
				page = int(pageIdUint)
			}

			limitStr := r.URL.Query().Get("limit")
			var limit = 50
			if len(limitStr) == 0 {
				limit = 50
			} else {
				limit, err := strconv.ParseInt(limitStr, 10, 64)
				if err != nil || limit == 0 {
					http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
					return
				}
			}

			log.Printf("GetMatches(season%+v, page%+v, limit: %+v)", season, page, limit)
			match, err := GetMatches(db, season, page, limit)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error retrieving match: %v", err), http.StatusInternalServerError)
				return
			}
			resType := r.URL.Query().Get("type")
			matchIdStr := r.URL.Query().Get("matchId")
			isFineMaster := r.URL.Query().Get("isFineMaster") == "true"
			matchId, err := strconv.ParseInt(matchIdStr, 10, 64)
			if err != nil {
				matchId = 0
			}
			if resType == "select" {
				matchComp := matchSelector(match, uint(matchId))
				matchComp.Render(GetContext(r, db), w)
			} else {
				matchComp := matchListPage(match, isOpen, isFineMaster, false)
				matchComp.Render(GetContext(r, db), w)
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func matchHistoryHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := GetContext(r, db)
		teamID := getTeamId(ctx)
		if teamID == 0 {
			warning("No active team").Render(ctx, w)
			return
		}

		page := 1
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			parsedPage, err := strconv.Atoi(pageStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
				return
			}
			page = parsedPage
		}
		if page < 1 {
			page = 1
		}

		limit := 20
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil || parsedLimit < 1 {
				http.Error(w, fmt.Sprintf("Error parsing limit %v", err), http.StatusBadRequest)
				return
			}
			limit = parsedLimit
		}

		matches, total, err := GetHistoricMatches(db, teamID, page, limit, time.Now())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving historic matches: %v", err), http.StatusInternalServerError)
			return
		}
		if page <= 1 {
			historicMatchesPanelContent(matches, page, limit, int(total)).Render(ctx, w)
			return
		}
		historicMatchesPage(matches, page, limit, int(total)).Render(ctx, w)
	}
}

func seasonBulkUpdateHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// get type from url param (not query string)
			updateType := chi.URLParam(r, "updateType")
			switch updateType {
			case "set-season":
				{
					err := r.ParseForm()
					if err != nil {
						warnComp := warning("Invalid form data")
						warnComp.Render(GetContext(r, db), w)
						return
					}
					days := r.FormValue("days")
					daysInt, err := strconv.ParseInt(days, 10, 64)
					if err != nil {
						warnComp := warning("Invalid days")
						warnComp.Render(GetContext(r, db), w)
						return
					}
					finesUpdated, playersUpdated, err := SetSeasonId(db, int(daysInt))

					if err != nil {
						er := warning(fmt.Sprintf("Error parsing form data 2 %+v", err))
						er.Render(GetContext(r, db), w)
					} else {

						success := success(fmt.Sprintf("SetSeasonId days:%d updated %d fine and %d players", daysInt, finesUpdated, playersUpdated))
						success.Render(GetContext(r, db), w)
					}
				}
			case "set-player-subs-outstanding-for-season":
				{
					activeSeason, err := GetActiveSeason(db)
					if err != nil {
						er := warning(fmt.Sprintf("Could not get active season %+v", err))
						er.Render(GetContext(r, db), w)
						return
					}
					seasonSubs := r.FormValue("seasonSubs")
					seasonSubsInt, err := strconv.ParseInt(seasonSubs, 10, 64)
					if err != nil {
						warnComp := warning("Invalid seasonSubs")
						warnComp.Render(GetContext(r, db), w)
						return
					}

					recordsUpdated, err := SetPlayerSubsOutstandingForActiveSeason(db, activeSeason.ID, int(seasonSubsInt))

					if err != nil {
						er := warning(fmt.Sprintf("Error parsing form data 2 %+v", err))
						er.Render(GetContext(r, db), w)
					} else {

						success := success(fmt.Sprintf("SetPlayerSubsOutstandingForSeason updated %d records", recordsUpdated))
						success.Render(GetContext(r, db), w)
					}
				}
			}
		case "GET":
			updateType := chi.URLParam(r, "updateType")
			switch updateType {
			case "set-season":

				activeSeason, err := GetActiveSeason(db)
				if err != nil {
					er := warning(fmt.Sprintf("Could not get active season %+v", err))
					er.Render(GetContext(r, db), w)
					return
				}

				setSeasonForm := setSeasonForm(activeSeason.Title)
				setSeasonForm.Render(GetContext(r, db), w)
				return
			default:
				http.Error(w, "Method Not Allowed (updateType needed)", http.StatusMethodNotAllowed)
				return
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func seasonSpecificHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			seasonIdStr := chi.URLParam(r, "seasonId")
			seasonId, err := strconv.ParseUint(seasonIdStr, 10, 64)

			if err != nil {
				warning := warning(fmt.Sprintf("Error parsing season ID: %v", err))
				warning.Render(GetContext(r, db), w)
				return
			}

			delErr := DeleteSeason(db, uint(seasonId))
			if delErr != nil {
				warning := warning(fmt.Sprintf("Error deleting season: %+v", delErr))
				warning.Render(GetContext(r, db), w)
				return
			}

			success := success(fmt.Sprintf("Deleted season: %d", seasonId))
			success.Render(GetContext(r, db), w)

		}
	}
}

func seasonHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			seasons, err := GetSeasons(db)
			if err != nil {
				http.Error(w, "Could not get seasons", http.StatusNotFound)
				return
			}
			displayType := r.URL.Query().Get("type")
			switch displayType {
			case "selector":
				seasonSelComp := selectSeasons(seasons)
				seasonSelComp.Render(GetContext(r, db), w)
				return
			case "button":
				seasonComp := manageSeasonsButton(false)
				seasonComp.Render(GetContext(r, db), w)
				return
			default:
				activeSeason, err := GetActiveSeason(db)
				if err != nil {
					warnComp := warning("Failed to GetActiveSeason")
					warnComp.Render(GetContext(r, db), w)
					return
				}

				seasonComp := manageSeasons(seasons, activeSeason)
				seasonComp.Render(GetContext(r, db), w)
			}

			return

		case "POST":
			err := r.ParseForm()
			if err != nil {
				warnComp := warning("Invalid form data")
				warnComp.Render(GetContext(r, db), w)
				return
			}

			title := r.FormValue("title")
			if title == "" {
				warnComp := warning("Title is required")
				warnComp.Render(GetContext(r, db), w)
				return
			}
			startDate := r.FormValue("startDate")
			seasonIdStr := r.FormValue("seasonId")
			isActive := r.FormValue("isActive") == "on"

			// Parse start date like 2024-08-01T09:25
			startDateParsed, err := time.Parse("2006-01-02T15:04", startDate)
			if err != nil {
				warnComp := warning("Invalid start date")
				warnComp.Render(GetContext(r, db), w)
				return
			}

			var season *Season
			if seasonIdStr != "" {
				// If seasonId is provided, attempt to parse and update the season
				seasonId, err := strconv.ParseUint(seasonIdStr, 10, 64)
				if err != nil {
					warnComp := warning("Invalid season ID")
					warnComp.Render(GetContext(r, db), w)
					return
				}

				// Fetch the season from the database to update it
				season, err = GetSeasonById(db, uint(seasonId))
				if err != nil {
					warnComp := warning("Season not found")
					warnComp.Render(GetContext(r, db), w)
					return
				}
				// Update the season fields
				season.Title = title
				season.StartDate = startDateParsed
				season.IsActive = isActive
			} else {
				// If seasonId is not provided, create a new season
				season = &Season{
					Title:     title,
					StartDate: startDateParsed,
					IsActive:  isActive,
				}
			}

			log.Printf("Season %+v", season)
			// Save the season (either update or create)
			savErr := SaveSeason(db, season)
			if savErr != nil {
				warnComp := warning("Could not save season")
				warnComp.Render(GetContext(r, db), w)
				return
			}
			if len(seasonIdStr) > 0 {
				successComp := success("Season updated")
				successComp.Render(GetContext(r, db), w)
			} else {
				successComp := successWithReload("Season created")
				successComp.Render(GetContext(r, db), w)
			}

			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func matchSportySyncHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		requestedTeamID, err := parseOptionalUint(r.URL.Query().Get("teamId"))
		if err != nil {
			warning("Invalid team ID").Render(ctx, w)
			return
		}
		if requestedTeamID > 0 {
			teamId = requestedTeamID
		}
		if teamId == 0 {
			warning("No active team - Sporty sync is unavailable").Render(ctx, w)
			return
		}

		team, err := GetTeam(db, teamId)
		if err != nil || team == nil {
			warning(fmt.Sprintf("Team not found (teamId=%d)", teamId)).Render(ctx, w)
			return
		}
		if !teamHasSportyConfig(*team) {
			warning("Sporty IDs are missing for this team").Render(ctx, w)
			return
		}

		result, err := syncSportyMatches(db, team, nil, false, includePastSportyMatches(r))
		if err != nil {
			renderMatchesManage(db, w, r, "/finemaster", fmt.Sprintf("Sporty sync failed: %v", err))
			return
		}

		renderMatchesManage(db, w, r, "/finemaster", sportySyncSummary(result))
	}
}

func matchHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	var matchId uint
	var err error
	return func(w http.ResponseWriter, r *http.Request) {
		var successMsg = ""

		switch r.Method {
		case "GET":

			matchIdStr := chi.URLParam(r, "matchId")

			if matchIdStr == "" {
				renderMatchesManage(db, w, r, r.Header.Get("Referrer"), "")
				return
			}
			var matchId64 uint64
			matchId64, err = strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
				return
			}

			matchId = uint(matchId64)

			match, err := GetMatchMetaGeneral(db, matchId)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error retrieving match: %v", err), http.StatusInternalServerError)
				return
			}

			var url = templ.SafeURL(r.Header.Get("Referrer"))

			var renderType = r.URL.Query().Get("type")
			actor := getLineupActor(r, db)
			notesData, notesErr := contextualNotesData(db, getTeamId(GetContext(r, db)), noteTargetMatch, match.Match.ID, "vs "+match.Match.Opponent, fmt.Sprintf("/match/%d", match.Match.ID))
			if notesErr != nil {
				successMsg = strings.TrimSpace(successMsg + " " + fmt.Sprintf("Could not load model notes: %v", notesErr))
			}
			notesData.ShowAdminNotes = actor.IsAdmin

			log.Printf("MatchHanlder - %s", renderType)
			switch renderType {
			case "form":

				matchComp := editMatch(url, *match, successMsg, notesData)
				matchComp.Render(GetContext(r, db), w)
				return
			case "list":
				matchComp := viewMatch(*match)
				matchComp.Render(GetContext(r, db), w)
				return
			default:
				matchComp := editMatchContainer(url, *match, successMsg, notesData)
				matchComp.Render(GetContext(r, db), w)
				return

			}

		case "POST":
			// Simplified example: Assume the response after a POST is a success message or error
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			seasonId := activeSeasonID(db)
			matchIdStr := r.FormValue("matchId")
			matchAction := strings.TrimSpace(r.FormValue("matchAction"))
			var matchId64 uint64 = 0
			if len(matchIdStr) > 0 {
				matchId64, err = strconv.ParseUint(matchIdStr, 10, 64)
				if err != nil {
					var msg = fmt.Sprintf("Error parsing match ID: (%s) %v", matchIdStr, err)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
				}
				if matchAction == "quick-fix-season-team" {
					teamID := getTeamId(GetContext(r, db))
					if teamID == 0 || seasonId == 0 {
						errComp := errMsg("Current season and active team are required to fix this match")
						errComp.Render(GetContext(r, db), w)
						return
					}
					if err := db.Model(&Match{}).Where("id = ?", uint(matchId64)).Updates(map[string]interface{}{
						"season_id": seasonId,
						"team_id":   teamID,
					}).Error; err != nil {
						errComp := errMsg(fmt.Sprintf("Could not fix match season/team: %v", err))
						errComp.Render(GetContext(r, db), w)
						return
					}
					successMsg = fmt.Sprintf("Match set to current season %d and team %d.", seasonId, teamID)
					var url = templ.SafeURL(r.Header.Get("Referrer"))
					genMeta, err := GetMatchMetaGeneral(db, uint(matchId64))
					if err != nil || genMeta == nil {
						errComp := errMsg(fmt.Sprintf("GetMatchMetaGeneral failed %v", err))
						errComp.Render(GetContext(r, db), w)
						return
					}
					actor := getLineupActor(r, db)
					notesData, notesErr := contextualNotesData(db, teamID, noteTargetMatch, genMeta.Match.ID, "vs "+genMeta.Match.Opponent, fmt.Sprintf("/match/%d", genMeta.Match.ID))
					if notesErr != nil {
						successMsg = strings.TrimSpace(successMsg + " " + fmt.Sprintf("Could not load model notes: %v", notesErr))
					}
					notesData.ShowAdminNotes = actor.IsAdmin
					editMatch(url, *genMeta, successMsg, notesData).Render(GetContext(r, db), w)
					return
				}

				if err := r.ParseForm(); err != nil {
					http.Error(w, "Invalid form data", http.StatusBadRequest)
					return
				}

				var form MatchForm
				if err := decoder.Decode(&form, r.PostForm); err != nil {
					var msg = fmt.Sprintf("Error parsing MatchForm: %v", err.Error())
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
				}

				/*	form = NewMatchForm{
					MatchId: matchId64,
					Location:  r.FormValue("location"),
					StartTime: r.FormValue("startTime"),
					Opponent:  r.FormValue("opponent"),
					Subtitle:  r.FormValue("subtitle"),
					EventTypeInjury: r.FormValue("eventTypeInjury"),
				}*/

				playerOfDayStr := r.FormValue("playerOfTheDay")
				if len(playerOfDayStr) > 0 {
					playerOfDayId, err := strconv.ParseUint(playerOfDayStr, 10, 64)
					if err != nil {
						var msg = fmt.Sprintf("Error parsing playerOfTheDay ID: \"%s\" %v", playerOfDayStr, err)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
					}

					form.PlayerOfTheDay = playerOfDayId
				}

				dudOfDayStr := r.FormValue("dudOfTheDay")
				if len(dudOfDayStr) > 0 {
					dudOfDayId, err := strconv.ParseUint(dudOfDayStr, 10, 64)
					if err != nil {
						var msg = fmt.Sprintf("Error parsing dudOfTheDay ID: %v", err)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
					}

					form.DudOfTheDay = dudOfDayId
				}
				log.Printf("\nUPDATE UPDATE UPDATE  %+v\n", form)

				// Create a new match based on the form data
				var existingMatch *Match
				existingMatch, err = GetMatch(db, uint(matchId64))
				if err != nil {
					var msg = fmt.Sprintf("GetMatch failed %v", err)
					log.Print(msg)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
					return
				}
				match := Match{
					TeamID:          existingMatch.TeamID,
					LineupID:        existingMatch.LineupID,
					Location:        form.Location,
					Opponent:        form.Opponent,
					Subtitle:        form.Subtitle,
					SeasonId:        uint64(seasonId),
					PlayerOfTheDay:  uint(form.PlayerOfTheDay),
					DudOfTheDay:     uint(form.DudOfTheDay),
					MatchLat:        form.MatchLat,
					MatchLng:        form.MatchLng,
					SportyFixtureID: existingMatch.SportyFixtureID,
				}
				var selectedLineup *Lineup
				if r.PostForm.Has("lineupId") && strings.TrimSpace(form.LineupID) != "" {
					lineupID, err := strconv.ParseUint(form.LineupID, 10, 64)
					if err != nil {
						var msg = fmt.Sprintf("Error parsing line-up ID: %v", err)
						log.Print(msg)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
						return
					}
					var lineup Lineup
					if err := db.Where("id = ? AND team_id = ? AND status = ?", uint(lineupID), existingMatch.TeamID, lineupStatusLive).First(&lineup).Error; err != nil {
						var msg = fmt.Sprintf("Selected line-up is not live for this team: %v", err)
						log.Print(msg)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
						return
					}
					selectedLineup = &lineup
					match.LineupID = lineup.ID
				}

				if len(form.StartTime) > 0 {
					// Parse start time
					startTime, err := time.Parse("2006-01-02T15:04", form.StartTime)
					if err != nil {
						var msg = fmt.Sprintf("Invalid start time format (%s) %v", form.StartTime, err)
						log.Print(msg)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
					} else {
						match.StartTime = &startTime

						fines, err := GetFinesByMatchId(db, uint(matchId64))
						if err != nil {
							var msg = fmt.Sprintf("GetFinesByMatchId failed %v", err)
							log.Print(msg)
							errComp := errMsg(msg)
							errComp.Render(GetContext(r, db), w)
						}

						for _, fine := range fines {
							err = SetFineStartTime(db, fine.ID, startTime)
							if err != nil {
								var msg = fmt.Sprintf("SetFineStartTime failed %v", err)
								log.Print(msg)
								errComp := errMsg(msg)
								errComp.Render(GetContext(r, db), w)
							}
						}
					}
				}

				if form.MatchId > 0 {
					match.ID = uint(form.MatchId)
				}
				log.Printf("SaveMatch %+v", match)
				matchId, err = SaveMatch(db, &match)
				if err != nil {
					var msg = fmt.Sprintf("SaveMatch 1 failed %v", err)
					log.Print(msg)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
					return
				}
				if selectedLineup != nil {
					if err := syncLineupMatch(db, *selectedLineup, matchId); err != nil {
						var msg = fmt.Sprintf("Could not link line-up to match: %v", err)
						log.Print(msg)
						errComp := errMsg(msg)
						errComp.Render(GetContext(r, db), w)
						return
					}
				}
				unavailablePlayerIDs, err := parseUnavailablePlayerIDs(r)
				if err != nil {
					var msg = fmt.Sprintf("Error parsing unavailable players: %v", err)
					log.Print(msg)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
					return
				}
				if err := SaveUnavailablePlayersForMatch(db, existingMatch.TeamID, matchId, unavailablePlayerIDs); err != nil {
					var msg = fmt.Sprintf("SaveUnavailablePlayersForMatch failed %v", err)
					log.Print(msg)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
					return
				}

				log.Printf("\n\nFORM: %+v", form)

				if len(form.EventTypeInjury) > 0 {
					for _, injPlayerId := range form.EventTypeInjury {
						err = SaveMatchEvent(db, &MatchEvent{
							MatchId:   uint64(matchId),
							EventName: "",
							EventType: "injury",
							PlayerId:  uint(injPlayerId),
						})
						if err != nil {
							log.Printf("\n\n\nINJURYINJURYINJURY ERROR %v", err)
						}

					}
				}

				if len(form.EventTypeGoal) > 0 {
					for _, injPlayerId := range form.EventTypeGoal {
						err = SaveMatchEvent(db, &MatchEvent{
							MatchId:   uint64(matchId),
							EventName: "",
							EventType: "goal",
							PlayerId:  uint(injPlayerId),
						})
						if err != nil {
							log.Printf("\n\n\nGOAL ERROR %v", err)
						}

					}
				}

				if len(form.EventTypeAssist) > 0 {
					for _, injPlayerId := range form.EventTypeAssist {
						err = SaveMatchEvent(db, &MatchEvent{
							MatchId:   uint64(matchId),
							EventName: "",
							EventType: "assist",
							PlayerId:  uint(injPlayerId),
						})
						if err != nil {
							log.Printf("\n\n\nGOAL ERROR %v", err)
						}

					}
				}

				if len(form.EventTypeConcededGoal) > 0 {
					log.Printf("\n\n\nIN HERE!!!!!!\n\n\n")
					for range form.EventTypeConcededGoal {
						err = SaveMatchEvent(db, &MatchEvent{
							MatchId:   uint64(matchId),
							EventName: "",
							EventType: "conceded-goal",
							PlayerId:  0,
						})
						if err != nil {
							log.Printf("\n\n\nGOAL ERROR %v", err)
						} else {
							log.Printf("\n\n\nconceded GOAL GOAL saved matchId:%d %s", matchId, "injury")
						}

					}
				}

				log.Printf("\n\nMATCH: %+v", match)
				successMsg = fmt.Sprintf("match %s updated (%d) [potd:%s dud:%s]", form.Location, matchId64, playerOfDayStr, dudOfDayStr)

			} else {
				log.Printf("\nCREATE CREATE CREATE \n")
				createForm := MatchForm{
					MatchId:   0,
					Location:  r.FormValue("location"),
					StartTime: r.FormValue("startTime"),
					Opponent:  r.FormValue("opponent"),
					Subtitle:  r.FormValue("subtitle"),
				}

				startTimeTime, startErr := time.Parse("2006-01-02T15:04", createForm.StartTime)
				if startErr != nil {
					errComp := errMsg("Could not parse time string")
					errComp.Render(GetContext(r, db), w)
					return
				}
				log.Printf("SaveMatch CREATE %+v", createForm)
				matchId, err = SaveMatch(db, &Match{
					TeamID:    getTeamId(GetContext(r, db)),
					Location:  createForm.Location,
					Opponent:  createForm.Opponent,
					Subtitle:  createForm.Subtitle,
					StartTime: &startTimeTime,
					SeasonId:  uint64(seasonId),
				})
				if err != nil {
					var msg = fmt.Sprintf("SaveMatch 2 failed %v", err)
					log.Print(msg)
					errComp := errMsg(msg)
					errComp.Render(GetContext(r, db), w)
					return
				}
				successMsg = fmt.Sprintf("New match created (%d)", matchId)

			}

			var url = templ.SafeURL(r.Header.Get("Referrer"))
			genMeta, err := GetMatchMetaGeneral(db, uint(matchId))
			if err != nil || genMeta == nil {
				var msg = fmt.Sprintf("GetMatchMetaGeneral failed %v", err)
				log.Print(msg)
				errComp := errMsg(msg)
				errComp.Render(GetContext(r, db), w)
			}

			actor := getLineupActor(r, db)
			notesData, notesErr := contextualNotesData(db, getTeamId(GetContext(r, db)), noteTargetMatch, genMeta.Match.ID, "vs "+genMeta.Match.Opponent, fmt.Sprintf("/match/%d", genMeta.Match.ID))
			if notesErr != nil {
				successMsg = strings.TrimSpace(successMsg + " " + fmt.Sprintf("Could not load model notes: %v", notesErr))
			}
			notesData.ShowAdminNotes = actor.IsAdmin
			matchComp := editMatch(url, *genMeta, successMsg, notesData)
			matchComp.Render(GetContext(r, db), w)
			return

		case "DELETE":
			matchIdStr := chi.URLParam(r, "matchId")
			matchId64, err := strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
				return
			}

			if matchIdStr != "" {
				err := DeleteMatch(db, uint(matchId64))
				if err != nil {
					errComp := errMsg("Could not DeleteMatch")
					errComp.Render(GetContext(r, db), w)
				}

			}
			matchComp := success(F("deleted match: %d", matchId64))
			matchComp.Render(GetContext(r, db), w)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GetMatchAndEvents(db *gorm.DB, matchId uint) (*MatchState, []MatchEvent, error) {
	matchEvents, getFErr := GetMatchEvents(db, matchId)
	if getFErr != nil {
		return nil, []MatchEvent{}, getFErr
	}

	players, err := GetPlayers(db, 1, 9999)
	if err != nil {
		log.Printf("Invalid GetPlayers")
		return nil, []MatchEvent{}, err
	}

	initialState := MatchState{MatchID: matchId, PlayersOn: []PlayerState{}, ScoreFor: 0, ScoreAgainst: 0}
	matchState := ConstructMatchState(matchEvents, initialState, 90, players)

	return &matchState, matchEvents, nil

}

func ConstructMatchState(events []MatchEvent, currentState MatchState, currentTime int, allPlayers []Player) MatchState {
	if len(events) == 0 {
		return currentState
	}

	event := events[0]
	var player Player
	if event.PlayerId > 0 {

		for _, p := range allPlayers {
			if p.ID == uint(event.PlayerId) {
				player = p
			}
		}
	}

	updatedState := UpdateStateBasedOnEvent(currentState, event, currentTime, &player)

	return ConstructMatchState(events[1:], updatedState, currentTime, allPlayers)
}

func UpdateStateBasedOnEvent(currentState MatchState, event MatchEvent, currentTime int, player *Player) MatchState {
	switch event.EventType {
	case "subbed-on":
		currentState.PlayersOn = append(currentState.PlayersOn, PlayerState{PlayerName: player.Name, PlayerId: player.ID, TimePlayed: currentTime - event.EventMinute})
	case "subbed-off":
		// Logic to remove player and update time played
	case "goal":
		// Increment player's goal count and the total score
		currentState.ScoreFor += 1
	case "assist":
		// Increment player's assist count
	case "own-goal":
		currentState.ScoreAgainst += 1
	case "conceded-goal":
		currentState.ScoreAgainst += 1
	case "attended-training":
		currentState.TrainingTotalNumbers += 1
	}

	// Update players' time played if they are on the field
	for i, player := range currentState.PlayersOn {
		if player.PlayerId == event.PlayerId {
			if event.EventType == "subbed-off" {
				// Remove player or update time
			} else if event.EventType == "goal" {
				currentState.PlayersOn[i].Goals += 1
			} else if event.EventType == "assist" {
				currentState.PlayersOn[i].Assists += 1
			}
		}
	}

	return currentState
}
