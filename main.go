package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var decoder = schema.NewDecoder()

type Context struct {
	Title                string
	UseRoles             bool
	UseMatchEventTracker bool
}

func playerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		var activeSeasonId uint = 0
		activeSeason, err := GetActiveSeason(db)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else if activeSeason != nil {
			activeSeasonId = activeSeason.ID
		}
		switch r.Method {
		case "GET":
			displayType := r.URL.Query().Get("type")

			if displayType == "role-selector" {
				playerIdStr := r.URL.Query().Get("playerId")
				playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("role-selector Error playerHandler playerId %v", err), http.StatusBadRequest)
					return
				}
				if playerId == 0 {
					http.Error(w, fmt.Sprintf("role-selector Error playerHandler playerId %v", err), http.StatusBadRequest)
					return
				}

				playerIds := []uint64{
					playerId,
				}
				playersWithFines, err := GetPlayersWithFines(db, activeSeasonId, teamId, playerIds)
				if err != nil || len(playersWithFines) == 0 {
					log.Printf("Error fetching players with fines: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				playerList := playerRoleSelector(playersWithFines[0], config, nil)
				playerList.Render(ctx, w)
				return
			} else if displayType == "super-input" {

				playerIdStr := r.URL.Query().Get("playerId")
				inputType := r.URL.Query().Get("inputType")

				var playerId uint64
				var err error
				if len(playerIdStr) > 0 {
					playerId, err = strconv.ParseUint(playerIdStr, 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("super-input - Error playerHandler playerId \"%s\" [%d] %v", playerIdStr, len(playerIdStr), err), http.StatusBadRequest)
						return
					}
				}
				players, err := GetPlayers(db, 0, 999)
				if err != nil || len(players) == 0 {
					log.Printf("Error fetching players with fines: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				playerList := playerInputSelector(players, uint(playerId), inputType)
				playerList.Render(ctx, w)
			} else {
				warnings := warning("Playerhanlder - Method not allowed")
				warnings.Render(ctx, w)
				return
			}

		case "POST":
			{
				if err := r.ParseForm(); err != nil {
					log.Printf("Error parsing form data in playerHandler POST: %v", err)
					http.Error(w, "Bad Request savePlayerHandler ParseForm", http.StatusBadRequest)
					return
				}
				playerAction := r.FormValue("playerAction")
				delete(r.PostForm, "playerAction")

				var player Player
				// Use the decoder to populate the player struct
				if err := decoder.Decode(&player, r.PostForm); err != nil {
					log.Printf("Error decoding form into player struct: %v", err)
					http.Error(w, "Bad Request savePlayerHandler Decode", http.StatusBadRequest)
					return
				}
				if player.ID > 0 {
					existingPlayer, err := GetPlayerByID(db, player.ID)
					if err != nil {
						log.Printf("Error fetching player before update: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
					player.Active = existingPlayer.Active
				}
				switch playerAction {
				case "activate":
					player.Active = true
				case "deactivate":
					player.Active = false
				case "":
					if activeValues, ok := r.PostForm["Active"]; ok && len(activeValues) > 0 {
						active, err := strconv.ParseBool(activeValues[len(activeValues)-1])
						if err != nil {
							http.Error(w, "Bad Request savePlayerHandler Active", http.StatusBadRequest)
							return
						}
						player.Active = active
					}
				default:
					http.Error(w, "Bad Request savePlayerHandler playerAction", http.StatusBadRequest)
					return
				}

				// Now, `player` is populated with values from the form
				// Save the player using your existing logic
				if err := SavePlayer(db, &player); err != nil {
					log.Printf("Error saving player: %v", err)
					http.Error(w, "Internal Server Error 1", http.StatusInternalServerError)
					return
				}

				//		w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))

				// Optionally, you can set the status code to 200 OK or any appropriate status
				playerIds := []uint64{
					uint64(player.ID),
				}

				playersWithFines, err := GetPlayersWithFines(db, activeSeasonId, teamId, playerIds)
				if err != nil || len(playersWithFines) == 0 {
					log.Printf("Error fetching players with fines: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				playerList := playerRoleSelector(playersWithFines[0], config, successWithReload("Updated player"))
				playerList.Render(ctx, w)
				return

			}
		case "DELETE":
			{
				log.Println("PlayerHanlder")
				if err := r.ParseForm(); err != nil {
					log.Printf("Error parsing form data in playerHanlder DELETE: %v", err)
					http.Error(w, "Bad Request savePlayerHandler ParseForm", http.StatusBadRequest)
					return
				}

				playerIdStr := r.URL.Query().Get("playerId")
				playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
				if err != nil {
					log.Printf("Error get playerIdStr fine: %v", err)
					http.Error(w, "Internal playerIdStr Error", http.StatusInternalServerError)
					return
				}

				log.Printf("PlayerHanlder %d", playerId)

				// Now, `player` is populated with values from the form
				// Save the player using your existing logic
				if err := DeletePlayer(db, uint(playerId)); err != nil {
					log.Printf("Error saving player: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))

				// Optionally, you can set the status code to 200 OK or any appropriate status
				w.WriteHeader(http.StatusOK)
				return
			}
		}

	}
}
func fineContestHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "fineContestHandler - Invalid form data", http.StatusBadRequest)
			return
		}

		// Extract and validate form data
		fineID, err := strconv.ParseUint(r.FormValue("fid"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid fid - %s", r.FormValue("fid")), http.StatusBadRequest)
			return
		}
		contest := r.FormValue("contest")

		err = UpdateFineContestByID(db, uint(fineID), contest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid UpdateFineContestByID - %v", err), http.StatusBadRequest)
			return
		}

		success := success("Added Contest")
		success.Render(GetContext(r, db), w)
	}
}

func fineContextHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "fineContextHandler - Invalid form data", http.StatusBadRequest)
			return
		}

		// Extract and validate form data
		fineID, err := strconv.ParseUint(r.FormValue("fid"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid fid - %s - %v", r.FormValue("fid"), err), http.StatusBadRequest)
			return
		}

		matchIdStr := r.FormValue("matchId")
		var matchId uint64 = 0
		var fineAt *time.Time
		if matchIdStr != "NA" {
			matchId, err = strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid matchId ID - %s", matchIdStr), http.StatusBadRequest)
				return
			}

			match, getMatchErr := GetMatch(db, uint(matchId))
			if getMatchErr != nil {
				http.Error(w, fmt.Sprintf("Invalid matchId ID - %s", matchIdStr), http.StatusBadRequest)
				return
			} else {
				fineAt = match.StartTime
			}
		}

		context := r.FormValue("context")

		err = UpdateFineContextByID(db, uint(fineID), uint(matchId), context, fineAt)
		log.Printf("fineContextHanlder-3")
		if err != nil {
			http.Error(w, "Invalid UpdateFineContextByID ID", http.StatusBadRequest)
			return
		}

		success := contextSuccess(uint(matchId), context, fineAt)
		success.Render(GetContext(r, db), w)
	}
}

func courtHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				full := r.URL.Query().Get("full") == "true"

				if full {
					courtComp := courtManage(true)
					courtComp.Render(GetContext(r, db), w)
					return
				}

				courtComp := courtManage(false)
				courtComp.Render(GetContext(r, db), w)
				return
			}
		default:
			warning := warning("Method not allowed")
			warning.Render(GetContext(r, db), w)
			return
		}
	}
}

func fineSetCourtSessionOrderHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "fineSetCourtSessionOrderHandler - Invalid form data", http.StatusBadRequest)
			return
		}

		// Extract and validate form data
		fineID, err := strconv.ParseUint(r.FormValue("fid"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid fid - %s", r.FormValue("fid")), http.StatusBadRequest)
			return
		}
		courtSessionOrderInt, err := strconv.ParseUint(r.FormValue("courtSessionOrder"), 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid courtSessionOrder - %v", err), http.StatusBadRequest)
			return
		}
		courtSessionNotes := r.FormValue("courtSessionNote")
		log.Printf("fineSetCourtSessionOrderHandler %d courtSessionNote:[%s]", courtSessionOrderInt, courtSessionNotes)

		err = UpdateFineCourtSessionOrderByID(db, uint(fineID), uint(courtSessionOrderInt), courtSessionNotes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid UpdateFineCourtSessionOrderByID - %v", err), http.StatusBadRequest)
			return
		}

		success := success(fmt.Sprintf("Added Court Session Order #%d to fine %d", courtSessionOrderInt, fineID))
		success.Render(GetContext(r, db), w)
	}

}

func fineEditHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// Extract fine ID from query parameters
			fineIDStr := chi.URLParam(r, "fid")
			fineID, err := strconv.ParseUint(fineIDStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid fine ID - %s", fineIDStr), http.StatusBadRequest)
				return
			}

			// Fetch the fine details from the database
			fine, err := GetFineByID(db, uint(fineID))
			if err != nil {
				http.Error(w, "GetFineByID not found", http.StatusNotFound)
				return
			}

			// Fetch associated player details (assuming GetPlayerByID is a function you have)
			player, err := GetPlayerByID(db, fine.PlayerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("GetPlayerByID Player not found - %d", fine.PlayerID), http.StatusNotFound)
				return
			}

			match := &Match{}
			if fine.MatchId > 0 {
				match, err = GetMatch(db, uint(fine.MatchId))
				if err != nil {
					http.Error(w, fmt.Sprintf("GetMatch Match not found - %d", fine.MatchId), http.StatusNotFound)
					return

				} else if match.StartTime != nil {
					fine.FineAt = *match.StartTime
				} else {
					fine.FineAt = time.Now()
				}

			}

			hideFineImageFeature := true

			// Prepare the data for rendering
			fineWithPlayer := FineWithPlayer{
				Fine:   *fine,
				Player: *player,
				Match:  *match,
			}

			isEdit := r.URL.Query().Get("isEdit")
			isContest := r.URL.Query().Get("isContest")
			isContext := r.URL.Query().Get("isContext")

			//	isFineMasterStr := r.URL.Query().Get("isFineMaster")

			/*var isFineMaster bool = isFineMasterStr == "true"
			if isFineMaster {
				log.Printf("permision Error fetching fineEditHandler")
				http.Error(w, "Not this time mate.", http.StatusInternalServerError)
				return
			}*/

			var isFineMaster bool = true

			if isEdit == "true" {
				fineEditRow := fineEditRow(fineWithPlayer, isFineMaster, hideFineImageFeature)
				fineEditRow.Render(GetContext(r, db), w)
			} else if isEdit == "fineEditDiv" {
				fineEditRow := fineEditDiv(fineWithPlayer, isFineMaster, hideFineImageFeature)
				fineEditRow.Render(GetContext(r, db), w)
			} else if isEdit == "form" {
				fineFormRow := fineEditForm(fineWithPlayer, isFineMaster, hideFineImageFeature)
				fineFormRow.Render(GetContext(r, db), w)
			} else if isContest == "true" {

				fineContestRow := fineContestRow(fineWithPlayer)
				fineContestRow.Render(GetContext(r, db), w)
			} else if isContext == "true" {
				matches, err := GetMatches(db, activeSeasonID(db), 0, 9999)
				if err != nil {
					http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
					return
				}
				fineContestRow := fineContextRow(fineWithPlayer, matches)
				fineContestRow.Render(GetContext(r, db), w)
			} else {
				fineRowComp := fineRow(isFineMaster, fineWithPlayer)
				fineRowComp.Render(GetContext(r, db), w)
			}
			return
		case "POST":
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "fineEditHandler - Invalid form data", http.StatusBadRequest)
				return
			}

			// Extract and validate form data
			fineID, err := strconv.ParseUint(r.FormValue("fid"), 10, 64)
			if err != nil {
				http.Error(w, "fineEditHandler - Invalid fine ID", http.StatusBadRequest)
				return
			}

			amountStr := r.FormValue("amount")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				http.Error(w, "Invalid amount", http.StatusBadRequest)
				return
			}

			reason := r.FormValue("reason")
			approvedStr := r.FormValue("approved")
			var approved bool
			if len(approvedStr) == 0 {
				approved = config.DefaultToApproved
			} else {
				approved = approvedStr == "true" //todo: on?
			}

			playerId, err := strconv.ParseUint(r.FormValue("playerId"), 10, 64)
			if err != nil {
				errComp := errMsg(F("Invalid playerId", r.FormValue("playerId")))
				errComp.Render(GetContext(r, db), w)
			}

			activeSeason, err := GetActiveSeason(db)
			var seasonId uint
			if err != nil {
				errComp := errMsg("Cannot select get the active season?")
				errComp.Render(GetContext(r, db), w)
			} else if activeSeason != nil {
				seasonId = activeSeason.ID
			}

			matchId, err := strconv.ParseUint(r.FormValue("matchId"), 10, 64)
			if err != nil {
				errComp := errMsg(F("Invalid matchId", r.FormValue("matchId")))
				errComp.Render(GetContext(r, db), w)
			}

			context := r.FormValue("context")

			// Update the fine in the database
			fine := Fine{
				Model:    gorm.Model{ID: uint(fineID)},
				PlayerID: uint(playerId),
				SeasonID: seasonId,
				Amount:   amount,
				Reason:   reason,
				Context:  context,
				Approved: approved,
			}

			if matchId != 0 {
				match, err := GetMatch(db, uint(matchId))
				if err != nil {
					errComp := errMsg("Cannot select get the active match?")
					errComp.Render(GetContext(r, db), w)
				} else if match != nil {
					fine.FineAt = *match.StartTime
				}
			} else {
				activeMatch, actMatchERr := GetActiveMatch(db)
				if actMatchERr != nil {
					log.Printf("fineEditHandler - Error fetching active match: %v", actMatchERr)

					errComp := errMsg("Cannot select get the active match?")
					errComp.Render(GetContext(r, db), w)
				} else if activeMatch != nil {
					fine.FineAt = *activeMatch.StartTime
				}
			}

			if err := SaveFine(db, &fine); err != nil {
				errComp := errMsg("Failed to update fine")
				errComp.Render(GetContext(r, db), w)
				return
			} else {
				success := success("Fine updated")
				success.Render(GetContext(r, db), w)
				log.Printf("\n\nSAve fine %+v\n\n", fine)
			}

			player, err := GetPlayerByID(db, fine.PlayerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
				return
			}

			match := &Match{}
			if fine.MatchId > 0 {
				match, _ = GetMatch(db, uint(fine.MatchId))
			}

			// Prepare the data for rendering
			fineWithPlayer := FineWithPlayer{
				Fine:   fine,
				Player: *player,
				Match:  *match,
			}

			fineRowComp := fineRow(true, fineWithPlayer)
			fineRowComp.Render(GetContext(r, db), w)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func fineImageHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := ""
		switch r.Method {
		case "POST":
			{

				fineIDParam := chi.URLParam(r, "fid")
				fineID, err := strconv.Atoi(fineIDParam)
				if err != nil {
					http.Error(w, "Invalid fine ID", http.StatusBadRequest)
					return
				}

				file, header, err := r.FormFile("image")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer file.Close()

				filename := header.Filename
				imageData, err := ioutil.ReadAll(file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Find the fine
				var fine Fine
				if err := db.First(&fine, fineID).Error; err != nil {
					http.Error(w, "Fine not found", http.StatusNotFound)
					return
				}

				// Create the image and associate it with the fine
				image := FineImage{
					Filename: filename,
					Data:     imageData,
					FineID:   uint(fineID),
				}

				if result := db.Create(&image); result.Error != nil {
					http.Error(w, result.Error.Error(), http.StatusInternalServerError)
					return
				}
				message = fmt.Sprintf("Image uploaded successfully: %s", filename)
				//	warningComp := warning("upload success")
				//	warningComp.Render(GetContext(r, db), w)
			}
		case "GET":
			{
				fineIDParam := chi.URLParam(r, "fid")
				fineID, err := strconv.Atoi(fineIDParam)

				if err != nil {
					errorComp := warning(fmt.Sprintf("Invalid fine ID: %v", err))
					errorComp.Render(GetContext(r, db), w)
				}

				fineImgs, err := GetFineImages(db, uint(fineID))
				if err != nil {
					errorComp := warning(fmt.Sprintf("Invalid fine ID: %v", err))
					errorComp.Render(GetContext(r, db), w)
				}

				displayType := r.URL.Query().Get("displayType")

				switch displayType {
				case "bigFineImage":
					fineImgComp := bigFineImages(fineImgs, uint(fineID))
					fineImgComp.Render(GetContext(r, db), w)
				default:
				case "fineImage":
					fineImgComp := fineImages(fineImgs, uint(fineID), message)
					fineImgComp.Render(GetContext(r, db), w)
				}

			}
		}
	}
}

type PlayerFinesTotal struct {
	PaymentDifference float64
	TotalFines        float64
	TotalPayments     float64
	PlayerPayment     []PlayerPayment
	Player            PlayerWithFines
}

func fineSummaryHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				ctx := GetContext(r, db)
				teamId := getTeamId(ctx)
				activeSeason, err := GetActiveSeason(db)
				if err != nil {
					http.Error(w, "Failed when getting active season", http.StatusNotFound)
					return
				}
				seasonId := uint(0)
				if activeSeason != nil {
					seasonId = activeSeason.ID
				}

				viewMode := r.URL.Query().Get("viewMode")
				if viewMode == "button" {
					fineList := fineSummaryButton("Open Total View", "summary", false)
					fineList.Render(ctx, w)
					return
				}

				if seasonId == 0 || teamId == 0 {
					warning("No active season/team - totals are unavailable").Render(ctx, w)
					return
				}

				playersWithFines, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
				if err != nil {
					log.Printf("Error fetching players with fines: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				playerPayments, err := GetPlayerPayments(db, seasonId, teamId)
				if err != nil {
					log.Printf("Error fetching player payments: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				var playerFinesTotals []PlayerFinesTotal

				var grandTotal float64

				for _, player := range playersWithFines {
					totalFines := 0.0
					totalPayments := 0.0
					thisPlayersPayments := []PlayerPayment{}
					for _, fine := range player.Fines {
						totalFines += fine.Amount
						grandTotal += fine.Amount
					}

					for _, payment := range playerPayments {
						if payment.PlayerID == player.ID {
							totalPayments += payment.Amount
							thisPlayersPayments = append(thisPlayersPayments, payment)
						}

					}

					// Store the player's total fines
					playerFinesTotals = append(playerFinesTotals, PlayerFinesTotal{
						Player:            player,
						TotalFines:        totalFines,
						TotalPayments:     totalPayments,
						PlayerPayment:     thisPlayersPayments,
						PaymentDifference: totalFines - totalPayments,
					})
				}

				fineSummary := fineTotals(playerFinesTotals, grandTotal)
				fineSummary.Render(ctx, w)
			}
		case "POST":
			{

				fTot, pfTot, err := SetAllFineAmounts(db, 2.0)

				if err != nil {
					er := warning(fmt.Sprintf("Error parsing form data 2 %+v", err))
					er.Render(GetContext(r, db), w)
				} else {

					success := success("Fine amounts updated fines: " + fmt.Sprintf("%d", fTot) + " preset fines: " + fmt.Sprintf("%d", pfTot))
					success.Render(GetContext(r, db), w)
				}
			}
		default:
			{
				er := warning("Method not allowed")
				er.Render(GetContext(r, db), w)
			}
		}
	}
}

type FineWithPlayer struct {
	Fine   Fine
	Player Player
	Match  Match
}

func fineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		if user, ok := requireGoogleAdminForTeam(db, r, teamId); !ok {
			if user == nil {
				w.Header().Set("HX-Redirect", "/finemaster/auth")
				http.Redirect(w, r, "/finemaster/auth", http.StatusSeeOther)
			} else {
				http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
			}
			return
		}
		activeSeason, _ := GetActiveSeason(db)
		seasonId := uint(0)
		if activeSeason != nil {
			seasonId = activeSeason.ID
		}

		switch r.Method {
		case "GET":
			{
				viewMode := r.URL.Query().Get("viewMode")
				if viewMode == "list-button" {
					fineList := finesListButton("Open Fine List", "list", false, false)
					fineList.Render(ctx, w)
					return
				} else if viewMode == "sheet-button" {
					fineList := finesListButton("Open Court Sheet", "sheet", false, false)
					fineList.Render(ctx, w)
					return
				} else if viewMode == "all-button" {
					allFineList := fineListButton(true)
					allFineList.Render(ctx, w)
					return
				}

				if seasonId == 0 || teamId == 0 {
					warning("No active season/team - fine list is unavailable").Render(ctx, w)
					return
				}

				var pageId = 0
				pageStr := r.URL.Query().Get("page")
				if len(pageStr) == 0 {
					pageId = 0
				} else {
					pageIdUint, err := strconv.ParseInt(pageStr, 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
						return
					}
					pageId = int(pageIdUint)
				}

				limitStr := r.URL.Query().Get("limit")
				var limit = 999999
				if len(limitStr) == 0 {
					limit = 999999
				} else {
					limit, err := strconv.ParseInt(limitStr, 10, 64)
					if err != nil || limit == 0 {
						http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
						return
					}
				}

				finemasterPage := false
				splitUrl := strings.Split(r.Header.Get("Referer"), "/")
				for _, urlBit := range splitUrl {
					if urlBit == "finemaster" {
						finemasterPage = true
					}
				}

				fineWithPlayers, err := GetFineWithPlayers(db, seasonId, teamId, 0, limit)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
					return
				}

				standAlone := r.URL.Query().Get("standAlone") == "true"

				full := r.URL.Query().Get("full") == "true"

				switch viewMode {
				case "sheet":
					log.Printf("fineHandler - fineListSheet %+v", fineWithPlayers)

					team, err := GetTeam(db, teamId)
					if err != nil || team == nil {
						warning(fmt.Sprintf("No active team")).Render(ctx, w)
						return
					}

					fineListSheet := fineListSheet(team, fineWithPlayers, standAlone, full)
					fineListSheet.Render(ctx, w)
					return
				default:
					log.Printf("fineHandler - fineList ")

					fineList := fineList(fineWithPlayers, pageId, 0, finemasterPage, false)
					fineList.Render(ctx, w)
				}

				return
			}
		case "POST":
			{

				if err := r.ParseForm(); err != nil {
					http.Error(w, "fineHandler - Invalid form data", http.StatusBadRequest)
					return
				}
				if seasonId == 0 || teamId == 0 {
					warning("No active season/team - cannot save fine").Render(ctx, w)
					http.Error(w, "No active season/team", http.StatusBadRequest)
					return
				}
				createdFines := []Fine{}
				createdPFines := []PresetFine{}

				playerIdStr := r.FormValue("playerId")
				context := r.FormValue("context")

				reason := r.FormValue("reason")

				displayOrderStr := r.FormValue("displayOrder")

				displayOrder, err := strconv.ParseInt(displayOrderStr, 10, 64)
				if err != nil || displayOrder == 0 {
					http.Error(w, fmt.Sprintf("Error parsing displayOrder %v", err), http.StatusBadRequest)
					return
				}
				isKudos := r.FormValue("isKudos") == "true"

				icon := r.FormValue("icon")

				if len(playerIdStr) == 0 || r.FormValue("fineOption") == "applyAgain" {

					amountStr := r.FormValue("amount")
					if len(amountStr) == 0 || amountStr == "0" {
						amountStr = "2"
					}
					amount, err := strconv.ParseFloat(amountStr, 64) // 64 specifies the bit size of the float type
					if err != nil {
						// Handle the error if the conversion fails
						http.Error(w, "Invalid amount", http.StatusBadRequest)
						return
					}

					log.Printf("%+v %+v", r.FormValue("fineOption"), r.FormValue("fineOption") == "applyAgain")

					if r.FormValue("fineOption") == "applyAgain" {
						var suggestedPFine = &PresetFine{
							Reason:       reason,
							Approved:     false,
							Amount:       amount,
							Context:      context,
							DisplayOrder: displayOrder,
							IsKudos:      isKudos,
							Icon:         icon,
						}
						err := SavePresetFine(db, suggestedPFine)
						if err != nil {
							// Handle the error if the conversion fails
							http.Error(w, "SavePresetFine failed", http.StatusBadRequest)
							return
						} else {
							createdPFines = append(createdPFines, *suggestedPFine)
						}
					}
				}
				if len(playerIdStr) > 0 {

					playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
					if err != nil {
						log.Printf(" fineHanlder - POST -  Error get playerIdStr fine: %v", err)
						http.Error(w, "Internal playerIdStr Error", http.StatusInternalServerError)
						return
					}

					approved := r.FormValue("approved") == "on"

					presetFineIds := r.Form["presetFineId"]

					for _, pfIdStr := range presetFineIds {
						if len(pfIdStr) == 0 {
							http.Error(w, "Invalid form data", http.StatusBadRequest)
							return
						}
						amountStr := r.FormValue("amount")
						if len(amountStr) == 0 || amountStr == "0" {
							amountStr = "2"
						}
						amount, err := strconv.ParseFloat(amountStr, 64) // 64 specifies the bit size of the float type
						if err != nil {
							// Handle the error if the conversion fails
							errComp := errMsg("Invalid amount")
							errComp.Render(GetContext(r, db), w)
						}
						if pfIdStr == "-1" && len(presetFineIds) > 1 {
							errComp := errMsg("Cannot select \"Fine is not listed here\" with others")
							errComp.Render(GetContext(r, db), w)
						}

						var fineAt = time.Now()
						activeMatch, actMatchERr := GetActiveMatch(db)
						if actMatchERr != nil {
							log.Printf("fineHandler - Error fetching active match: %v", actMatchERr)
							errComp := errMsg("Cannot select get the active match?")
							errComp.Render(GetContext(r, db), w)
						} else if activeMatch != nil {
							fineAt = *activeMatch.StartTime
						}

						if pfIdStr == "-1" {

							fine := Fine{
								Amount:   amount,
								Reason:   reason,
								Context:  context,
								PlayerID: uint(playerId),
								FineAt:   fineAt,
								Approved: approved,
								SeasonID: seasonId,
								TeamID:   teamId,
							}

							if err := SaveFine(db, &fine); err != nil {
								log.Printf("Error saving fine: %v", err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							} else {
								createdFines = append(createdFines, fine)
							}
						} else {

							// Parse the string ID to an unsigned integer
							pfId, err := strconv.ParseUint(pfIdStr, 10, 64)
							if err != nil {
								http.Error(w, "fine handler - POST Invalid strconv.ParseUint(pfIdStr", http.StatusBadRequest)
								return
							}

							// Assuming a function GetPresetFineByID that returns a *PresetFine struct for a given ID
							presetFine, err := GetPresetFine(db, pfId, "")
							if err != nil {
								http.Error(w, "Invalid GetPresetFine data", http.StatusBadRequest)
								return
							}

							if len(context) == 0 {
								context = presetFine.Context
							}

							fine := Fine{
								Amount:   presetFine.Amount,
								Reason:   presetFine.Reason,
								FineAt:   time.Now(),
								Context:  context,
								PlayerID: uint(playerId),
								Approved: approved,
								SeasonID: seasonId,
								TeamID:   teamId,
							}

							if err := SaveFine(db, &fine); err != nil {
								log.Printf("Error saving fine: %v", err)
								http.Error(w, "Internal Server Error", http.StatusInternalServerError)
								return
							} else {
								createdFines = append(createdFines, fine)
							}
						}

					}
				}

				/*	dontRedirect := r.FormValue("dontRedirect")
					if dontRedirect != "true" {
						w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
					}*/

				success := fineAddRes(createdFines, createdPFines)
				success.Render(GetContext(r, db), w)
				return
			}
		case "DELETE":

			fIDStr := r.URL.Query().Get("fid")

			fID, err := strconv.ParseUint(fIDStr, 10, 64)
			if err != nil || fID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine id %v", err), http.StatusBadRequest)
				return
			}

			if err := DeleteFineByID(db, uint(fID)); err != nil {
				warning := warning(fmt.Sprintf("Error deleting preset fine: %v", err))
				warning.Render(GetContext(r, db), w)
				return
			}
			//w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
			success := success("fine deleted")
			success.Render(GetContext(r, db), w)
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

type PlayerWithDrinks struct {
	TotalCurrent   float64
	TotalPaid      float64
	TotalDrinks    float64
	Player         Player
	DrinkPurchases []DrinkPurchase
	DrinkPayments  []DrinkPayment
}

func drinkHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{

				players, err := GetPlayers(db, 0, 999999)
				if err != nil {
					http.Error(w, "Error fetching players", http.StatusInternalServerError)
					return
				}

				drinkPayments, err := GetDrinkPayments(db)
				if err != nil {
					err := errMsg(fmt.Sprintf("Error fetching drink payments: %v", err))
					err.Render(GetContext(r, db), w)
					return
				}

				drinkTypes, err := GetDrinkTypes(db)
				if err != nil {
					err := errMsg(fmt.Sprintf("Error fetching drink types: %v", err))
					err.Render(GetContext(r, db), w)
					return
				}

				drinkPurchases, err := GetDrinkPurchases(db)
				if err != nil {
					http.Error(w, "Error fetching drink purchases", http.StatusInternalServerError)
					return
				}

				playersWithDrinks := []PlayerWithDrinks{}
				for _, player := range players {
					playerWithDrinks := PlayerWithDrinks{
						Player:       player,
						TotalCurrent: 0,
						TotalPaid:    0,
						TotalDrinks:  0,
					}
					for _, drinkPurchase := range drinkPurchases {
						if drinkPurchase.PlayerId == player.ID {
							playerWithDrinks.TotalCurrent += drinkPurchase.Amount
							playerWithDrinks.TotalDrinks += 1
						}
					}
					for _, drinkPayment := range drinkPayments {
						if drinkPayment.PlayerId == player.ID {
							playerWithDrinks.TotalPaid += drinkPayment.Amount
							playerWithDrinks.TotalCurrent -= drinkPayment.Amount
						}
					}
					playersWithDrinks = append(playersWithDrinks, playerWithDrinks)
				}

				drinkPurchasesComp := DrinkPurchases(playersWithDrinks, drinkTypes)
				drinkPurchasesComp.Render(GetContext(r, db), w)
				return

			}
		case "POST":
			{
				drinkType := DrinkType{
					Name: "Water",
				}

				if err := AddDrinkType(db, &drinkType); err != nil {
					http.Error(w, "Error adding drink type", http.StatusInternalServerError)
					return
				}

				success := success("Drink type added")
				success.Render(GetContext(r, db), w)
			}
		}
	}
}

func drinkTypeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				drinkTypes, err := GetDrinkTypes(db)
				if err != nil {
					http.Error(w, "Error fetching drink types", http.StatusInternalServerError)
					return
				}

				drinkTypesComp := DrinkTypes(drinkTypes)
				drinkTypesComp.Render(GetContext(r, db), w)
			}
		case "POST":
			{
				drinkType := DrinkType{
					Name: "Water",
				}

				if err := AddDrinkType(db, &drinkType); err != nil {
					http.Error(w, "Error adding drink type", http.StatusInternalServerError)
					return
				}

				success := success("Drink type added")
				success.Render(GetContext(r, db), w)
			}
		}
	}
}

func adminHandler(db *gorm.DB) http.HandlerFunc {
	DEFAULT_PASS := "pass"

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{

				if err := r.ParseForm(); err != nil {
					http.Error(w, "fineHandler - Invalid form data", http.StatusBadRequest)
					return
				}
				if googleAuthEnabled() {
					handleGoogleAdminAuth(db, w, r)
					return
				}
				password := r.FormValue("password")
				teamIdStr := r.FormValue("teamId")
				pass := os.Getenv("PASS")
				if pass == "" {
					pass = DEFAULT_PASS
				}

				if password != pass {
					success := warning(fmt.Sprintf("Sorry Invalid password (%v) Hint: ends with a number", password))
					success.Render(GetContext(r, db), w)
					http.Error(w, "Invalid password", http.StatusBadRequest)
					return
				}

				if len(teamIdStr) == 0 {
					warn := warning("Please select a team")
					warn.Render(GetContext(r, db), w)
					http.Error(w, "Missing teamId", http.StatusBadRequest)
					return
				}

				teamIdInt, err := strconv.ParseUint(teamIdStr, 10, 64)
				if err != nil {
					warn := warning(fmt.Sprintf("Invalid teamId: %v", err))
					warn.Render(GetContext(r, db), w)
					http.Error(w, "Invalid teamId", http.StatusBadRequest)
					return
				}

				team, err := GetTeam(db, uint(teamIdInt))
				if err != nil || team == nil {
					warn := warning(fmt.Sprintf("Team not found (teamId=%s)", teamIdStr))
					warn.Render(GetContext(r, db), w)
					http.Error(w, "Team not found", http.StatusBadRequest)
					return
				}

				// Persist admin login for this team
				http.SetCookie(w, &http.Cookie{
					Name:     adminTokenCookieName,
					Value:    makeAdminToken(team.ID, adminTokenSecret()),
					Path:     "/",
					HttpOnly: true,
					Secure:   false, // Set to true in production with HTTPS
					SameSite: http.SameSiteStrictMode,
					MaxAge:   60 * 60 * 24 * 30, // 30 days
				})

				var url = "/finemaster"
				log.Printf("REdirecting to %s", url)
				w.Header().Set("HX-Redirect", url)
				w.Header().Set("HX-Reload", "true")
			}
		case "DELETE":

			fIDStr := r.URL.Query().Get("fid")

			fID, err := strconv.ParseUint(fIDStr, 10, 64)
			if err != nil || fID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine id %v", err), http.StatusBadRequest)
				return
			}

			if err := DeleteFineByID(db, uint(fID)); err != nil {
				log.Printf("Error deleting preset fine: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
			w.WriteHeader(http.StatusOK)
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func fineAddHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		activeSeason, _ := GetActiveSeason(db)
		seasonId := uint(0)
		if activeSeason != nil {
			seasonId = activeSeason.ID
		}

		var newFines []uint
		if r.Method == "POST" {

			var input struct {
				Reason string `json:"reason"`
			}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var suggestedPFine = &PresetFine{
				Reason:   input.Reason,
				Approved: true,
				Amount:   2,
			}
			err := SavePresetFine(db, suggestedPFine)
			if err != nil {
				http.Error(w, "SavePresetFine failed", http.StatusBadRequest)
				return
			}
			newFines = append(newFines, suggestedPFine.ID)

			//success := success(fmt.Sprintf("Added Fine - %s", input.Reason))
			//success.Render(GetContext(r, db), w)
		}

		playersWithFines, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
		if err != nil {
			log.Printf("Error fetching players with fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		pFines, err := GetPresetFines(db, false)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		/*	var approvedPFines = []PresetFine{}
			var pendingPFines = []PresetFine{}

			for _, f := range pFines {
				if f.Approved {
					approvedPFines = append(approvedPFines, f)
				}else{
					pendingPFines = append(pendingPFines, f)
				}
			}*/
		fsComp := fineSuperSelect(playersWithFines, pFines, newFines, "2")
		fsComp.Render(ctx, w)
	}
}

func fineMultiHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "fineMultiHandler - Invalid form data", http.StatusBadRequest)
			return
		}

		pfineIDs := r.Form["pfines[]"]
		playerIDs := r.Form["players[]"]
		context := r.FormValue("context")
		savedFines := []Fine{}

		warnStr := "No players selected"

		log.Printf("pfineIDs: %+v %d", len(pfineIDs), len(playerIDs))

		for _, pfineIDStr := range pfineIDs {

			for _, playerIDStr := range playerIDs {
				playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
				if err != nil {
					errComp := errMsg(F("Invalid player ID: [%s]", playerIDStr))
					errComp.Render(GetContext(r, db), w)
					return
				} else {
					warnStr = ""
				}
				fine, err := GetFineFromPreset(db, pfineIDStr)
				if err != nil {
					errComp := errMsg(F("Could not GetFineFromPreset: [%s]", pfineIDStr))
					errComp.Render(GetContext(r, db), w)
					return
				}
				fine.PlayerID = uint(playerID)
				fine.Context = context

				activeMatch, err := GetActiveMatch(db)
				if err != nil {
					log.Printf("fineMultiHandler - Error fetching active match: %v", err)
				}

				if activeMatch != nil {
					fine.MatchId = activeMatch.ID
					if activeMatch.StartTime != nil {
						fine.FineAt = *activeMatch.StartTime
					} else {
						fine.FineAt = time.Now()
					}
				} else {
					warnStr = "No active match (using current time as fine time)"
					fine.FineAt = time.Now()
				}

				ctx := GetContext(r, db)
				activeSeason, err := GetActiveSeason(db)
				if err != nil {
					errComp := errMsg(F("Could not get active season %v", err))
					errComp.Render(ctx, w)
				} else if activeSeason != nil {
					fine.SeasonID = uint(activeSeason.ID)
				} else {
					warnStr = warnStr + "\nNo active season"
				}

				fine.TeamID = getTeamId(ctx)
				if fine.TeamID == 0 || fine.SeasonID == 0 {
					errMsg("No active season/team - cannot save fine").Render(ctx, w)
					http.Error(w, "No active season/team", http.StatusBadRequest)
					return
				}

				fine.Approved = config.DefaultToApproved

				fine.CreatedAt = time.Now()
				err = SaveFine(db, fine)
				if err != nil {
					http.Error(w, "Invalid player ID", http.StatusBadRequest)
					errComp := errMsg(F("Could not Save Fine %+v", fine))
					errComp.Render(ctx, w)
					return
				} else {
					savedFines = append(savedFines, *fine)
				}
			}
		}

		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		activeSeason2, _ := GetActiveSeason(db)
		seasonId := uint(0)
		if activeSeason2 != nil {
			seasonId = activeSeason2.ID
		}
		playersWithFines, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
		if err != nil {
			log.Printf("Error fetching players with fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		pFines, err := GetPresetFines(db, true)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res := fineSuperSelectResults(playersWithFines, pFines, savedFines, warnStr)
		res.Render(ctx, w)
	}
}

func teamCourtHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				if err := r.ParseForm(); err != nil {
					warning := warning("specificTeamHandler - Invalid form data")
					warning.Render(GetContext(r, db), w)
					return
				}

				teamId := r.FormValue("ID")
				if len(teamId) == 0 {
					warning := warning("No team ID provided")
					warning.Render(GetContext(r, db), w)
					return
				}

				teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
				if err != nil {
					warning := warning(fmt.Sprintf("Error parsing team ID: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				courtNotes := r.FormValue("courtNotes")
				team, err := GetTeam(db, uint(teamIdInt))
				if err != nil {
					warning := warning(fmt.Sprintf("teamCourtHandlerPOST - Error fetching team: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				team.CourtNotes = courtNotes

				newTeam, err := SaveTeam(db, team)
				if err != nil {
					warning := warning(fmt.Sprintf("Error saving team %+v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				teamCourtNotes := teamCourtNotes(newTeam, "Updated court notes")
				teamCourtNotes.Render(GetContext(r, db), w)
				return

			}
		case "GET":
			{
				mst, err := GetMatchSeasonTeam(db)
				if err != nil {
					warning := warning(fmt.Sprintf("teamCourtHandler GET - Error fetching GetMatchSeasonTeam: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				viewMode := r.URL.Query().Get("viewMode")
				switch viewMode {
				case "court-notes-form":
					teamForm := teamCourtNotesForm(*mst.Team)
					teamForm.Render(GetContext(r, db), w)
					return
				case "court-notes":
					teamCourtNotes := teamCourtNotes(*mst.Team, "")
					teamCourtNotes.Render(GetContext(r, db), w)
					return
				case "court-notes-form-button":
					noteButton := teamCourtNotesButton(true)
					noteButton.Render(GetContext(r, db), w)
					return

				default:
					warning := warning("No viewMode provided (court-notes-form or court-notes needed)")
					warning.Render(GetContext(r, db), w)

				}
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func teamActiveMatchHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{
				if err := r.ParseForm(); err != nil {
					warning := warning("specificTeamHandler - Invalid form data")
					warning.Render(GetContext(r, db), w)
					return
				}

				teamId := r.FormValue("ID")
				if len(teamId) == 0 {
					warning := warning("No team ID provided")
					warning.Render(GetContext(r, db), w)
					return
				}

				teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
				if err != nil {
					warning := warning(fmt.Sprintf("Error parsing team ID: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				activeMatchOverrideId := r.FormValue("activeMatchOverrideId")
				team, err := GetTeam(db, uint(teamIdInt))
				if err != nil {
					warning := warning(fmt.Sprintf("activeMatch handler - teamCourtHandlerPOST - Error fetching team: for teamId input:[%s] %v", teamId, err))
					warning.Render(GetContext(r, db), w)
					return
				}

				activeMatchOverrideIdInt, err := strconv.ParseUint(activeMatchOverrideId, 10, 64)
				if err != nil {
					warning := warning(fmt.Sprintf("Error parsing activeMatchOverrideId: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				team.ActiveMatchIDOverride = uint(activeMatchOverrideIdInt)

				newTeam, err := SaveTeam(db, team)
				if err != nil {
					warning := warning(fmt.Sprintf("Error saving team %+v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
				if err != nil {
					warning := warning(fmt.Sprintf("teamActiveMatchHandler - Error fetching matches: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				success := teamEditForm(newTeam, matches, success(fmt.Sprintf("Active Match ID set to %s", activeMatchOverrideId)))
				success.Render(GetContext(r, db), w)
				return

			}
		case "GET":
			{
				teamId := r.URL.Query().Get("teamId")
				if len(teamId) == 0 {

					teamsArr, err := GetTeams(db, 9999, 0)
					if len(teamsArr) == 0 {
						warning := warning("No teams exist")
						warning.Render(GetContext(r, db), w)
						addForm := teamListButton(true)

						addForm.Render(GetContext(r, db), w)
						return
					} else if err != nil {
						warning := warning(fmt.Sprintf("Error fetching teams: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					} else if len(teamsArr) == 1 {
						teamId = fmt.Sprintf("%d", teamsArr[0].ID)
					}
				}
				/*id, err := strconv.ParseUint(teamId, 10, 64)
				if err != nil {
					warning := warning(fmt.Sprintf("Error parsing team ID: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}*/

				matchSeasonTeam, err := GetMatchSeasonTeam(db)
				if err != nil {
					warning := warning(fmt.Sprintf("Error fetching GetMatchSeasonTeam: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}
				if matchSeasonTeam.Team != nil {
					if upcomingMatch, lineup, err := upcomingMatchLineupForTeam(db, matchSeasonTeam.Team.ID); err == nil {
						matchSeasonTeam.Match = upcomingMatch
						if matchSeasonTeam.Match != nil && lineup != nil {
							matchSeasonTeam.Match.LineupID = lineup.ID
						}
					}
				}

				viewMode := r.URL.Query().Get("viewMode")
				switch viewMode {
				case "button":
					teamList := teamActiveMatchButton(true, matchSeasonTeam)
					teamList.Render(GetContext(r, db), w)
					return
				case "add-or-override":

					/*team, err := GetTeam(db, uint(id))
					if err != nil {
						warning := warning(fmt.Sprintf("Error fetching team: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}
					if team == nil {
						warning := warning(fmt.Sprintf("Error fetching team: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}


					if team.ActiveMatchIDOverride > 0 {
						aMatch, err := GetMatch(db, team.ActiveMatchIDOverride)
						if err != nil {
							warning := warning(fmt.Sprintf("Error fetching active match: %v", err))
							warning.Render(GetContext(r, db), w)
							return
						}

						if aMatch == nil {
							warning := warning(fmt.Sprintf("No active match found (via ActiveMatchIDOverride)"))
							warning.Render(GetContext(r, db), w)
							return
						} else {
							activeMatch = aMatch
						}

					}*/

					matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
					if err != nil {
						warning := warning(fmt.Sprintf("Error fetching matches: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					matchSeasonTeam, err := GetMatchSeasonTeam(db)
					if err != nil {
						warning := warning(fmt.Sprintf("Error fetching GetMatchSeasonTeam: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}
					if matchSeasonTeam.Team != nil {
						if upcomingMatch, lineup, err := upcomingMatchLineupForTeam(db, matchSeasonTeam.Team.ID); err == nil {
							matchSeasonTeam.Match = upcomingMatch
							if matchSeasonTeam.Match != nil && lineup != nil {
								matchSeasonTeam.Match.LineupID = lineup.ID
							}
						}
					}

					teamActiveMatchForm := teamActiveMatchAddOrOverrideForm(matchSeasonTeam.Team, matchSeasonTeam, matches, "")
					teamActiveMatchForm.Render(GetContext(r, db), w)
					return

				default:
					warning := warning("No viewMode provided (button needed)")
					warning.Render(GetContext(r, db), w)
					return
				}

			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func teamActivateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		teamId := r.URL.Query().Get("teamId")
		if len(teamId) == 0 {
			warning("No team ID provided").Render(GetContext(r, db), w)
			return
		}

		teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
		if err != nil {
			warning(fmt.Sprintf("Error parsing team ID: %v", err)).Render(GetContext(r, db), w)
			return
		}

		team, err := GetTeam(db, uint(teamIdInt))
		if err != nil || team == nil {
			warning(fmt.Sprintf("Team not found (teamId=%s)", teamId)).Render(GetContext(r, db), w)
			return
		}
		if user, ok := requireGoogleAdminForTeam(db, r, team.ID); !ok {
			if user == nil {
				http.Error(w, "Admin auth required", http.StatusUnauthorized)
			} else {
				http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     adminTokenCookieName,
			Value:    makeAdminToken(team.ID, adminTokenSecret()),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   60 * 60 * 24 * 30,
		})
		if user, _, ok := currentAdminUser(r, db); ok {
			if err := activateAdminTeam(w, *user, *team); err != nil {
				warning(fmt.Sprintf("Error saving active team: %v", err)).Render(GetContext(r, db), w)
				return
			}
		} else {
			session, _ := store.Get(r, "session-name")
			setTeamSessionValues(*team, session)
			if err := session.Save(r, w); err != nil {
				warning(fmt.Sprintf("Error saving active team: %v", err)).Render(GetContext(r, db), w)
				return
			}
		}

		teams, err := GetTeams(db, 9999, 0)
		if err != nil {
			http.Error(w, "Error fetching teams", http.StatusInternalServerError)
			return
		}

		matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
		if err != nil {
			warning(fmt.Sprintf("teamActivateHandler - Error fetching matches: %v", err)).Render(GetContext(r, db), w)
			return
		}

		teamList(teams, matches, team.ID).Render(GetContext(r, db), w)
	}
}

func activateAdminTeam(w http.ResponseWriter, user AdminUser, team Team) error {
	setAdminSessionCookie(w, user.ID, team.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     adminTokenCookieName,
		Value:    makeAdminToken(team.ID, adminTokenSecret()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 30,
	})
	return nil
}

func adminTeamAccessRequestHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		user, _, ok := currentAdminUser(r, db)
		if !ok {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		teamID64, err := strconv.ParseUint(r.FormValue("teamId"), 10, 64)
		if err != nil || teamID64 == 0 {
			http.Error(w, "Invalid team", http.StatusBadRequest)
			return
		}
		team, err := GetTeam(db, uint(teamID64))
		if err != nil {
			http.Error(w, "Team not found", http.StatusBadRequest)
			return
		}
		if canAdminAccessTeam(db, user.ID, team.ID) {
			if err := activateAdminTeam(w, *user, *team); err != nil {
				http.Error(w, "Could not activate team", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Redirect", "/")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err := RequestAdminAccess(db, user.ID, team.ID, adminRoleTeamAdmin); err != nil {
			http.Error(w, "Could not request admin access", http.StatusInternalServerError)
			return
		}
		success("Access request sent. A team admin or super admin can approve it.").Render(GetContext(r, db), w)
	}
}

func adminTeamSwitchHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		user, _, ok := currentAdminUser(r, db)
		if !ok {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		teamID64, err := strconv.ParseUint(r.FormValue("teamId"), 10, 64)
		if err != nil || teamID64 == 0 {
			http.Error(w, "Invalid team", http.StatusBadRequest)
			return
		}
		team, err := GetTeam(db, uint(teamID64))
		if err != nil || team == nil {
			http.Error(w, "Team not found", http.StatusBadRequest)
			return
		}
		if !canAdminAccessTeam(db, user.ID, team.ID) {
			http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
			return
		}
		if err := activateAdminTeam(w, *user, *team); err != nil {
			http.Error(w, "Could not activate team", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func teamHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var adminUser *AdminUser
		if googleAuthEnabled() {
			user, _, ok := currentAdminUser(r, db)
			if !ok {
				w.Header().Set("HX-Redirect", "/finemaster/auth")
				http.Redirect(w, r, "/finemaster/auth", http.StatusSeeOther)
				return
			}
			adminUser = user
		}
		switch r.Method {
		case "GET":
			{
				ctx := GetContext(r, db)
				activeTeamID := getTeamId(ctx)
				viewMode := r.URL.Query().Get("viewMode")
				switch viewMode {
				case "button":
					teamList := teamListButton(true)
					teamList.Render(ctx, w)
					return

				case "add":
					team := teamAddForm()
					team.Render(GetContext(r, db), w)
					return
				case "edit-button":
					teamId := r.URL.Query().Get("teamId")
					if len(teamId) == 0 {
						warning := warning("No team ID provided")
						warning.Render(GetContext(r, db), w)
						return
					}

					teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
					if err != nil {
						warning := warning(fmt.Sprintf("Error parsing team ID: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					team, err := GetTeam(db, uint(teamIdInt))
					if err != nil {
						warning := warning(fmt.Sprintf("teamHandler GET - Error fetching team: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}
					if adminUser != nil && !canAdminAccessTeam(db, adminUser.ID, team.ID) {
						http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
						return
					}

					teamEditFormButton := teamEditFormButton(*team, false)
					teamEditFormButton.Render(GetContext(r, db), w)
					return
				case "edit":
					teamId := r.URL.Query().Get("teamId")
					if len(teamId) == 0 {
						warning := warning("No team ID provided")
						warning.Render(GetContext(r, db), w)
						return
					}

					teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
					if err != nil {
						warning := warning(fmt.Sprintf("Error parsing team ID: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					team, err := GetTeam(db, uint(teamIdInt))
					if err != nil {
						warning := warning(fmt.Sprintf("teamHandler GET - Error fetching team: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}
					if adminUser != nil && !canAdminAccessTeam(db, adminUser.ID, team.ID) {
						http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
						return
					}

					matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
					if err != nil {
						warning := warning(fmt.Sprintf("teamHandler - Error fetching matches: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					editTeam := teamEditForm(*team, matches, nil)
					editTeam.Render(GetContext(r, db), w)
					return
				case "list":
				default:
					teams, err := GetTeams(db, 9999, 0)
					if err != nil {
						http.Error(w, "Error fetching teams", http.StatusInternalServerError)
						return
					}
					if adminUser != nil {
						teams = FilterTeamsForAdminUser(db, adminUser.ID, teams)
					}

					matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
					if err != nil {
						warning := warning(fmt.Sprintf("teamHandler - Error fetching matches: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					includeFrame := r.URL.Query().Get("includeFrame") == "true"

					if includeFrame {
						teamList := teamListFrame(teams, matches, activeTeamID)
						teamList.Render(ctx, w)
						return
					}
					teamList := teamList(teams, matches, activeTeamID)
					teamList.Render(ctx, w)
					return

				}
			}
		case "PUT":
			{
				if err := r.ParseForm(); err != nil {
					warning := warning("teamHandler - Invalid form data")
					warning.Render(GetContext(r, db), w)
					return
				}

				if ID := r.FormValue("ID"); len(ID) > 0 {
					id, err := strconv.ParseUint(ID, 10, 64)
					if err != nil {
						waring := warning(fmt.Sprintf("Error parsing team ID: %v", err))
						waring.Render(GetContext(r, db), w)
						return
					}

					teamName := r.FormValue("teamName")
					teamKey := r.FormValue("teamKey")
					teamMemberPass := r.FormValue("teamMemberPass")
					teamAdminPass := r.FormValue("teamAdminPass")
					activeMatchOverrideID, err := parseOptionalUint(r.FormValue("activeMatchOverrideId"))
					if err != nil {
						waring := warning(fmt.Sprintf("Error parsing active match override ID: %v", err))
						waring.Render(GetContext(r, db), w)
						return
					}
					sportyCompetitionID, err := parseOptionalUint(r.FormValue("sportyCompetitionId"))
					if err != nil {
						warning(fmt.Sprintf("Error parsing Sporty competition ID: %v", err)).Render(GetContext(r, db), w)
						return
					}
					sportyGradeID, err := parseOptionalUint(r.FormValue("sportyGradeId"))
					if err != nil {
						warning(fmt.Sprintf("Error parsing Sporty grade ID: %v", err)).Render(GetContext(r, db), w)
						return
					}
					sportyTeamID, err := parseOptionalUint(r.FormValue("sportyTeamId"))
					if err != nil {
						warning(fmt.Sprintf("Error parsing Sporty team ID: %v", err)).Render(GetContext(r, db), w)
						return
					}
					showFineAddOnHomePage := r.FormValue("showFineAddOnHomePage") == "on"
					showPitchMatchOnHomePage := r.FormValue("showPitchMatchOnHomePage") == "on"
					showCourtSheetOnHomePage := r.FormValue("showCourtSheetOnHomePage") == "on"
					enablePublicFeedbackForm := r.FormValue("enablePublicFeedbackForm") == "on"
					showCourtTotals := r.FormValue("showCourtTotals") == "on"
					enableFinesModule := r.FormValue("enableFinesModule") == "on"
					enableLineupsModule := r.FormValue("enableLineupsModule") == "on"
					enableMatchesModule := r.FormValue("enableMatchesModule") == "on"
					enablePlayersModule := r.FormValue("enablePlayersModule") == "on"
					enablePaymentsModule := r.FormValue("enablePaymentsModule") == "on"
					enableCourtModule := r.FormValue("enableCourtModule") == "on"
					enableLeaderboardModule := r.FormValue("enableLeaderboardModule") == "on"
					allowAdminRegistration := r.FormValue("allowAdminRegistration") == "on"
					sportyAutoMatchSync := r.FormValue("sportyAutoMatchSync") == "on"
					lineupPlayerCount, err := parseLineupPlayerCount(r.FormValue("lineupPlayerCount"))
					if err != nil {
						warning(err.Error()).Render(GetContext(r, db), w)
						return
					}

					team, err := GetTeam(db, uint(id))
					if err != nil {
						warning(fmt.Sprintf("Error fetching team %+v", err)).Render(GetContext(r, db), w)
						return
					}
					if adminUser != nil && !canAdminAccessTeam(db, adminUser.ID, team.ID) {
						http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
						return
					}
					team.TeamName = teamName
					team.TeamKey = teamKey
					team.TeamAdminPass = teamAdminPass
					team.TeamMemberPass = teamMemberPass
					team.ActiveMatchIDOverride = activeMatchOverrideID
					team.ShowFineAddOnHomePage = showFineAddOnHomePage
					team.ShowPitchMatchOnHomePage = showPitchMatchOnHomePage
					team.ShowCourtSheetOnHomePage = showCourtSheetOnHomePage
					team.EnablePublicFeedbackForm = enablePublicFeedbackForm
					team.ShowCourtTotals = showCourtTotals
					team.EnableFinesModule = enableFinesModule
					team.EnableLineupsModule = enableLineupsModule
					team.EnableMatchesModule = enableMatchesModule
					team.EnablePlayersModule = enablePlayersModule
					team.EnablePaymentsModule = enablePaymentsModule
					team.EnableCourtModule = enableCourtModule
					team.EnableLeaderboardModule = enableLeaderboardModule
					team.AllowAdminRegistration = allowAdminRegistration
					team.LineupPlayerCount = lineupPlayerCount
					team.SportyCompetitionID = sportyCompetitionID
					team.SportyOrgID = 0
					team.SportyGradeID = sportyGradeID
					team.SportyTeamID = sportyTeamID
					team.SportyAutoMatchSync = sportyAutoMatchSync

					newTeam, err := SaveTeam(db, team)
					if err != nil {
						warning := warning(fmt.Sprintf("Error saving team %+v", err))
						warning.Render(GetContext(r, db), w)
						return
					}
					var syncResult *SportySyncResult
					var syncErr error
					if newTeam.SportyAutoMatchSync {
						syncResult, syncErr = SyncSportyMatches(db, &newTeam)
					}
					session, _ := store.Get(r, "session-name")
					saveTeamToSession(r, newTeam, session)

					matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
					if err != nil {
						warning := warning(fmt.Sprintf("teamHandler - Error fetching matches: %v", err))
						warning.Render(GetContext(r, db), w)
						return
					}

					if syncErr != nil {
						teamEditForm(newTeam, matches, warning(fmt.Sprintf("Team updated, but Sporty sync failed: %v", syncErr))).Render(GetContext(r, db), w)
					} else if syncResult != nil {
						teamEditForm(newTeam, matches, sportySyncResultMessage("Team updated.", syncResult, true)).Render(GetContext(r, db), w)
					} else {
						teamEditForm(newTeam, matches, successWithHomepageLink("Team updated")).Render(GetContext(r, db), w)
					}
					return
				}

			}
		case "POST":
			{
				if err := r.ParseForm(); err != nil {
					warning := warning("teamHandler - Invalid form data")
					warning.Render(GetContext(r, db), w)
					return
				}

				teamName := r.FormValue("teamName")
				teamKey := r.FormValue("teamKey")
				teamMemberPass := r.FormValue("teamMemberPass")
				teamAdminPass := r.FormValue("teamAdminPass")
				sportyCompetitionID, err := parseOptionalUint(r.FormValue("sportyCompetitionId"))
				if err != nil {
					warning(fmt.Sprintf("Error parsing Sporty competition ID: %v", err)).Render(GetContext(r, db), w)
					return
				}
				sportyGradeID, err := parseOptionalUint(r.FormValue("sportyGradeId"))
				if err != nil {
					warning(fmt.Sprintf("Error parsing Sporty grade ID: %v", err)).Render(GetContext(r, db), w)
					return
				}
				sportyTeamID, err := parseOptionalUint(r.FormValue("sportyTeamId"))
				if err != nil {
					warning(fmt.Sprintf("Error parsing Sporty team ID: %v", err)).Render(GetContext(r, db), w)
					return
				}
				showFineAddOnHomePage := r.FormValue("showFineAddOnHomePage") == "on"
				showPitchMatchOnHomePage := r.FormValue("showPitchMatchOnHomePage") == "on"
				showCourtSheetOnHomePage := r.FormValue("showCourtSheetOnHomePage") == "on"
				enablePublicFeedbackForm := r.FormValue("enablePublicFeedbackForm") == "on"
				showCourtTotals := r.FormValue("showCourtTotals") == "on"
				enableFinesModule := r.FormValue("enableFinesModule") == "on"
				enableLineupsModule := r.FormValue("enableLineupsModule") == "on"
				enableMatchesModule := r.FormValue("enableMatchesModule") == "on"
				enablePlayersModule := r.FormValue("enablePlayersModule") == "on"
				enablePaymentsModule := r.FormValue("enablePaymentsModule") == "on"
				enableCourtModule := r.FormValue("enableCourtModule") == "on"
				enableLeaderboardModule := r.FormValue("enableLeaderboardModule") == "on"
				allowAdminRegistration := r.FormValue("allowAdminRegistration") == "on"
				sportyAutoMatchSync := r.FormValue("sportyAutoMatchSync") == "on"
				lineupPlayerCount, err := parseLineupPlayerCount(r.FormValue("lineupPlayerCount"))
				if err != nil {
					warning(err.Error()).Render(GetContext(r, db), w)
					return
				}

				team := Team{
					TeamName:                 teamName,
					TeamKey:                  teamKey,
					TeamAdminPass:            teamAdminPass,
					TeamMemberPass:           teamMemberPass,
					ShowFineAddOnHomePage:    showFineAddOnHomePage,
					ShowPitchMatchOnHomePage: showPitchMatchOnHomePage,
					ShowCourtSheetOnHomePage: showCourtSheetOnHomePage,
					EnablePublicFeedbackForm: enablePublicFeedbackForm,
					ShowCourtTotals:          showCourtTotals,
					EnableFinesModule:        enableFinesModule,
					EnableLineupsModule:      enableLineupsModule,
					EnableMatchesModule:      enableMatchesModule,
					EnablePlayersModule:      enablePlayersModule,
					EnablePaymentsModule:     enablePaymentsModule,
					EnableCourtModule:        enableCourtModule,
					EnableLeaderboardModule:  enableLeaderboardModule,
					AllowAdminRegistration:   allowAdminRegistration,
					LineupPlayerCount:        lineupPlayerCount,
					SportyCompetitionID:      sportyCompetitionID,
					SportyOrgID:              0,
					SportyGradeID:            sportyGradeID,
					SportyTeamID:             sportyTeamID,
					SportyAutoMatchSync:      sportyAutoMatchSync,
				}

				team, err = SaveTeam(db, &team)
				if err != nil {
					warning := warning(fmt.Sprintf("Error saving team %+v", err))
					warning.Render(GetContext(r, db), w)
					return
				}
				if adminUser != nil && !isSuperAdmin(db, adminUser.ID) {
					if err := GrantAdminRole(db, adminUser.ID, team.ID, adminRoleTeamAdmin); err != nil {
						warning(fmt.Sprintf("Team created, but admin role assignment failed: %v", err)).Render(GetContext(r, db), w)
						return
					}
				}
				var syncResult *SportySyncResult
				var syncErr error
				if team.SportyAutoMatchSync {
					syncResult, syncErr = SyncSportyMatches(db, &team)
				}
				session, _ := store.Get(r, "session-name")
				saveTeamToSession(r, team, session)

				matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
				if err != nil {
					warning := warning(fmt.Sprintf("teamHandler - Error fetching matches: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}

				if syncErr != nil {
					teamEditForm(team, matches, warning(fmt.Sprintf("Team created, but Sporty sync failed: %v", syncErr))).Render(GetContext(r, db), w)
				} else if syncResult != nil {
					teamEditForm(team, matches, sportySyncResultMessage("Team created.", syncResult, true)).Render(GetContext(r, db), w)
				} else {
					teamEditForm(team, matches, successWithHomepageLink("Team created")).Render(GetContext(r, db), w)
				}
				return
			}
		case "DELETE":
			{
				teamId := r.URL.Query().Get("teamId")
				if len(teamId) == 0 {
					warning := warning("No team ID provided")
					warning.Render(GetContext(r, db), w)
					return
				}
				teamIdInt, err := strconv.ParseUint(teamId, 10, 64)
				if err != nil {
					warning := warning(fmt.Sprintf("teanHandler Error parsing team ID: %v", err))
					warning.Render(GetContext(r, db), w)
					return
				}
				if adminUser != nil && !canAdminAccessTeam(db, adminUser.ID, uint(teamIdInt)) {
					http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
					return
				}

				delErr := DeleteTeam(db, uint(teamIdInt))
				if delErr != nil {
					warning := warning(fmt.Sprintf("Error deleting team: %v", delErr))
					warning.Render(GetContext(r, db), w)
					return
				}
				success := success("Team deleted")
				success.Render(GetContext(r, db), w)
				return
			}
		default:
			warning := warning("Method not allowed")
			warning.Render(GetContext(r, db), w)
			return
		}
	}

}

func fineApproveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			fIDStr := r.FormValue("fid")
			amountStr := r.FormValue("amount")
			approved := r.FormValue("approved") == "on"

			fID, err := strconv.ParseUint(fIDStr, 10, 64)
			if err != nil || fID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine ID: %v", err), http.StatusBadRequest)
				return
			}

			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing amount: %v", err), http.StatusBadRequest)
				return
			}
			if err := ApproveFine(db, uint(fID), amount, approved); err != nil {
				log.Printf("Error approving fine with specified amount: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if approved {
				success := success("Approved")
				success.Render(GetContext(r, db), w)
			} else {
				success := warning("Declined")
				success.Render(GetContext(r, db), w)
			}

			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func fineQuickHideHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			showOrHide := chi.URLParam(r, "showOrHide") == "show"

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			pfIDStr := r.FormValue("pfid")
			pass := r.FormValue("pass")
			// Convert 'fid' to uint64
			pfID, err := strconv.ParseUint(pfIDStr, 10, 64)
			if err != nil || pfID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine ID: %v", err), http.StatusBadRequest)
				return
			}
			if showOrHide {
				if err := QuickHideFine(db, uint(pfID), false); err != nil {
					log.Printf("Error approving fine with specified amount: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				if err := QuickHideFine(db, uint(pfID), true); err != nil {
					log.Printf("Error approving fine with specified amount: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			var url = fmt.Sprintf("/finemaster/%s?pf=true&upf=%d#pf-%d", pass, pfID, pfID)
			log.Printf("REdirecting to %s", url)
			w.Header().Set("HX-Redirect", url)
			w.Header().Set("HX-Reload", "true")
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func presetFineApproveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			{

				pfIDStr := r.URL.Query().Get("pfid")

				pfID, err := strconv.ParseUint(pfIDStr, 10, 64)
				if err != nil || pfID == 0 {
					http.Error(w, fmt.Sprintf("Error parsing fine id %v", err), http.StatusBadRequest)
					return
				}

				if err := ApprovePresetFine(db, uint(pfID)); err != nil {
					log.Printf("Error deleting preset fine: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				success := success("Approved (will appear in fine list)")
				success.Render(GetContext(r, db), w)
				return
			}
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GenerateUrl(baseURL string, queryParams *HomeQueryParams) (*string, error) {
	fullURL, err := GenerateURLWithParams(baseURL, queryParams)
	if err != nil {
		return nil, err
	}
	return &fullURL, nil
}

type HomeQueryParams struct {
	FinesOpen       bool `schema:"f"`
	FineListOpen    bool `schema:"fl"`
	PlayerOpen      bool `schema:"p"`
	PresetFinesOpen bool `schema:"pf"`
	MatchesOpen     bool `schema:"m"`
}

func presetFineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			pass := r.FormValue("pass")
			realPass := os.Getenv("PASS")
			if pass != realPass {
				log.Printf("key Error fetching presetFineHandler")
				http.Error(w, "Not this time mate. [3]", http.StatusInternalServerError)
				return
			}

			if err := r.ParseForm(); err != nil {
				http.Error(w, "presetFineHandler - Invalid form data", http.StatusBadRequest)
				return
			}

			// Manual assignment of form values to struct
			presetFine := PresetFine{
				Reason:   r.FormValue("reason"),
				Context:  r.FormValue("context"),
				Icon:     r.FormValue("icon"),
				IsKudos:  r.FormValue("isKudos") == "on",
				Approved: true,
			}

			pfIdStr := r.FormValue("pfid")
			pfID, err := strconv.ParseUint(pfIdStr, 10, 64)
			if err != nil || pfID > 0 {
				presetFine.ID = uint(pfID)
			}

			displayOrderStr := r.FormValue("displayOrder")
			var displayOrder uint64 = 0
			if len(displayOrderStr) > 0 {
				displayOrder, err = strconv.ParseUint(displayOrderStr, 10, 64)
				if err == nil {
					presetFine.DisplayOrder = int64(displayOrder)
				}
			}

			amountStr := r.FormValue("amount")
			if len(amountStr) == 0 {
				amountStr = "2"
			}

			// Parse Amount as float64 from form value
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
				presetFine.Amount = amount
			} else {
				log.Printf("Error parsing amount: %v", err)
				http.Error(w, "Invalid amount value", http.StatusBadRequest)
				return
			}

			if err := SavePresetFine(db, &presetFine); err != nil {
				log.Printf("Error saving preset fine: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// var url = fmt.Sprintf("/finemaster/%s?pf=true&upf=%d#pf-%d", pass, presetFine.ID, presetFine.ID)

			editComp := editPresetFine(fmt.Sprintf("%s/%s", finemasterBaseUrl, pass), pass, presetFine, fmt.Sprintf("Updated standard fine - %v", time.Now().Format(time.TimeOnly)))
			editComp.Render(GetContext(r, db), w)
			return
		case "DELETE":
			pfIDStr := r.URL.Query().Get("pfid")

			pfID, err := strconv.ParseUint(pfIDStr, 10, 64)
			if err != nil || pfID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing playerId %v", err), http.StatusBadRequest)
				return
			}

			if err := DeletePresetFineByID(db, uint(pfID)); err != nil {
				log.Printf("Error deleting preset fine: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
			w.WriteHeader(http.StatusOK)
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

type FineMasterQueryParams struct {
	FinesOpen         bool `schema:"f"`
	PlayerOpen        bool `schema:"p"`
	PresetFinesOpen   bool `schema:"pf"`
	PresetFineUpdated uint `schema:"upf"`
	FineList          bool `schema:"fl"`
	MatchesOpen       bool `schema:"m"`
}

func finemasterAuthHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if googleAuthEnabled() {
				if _, _, ok := currentAdminUser(r, db); ok {
					w.Header().Set("HX-Redirect", "/finemaster")
					return
				}
				teams, err := GetTeams(db, 9999, 0)
				if err != nil {
					log.Printf("Error checking teams: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				authForm := googleAdminAuthForm(os.Getenv("GOOGLE_CLIENT_ID"), teams, getTeamId(GetContext(r, db)), r.URL.Query().Get("mode") == "register")
				authForm.Render(GetContext(r, db), w)
				return
			}
			// Check if user is already authenticated
			if _, err := r.Cookie("admin-user"); err == nil {
				// User is authenticated, redirect to finemaster
				w.Header().Set("HX-Redirect", "/finemaster")
				return
			}

			// Check if any teams exist
			teams, err := GetTeams(db, 1, 0)
			if err != nil {
				log.Printf("Error checking teams: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if len(teams) == 0 {
				// No teams exist, show bootstrap form
				bootstrap := bootstrapForm()
				bootstrap.Render(GetContext(r, db), w)
				return
			}

			// Show authentication form
			authForm := teamAuthForm()
			authForm.Render(GetContext(r, db), w)

		case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			if googleAuthEnabled() {
				handleGoogleAdminAuth(db, w, r)
				return
			}

			teamKey := r.FormValue("teamKey")
			adminPassword := r.FormValue("teamAdminPass")

			// Authenticate team
			team, err := GetTeamByKeyAndPassword(db, teamKey, adminPassword)
			if err != nil {
				log.Printf("Authentication failed for team key %s: %v", teamKey, err)
				// Show error message
				authForm := teamAuthForm()
				authForm.Render(GetContext(r, db), w)
				return
			}

			// Set authentication cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "admin-user",
				Value:    fmt.Sprintf("%d", team.ID),
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production with HTTPS
				SameSite: http.SameSiteStrictMode,
			})

			// Redirect to finemaster
			w.Header().Set("HX-Redirect", "/finemaster")
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func handleGoogleAdminAuth(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	info, err := googleCredentialVerifier(r.FormValue("credential"))
	if err != nil {
		log.Printf("Google admin auth failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		warning("Google sign-in failed").Render(GetContext(r, db), w)
		return
	}

	teamID := parseAdminTeamID(r, db)
	firstUserCount, err := CountAdminUsers(db)
	if err != nil {
		http.Error(w, "Could not check admin users", http.StatusInternalServerError)
		return
	}
	register := r.FormValue("mode") == "register"
	requestAccess := r.FormValue("mode") == "request"
	if requestAccess {
		if teamID == 0 {
			http.Error(w, "No active team for access request", http.StatusBadRequest)
			return
		}
		user, _, err := SaveAdminUserFromGoogle(db, *info, true)
		if err != nil {
			http.Error(w, "Could not save access request user", http.StatusInternalServerError)
			return
		}
		if err := RequestAdminAccess(db, user.ID, teamID, adminRoleTeamAdmin); err != nil {
			http.Error(w, "Could not request admin access", http.StatusInternalServerError)
			return
		}
		success("Access request sent. A team admin or super admin can approve it.").Render(GetContext(r, db), w)
		return
	}
	allowCreate := firstUserCount == 0
	if register && !allowCreate {
		team, err := GetTeam(db, teamID)
		if err == nil && team != nil && team.AllowAdminRegistration {
			allowCreate = true
		}
	}

	user, created, err := SaveAdminUserFromGoogle(db, *info, allowCreate)
	if err != nil {
		log.Printf("Admin user auth failed for %s: %v", info.Email, err)
		w.WriteHeader(http.StatusForbidden)
		warning("Admin account is not registered for this team").Render(GetContext(r, db), w)
		return
	}

	if created {
		if firstUserCount == 0 {
			if err := GrantAdminRole(db, user.ID, 0, adminRoleSuperAdmin); err != nil {
				http.Error(w, "Could not grant super admin role", http.StatusInternalServerError)
				return
			}
		} else {
			if teamID == 0 {
				http.Error(w, "No active team for registration", http.StatusBadRequest)
				return
			}
			if err := GrantAdminRole(db, user.ID, teamID, adminRoleTeamAdmin); err != nil {
				http.Error(w, "Could not grant team admin role", http.StatusInternalServerError)
				return
			}
		}
	}

	if teamID == 0 {
		if team, err := FirstTeamForAdminUser(db, user.ID); err == nil && team != nil {
			teamID = team.ID
		}
	}
	if teamID > 0 && !canAdminAccessTeam(db, user.ID, teamID) {
		http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
		return
	}

	setAdminSessionCookie(w, user.ID, teamID)
	if teamID > 0 {
		http.SetCookie(w, &http.Cookie{
			Name:     adminTokenCookieName,
			Value:    makeAdminToken(teamID, adminTokenSecret()),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   60 * 60 * 24 * 30,
		})
	}
	w.Header().Set("HX-Redirect", "/finemaster")
}

func parseAdminTeamID(r *http.Request, db *gorm.DB) uint {
	for _, key := range []string{"teamId", "teamID"} {
		if value := r.FormValue(key); value != "" {
			if id64, err := strconv.ParseUint(value, 10, 64); err == nil {
				return uint(id64)
			}
		}
	}
	return getTeamId(GetContext(r, db))
}

func canManageAdminRole(db *gorm.DB, viewer AdminUser, role string, teamID uint) bool {
	if isSuperAdmin(db, viewer.ID) {
		return true
	}
	return (role == adminRoleTeamAdmin || role == adminRoleLineupAccess) && teamID > 0 && canAdminAccessTeam(db, viewer.ID, teamID)
}

func adminUserRoleHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := currentAdminUser(r, db)
		if !ok {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		adminUserID64, err := strconv.ParseUint(r.FormValue("adminUserId"), 10, 64)
		if err != nil || adminUserID64 == 0 {
			http.Error(w, "Invalid admin user", http.StatusBadRequest)
			return
		}
		role := r.FormValue("role")
		action := strings.TrimSpace(r.FormValue("action"))
		if action == "" {
			action = "grant"
		}
		teamID := uint(0)
		switch role {
		case adminRoleSuperAdmin:
			teamID = 0
		case adminRoleTeamAdmin, adminRoleLineupAccess:
			teamID64, err := strconv.ParseUint(r.FormValue("teamId"), 10, 64)
			if err != nil || teamID64 == 0 {
				http.Error(w, "Team-scoped role requires a team", http.StatusBadRequest)
				return
			}
			teamID = uint(teamID64)
		default:
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}
		if !canManageAdminRole(db, *user, role, teamID) {
			http.Error(w, "Admin user cannot manage this role", http.StatusForbidden)
			return
		}

		if _, err := GetAdminUser(db, uint(adminUserID64)); err != nil {
			http.Error(w, "Admin user not found", http.StatusBadRequest)
			return
		}
		if teamID > 0 {
			if _, err := GetTeam(db, teamID); err != nil {
				http.Error(w, "Team not found", http.StatusBadRequest)
				return
			}
		}
		switch action {
		case "grant":
			if err := GrantAdminRole(db, uint(adminUserID64), teamID, role); err != nil {
				http.Error(w, "Could not assign admin role", http.StatusInternalServerError)
				return
			}
		case "revoke":
			if role == adminRoleSuperAdmin || role == adminRoleTeamAdmin {
				count, err := CountAdminUsersWithRole(db, teamID, role)
				if err != nil {
					http.Error(w, "Could not check admin role count", http.StatusInternalServerError)
					return
				}
				if count <= 1 {
					http.Error(w, "Cannot remove the last admin for this role", http.StatusBadRequest)
					return
				}
			}
			if err := RevokeAdminRole(db, uint(adminUserID64), teamID, role); err != nil {
				http.Error(w, "Could not remove admin role", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Invalid role action", http.StatusBadRequest)
			return
		}

		accounts, teams, requests, err := adminAccountsViewData(db, *user)
		if err != nil {
			http.Error(w, "Could not load admin users", http.StatusInternalServerError)
			return
		}
		message := "Role assigned"
		if action == "revoke" {
			message = "Role removed"
		}
		adminUserAccounts(accounts, teams, requests, isSuperAdmin(db, user.ID), success(message)).Render(GetContext(r, db), w)
	}
}

func adminAccessRequestHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, ok := currentAdminUser(r, db)
		if !ok {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		requestID64, err := strconv.ParseUint(r.FormValue("requestId"), 10, 64)
		if err != nil || requestID64 == 0 {
			http.Error(w, "Invalid access request", http.StatusBadRequest)
			return
		}
		action := strings.TrimSpace(r.FormValue("action"))
		if action != "approve" && action != "deny" {
			http.Error(w, "Invalid request action", http.StatusBadRequest)
			return
		}
		var request AdminAccessRequest
		if err := db.Where("id = ? AND status = ?", requestID64, "pending").First(&request).Error; err != nil {
			http.Error(w, "Access request not found", http.StatusBadRequest)
			return
		}
		if !canManageAdminRole(db, *user, request.Role, request.TeamID) {
			http.Error(w, "Admin user cannot manage this request", http.StatusForbidden)
			return
		}
		status := "denied"
		if action == "approve" {
			status = "approved"
		}
		resolved, err := ResolveAdminAccessRequest(db, uint(requestID64), user.ID, status)
		if err != nil {
			http.Error(w, "Could not update access request", http.StatusInternalServerError)
			return
		}
		if action == "approve" {
			if err := GrantAdminRole(db, resolved.AdminUserID, resolved.TeamID, resolved.Role); err != nil {
				http.Error(w, "Could not grant admin role", http.StatusInternalServerError)
				return
			}
		}
		accounts, teams, requests, err := adminAccountsViewData(db, *user)
		if err != nil {
			http.Error(w, "Could not load admin users", http.StatusInternalServerError)
			return
		}
		message := "Access request denied"
		if action == "approve" {
			message = "Access request approved"
		}
		adminUserAccounts(accounts, teams, requests, isSuperAdmin(db, user.ID), success(message)).Render(GetContext(r, db), w)
	}
}

func adminAccountsViewData(db *gorm.DB, viewer AdminUser) ([]AdminUserAccountView, []Team, []AdminAccessRequestView, error) {
	accounts, err := ListAdminUserAccounts(db)
	if err != nil {
		return nil, nil, nil, err
	}
	teams, err := GetTeams(db, 9999, 0)
	if err != nil {
		return nil, nil, nil, err
	}
	requests, err := ListPendingAdminAccessRequests(db)
	if err != nil {
		return nil, nil, nil, err
	}
	if isSuperAdmin(db, viewer.ID) {
		return accounts, teams, requests, nil
	}
	teams = FilterTeamsForAdminUser(db, viewer.ID, teams)
	allowedTeams := map[uint]bool{}
	for _, team := range teams {
		allowedTeams[team.ID] = true
	}
	requestUserIDs := map[uint]bool{}
	filteredRequests := []AdminAccessRequestView{}
	for _, request := range requests {
		if allowedTeams[request.Request.TeamID] {
			filteredRequests = append(filteredRequests, request)
			requestUserIDs[request.User.ID] = true
		}
	}
	filteredAccounts := []AdminUserAccountView{}
	for _, account := range accounts {
		filteredRoles := []AdminUserRoleView{}
		for _, role := range account.Roles {
			if (role.Role == adminRoleTeamAdmin || role.Role == adminRoleLineupAccess) && allowedTeams[role.TeamID] {
				filteredRoles = append(filteredRoles, role)
			}
		}
		if len(filteredRoles) > 0 || requestUserIDs[account.User.ID] {
			account.Roles = filteredRoles
			filteredAccounts = append(filteredAccounts, account)
		}
	}
	return filteredAccounts, teams, filteredRequests, nil
}

func presetFineMasterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presetFineMasterHandler")
		warnStr := ""
		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		activeSeason, _ := GetActiveSeason(db)
		seasonId := uint(0)
		if activeSeason != nil {
			seasonId = activeSeason.ID
		}

		pass := r.URL.Query().Get("pass")
		/*realPass := os.Getenv("PASS")
		pass := r.URL.Query().Get("pass")
		if pass != realPass {
			log.Printf("Error fetching presetFineMasterHandler - key miss match '%s' %s", pass, r.Referer())
			http.Error(w, "Not this time mate.", http.StatusInternalServerError)
			return
		}*/

		decoder := schema.NewDecoder()
		queryParams := new(FineMasterQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			warn := warning(fmt.Sprintf("presetFineMasterHandler - Error decoding query params: %v", err))
			warn.Render(ctx, w)
			return
		}

		if seasonId == 0 || teamId == 0 {
			warning("No active season/team - finemaster fines view is unavailable").Render(ctx, w)
			return
		}

		playersWithFines, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
		if err != nil {
			msg := fmt.Sprintf("Error fetching players with fines: %v", err)
			log.Print(msg)
			warn := warning(msg)
			warn.Render(ctx, w)
			return
		}

		pFines, err := GetPresetFines(db, true)
		if err != nil {
			msg := fmt.Sprintf("Error retrieving preset fines: %v", err)
			log.Print(msg)
			warn := warning(msg)
			warn.Render(ctx, w)
			return
		}

		matches, err := GetManageMatches(db, teamId, seasonId, 0, 9999)
		if err != nil {
			msg := fmt.Sprintf("Error retrieving preset fines: %v", err)
			log.Print(msg)
			warn := warning(msg)
			warn.Render(ctx, w)
			return
		}

		fineWithPlayers, err := GetFineWithPlayers(db, seasonId, teamId, 0, 999999)
		if err != nil {
			msg := fmt.Sprintf("Error retrieving GetFineWithPlayers: %v", err)
			log.Print(msg)
			warn := warning(msg)
			warn.Render(ctx, w)
			return
		}

		matchSeasonTeam, err := GetMatchSeasonTeam(db)
		if err != nil {
			msg := fmt.Sprintf("Error retrieving GetFineWithPlayers: %v", err)
			log.Print(msg)
			warn := warning(msg)
			warn.Render(ctx, w)
			return
		}
		if teamId > 0 && (matchSeasonTeam.Team == nil || matchSeasonTeam.Team.ID != teamId) {
			team, err := GetTeam(db, teamId)
			if err != nil {
				msg := fmt.Sprintf("Error retrieving selected team: %v", err)
				log.Print(msg)
				warn := warning(msg)
				warn.Render(ctx, w)
				return
			}
			matchSeasonTeam.Team = team
		}
		showAdminUsers := false
		adminAccounts := []AdminUserAccountView{}
		adminTeams := []Team{}
		adminRequests := []AdminAccessRequestView{}
		adminIsSuperAdmin := false
		if googleAuthEnabled() {
			if adminUser, _, ok := currentAdminUser(r, db); ok {
				showAdminUsers = true
				adminIsSuperAdmin = isSuperAdmin(db, adminUser.ID)
				var err error
				adminAccounts, adminTeams, adminRequests, err = adminAccountsViewData(db, *adminUser)
				if err != nil {
					msg := fmt.Sprintf("Error retrieving admin users: %v", err)
					log.Print(msg)
					warn := warning(msg)
					warn.Render(ctx, w)
					return
				}
			}
		}

		log.Print("presetFineMasterHandler - Rendering")

		backupStatus := LitestreamBackupStatus{}
		if showAdminUsers {
			backupStatus = litestreamBackupStatus()
		}
		finemaster := finemaster(pass, playersWithFines, fineWithPlayers, pFines, matches, *queryParams, matchSeasonTeam, warnStr, showAdminUsers, adminAccounts, adminTeams, adminRequests, adminIsSuperAdmin, backupStatus)
		log.Print("presetFineMasterHandler - Rendered")

		// set url to "/finemaster"
		//w.Header().Set("HX-Redirect", "/finemaster")
		//w.Header().Set("HX-Reload", "true")

		finemaster.Render(ctx, w)
		return
	}
}

func main() {
	// Initialize the database
	db, err := DBInit()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	} else {
		log.Printf("DB connected.")
	}

	// Initialize the router
	r := setupRouter(db)

	// Get the port from environment variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	portStr := fmt.Sprintf(":%s", port)
	log.Printf("Using port (%s)", portStr)

	// Start the server with graceful shutdown
	startServer(r, portStr)
}

func getUserInfo(r *http.Request) (string, error) {
	// Example: Extract username from a header
	username := r.Header.Get("X-Username")
	if username == "" {
		return "", errors.New("username not found in request")
	}
	return username, nil
}

func checkWhitelist(db *gorm.DB, username string) (bool, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func AuthMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserInfo(r)
		if err != nil {
			log.Printf("getUserInfo err (%+v)", err)
			http.Error(w, "Request Auth", http.StatusUnauthorized)
			return
		}

		checkUser, err := checkWhitelist(db, user)
		if err != nil {
			log.Printf("checkWhitelist err (%+v)", err)
			http.Error(w, "Request Auth", http.StatusUnauthorized)
			return
		}

		if !checkUser {
			http.Error(w, "Request Auth (checkUser)", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}

func requireTeamFeature(db *gorm.DB, moduleName string, enabled func(Team) bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := getTeamId(GetContext(r, db))
		if teamID == 0 {
			next(w, r)
			return
		}
		team, err := GetTeam(db, teamID)
		if err != nil {
			http.Error(w, "Active team not found", http.StatusNotFound)
			return
		}
		if !enabled(*team) {
			http.Error(w, fmt.Sprintf("%s module is disabled for %s", moduleName, team.TeamName), http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// setupRouter initializes the HTTP routes and returns a router.
func setupRouter(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	finesEnabled := func(team Team) bool { return team.FinesModuleEnabled() }
	lineupsEnabled := func(team Team) bool { return team.LineupsModuleEnabled() }
	matchesEnabled := func(team Team) bool { return team.EnableMatchesModule }
	playersEnabled := func(team Team) bool { return team.EnablePlayersModule }
	paymentsEnabled := func(team Team) bool { return team.PaymentsModuleEnabled() }
	courtEnabled := func(team Team) bool { return team.EnableCourtModule }

	r.HandleFunc("/players", requireTeamFeature(db, "Players", playersEnabled, playerHandler(db)))
	r.HandleFunc("/fines", requireTeamFeature(db, "Fines", finesEnabled, fineHandler(db)))
	r.HandleFunc("/fines/summary", requireTeamFeature(db, "Fines", finesEnabled, fineSummaryHandler(db)))

	r.HandleFunc("/admin", adminHandler(db))
	r.HandleFunc("/admin/users/roles", requireGoogleAdmin(db, adminUserRoleHandler(db)))
	r.HandleFunc("/admin/access-requests", requireGoogleAdmin(db, adminAccessRequestHandler(db)))
	r.HandleFunc("/admin/team-access-requests", adminTeamAccessRequestHandler(db))
	r.HandleFunc("/admin/team-switch", adminTeamSwitchHandler(db))

	r.HandleFunc("/fines/add", requireTeamFeature(db, "Fines", finesEnabled, fineAddHandler(db)))
	r.HandleFunc("/fines-multi", requireTeamFeature(db, "Fines", finesEnabled, fineMultiHandler(db)))
	r.HandleFunc("/fines/approve", requireTeamFeature(db, "Fines", finesEnabled, fineApproveHandler(db)))
	r.HandleFunc("/fines/edit/{fid}", requireTeamFeature(db, "Fines", finesEnabled, fineEditHandler(db)))
	r.HandleFunc("/fines/edit/{fid}/image", requireTeamFeature(db, "Fines", finesEnabled, fineImageHandler(db)))
	r.HandleFunc("/fines/contest", requireTeamFeature(db, "Fines", finesEnabled, fineContestHandler(db)))
	r.HandleFunc("/fines/context", requireTeamFeature(db, "Fines", finesEnabled, fineContextHandler(db)))
	r.HandleFunc("/fines/court-display-order", requireTeamFeature(db, "Court", courtEnabled, fineSetCourtSessionOrderHandler(db)))

	r.HandleFunc("/court", requireTeamFeature(db, "Court", courtEnabled, courtHandler(db)))

	r.HandleFunc("/preset-fines", requireTeamFeature(db, "Fines", finesEnabled, presetFineHandler(db)))
	r.HandleFunc("/preset-fines/approve", requireTeamFeature(db, "Fines", finesEnabled, presetFineApproveHandler(db)))
	r.HandleFunc("/preset-fines/{showOrHide}", requireTeamFeature(db, "Fines", finesEnabled, fineQuickHideHandler(db)))
	r.HandleFunc("/finemaster/auth", finemasterAuthHandler(db))
	r.HandleFunc("/finemaster", requireGoogleAdmin(db, presetFineMasterHandler(db)))
	r.HandleFunc("/finemaster/{pass}", requireGoogleAdmin(db, presetFineMasterHandler(db))) // Keep for backward compatibility
	r.HandleFunc("/", homeHandler(db))
	r.HandleFunc("/match-list", requireTeamFeature(db, "Matches", matchesEnabled, matchListHandler(db)))
	r.HandleFunc("/matches/history", requireTeamFeature(db, "Matches", matchesEnabled, matchHistoryHandler(db)))
	r.HandleFunc("/matches/sporty-sync", requireTeamFeature(db, "Matches", matchesEnabled, matchSportySyncHandler(db)))
	r.HandleFunc("/match/{matchId}", requireTeamFeature(db, "Matches", matchesEnabled, matchHandler(db)))
	r.HandleFunc("/match", requireTeamFeature(db, "Matches", matchesEnabled, matchHandler(db)))
	r.HandleFunc("/feedback", feedbackHandler(db))
	r.HandleFunc("/match-url/{matchURLSlug}", publicMatchURLHandler(db))
	r.HandleFunc("/match-url/{matchURLSlug}/feedback", publicMatchFeedbackHandler(db))
	r.HandleFunc("/notes", notesHandler(db))
	r.HandleFunc("/notes/context", notesContextHandler(db))
	r.HandleFunc("/notes/{noteId}", noteDetailHandler(db))
	r.HandleFunc("/lineups", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupListHandler(db)))
	r.HandleFunc("/lineups/login", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupLoginHandler(db)))
	r.HandleFunc("/lineups/new", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupNewHandler(db)))
	r.HandleFunc("/lineups/{lineupId}/share", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupShareHandler(db)))
	r.HandleFunc("/lineups/{lineupId}/errors", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupErrorsHandler(db)))
	r.HandleFunc("/lineups/{lineupId}/slot", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupSlotHandler(db)))
	r.HandleFunc("/lineups/{lineupId}/copy", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupCopyHandler(db)))
	r.HandleFunc("/lineups/{lineupId}", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, lineupDetailHandler(db)))
	r.HandleFunc("/formations", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, formationListHandler(db)))
	r.HandleFunc("/formations/new", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, formationNewHandler(db)))
	r.HandleFunc("/formations/{formationId}/live", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, formationLiveHandler(db)))
	r.HandleFunc("/formations/{formationId}/edit", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, formationEditHandler(db)))
	r.HandleFunc("/formations/{formationId}/copy", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, formationCopyHandler(db)))
	r.HandleFunc("/match-day", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, matchDayHandler(db)))
	r.HandleFunc("/match-day/{lineupId}/action", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, matchDayActionHandler(db)))
	r.HandleFunc("/match-day/{lineupId}", requireTeamFeature(db, "Line-ups and formations", lineupsEnabled, matchDayHandler(db)))

	r.HandleFunc("/playersName", requireTeamFeature(db, "Players", playersEnabled, playerNamesHandler(db)))
	r.HandleFunc("/season/{seasonId}/payments", requireTeamFeature(db, "Payments", paymentsEnabled, playerPayments(db)))

	r.HandleFunc("/match/{matchId}/event", requireTeamFeature(db, "Matches", matchesEnabled, matchEventHandler(db)))
	r.HandleFunc("/match/{matchId}/event/{eventId}", requireTeamFeature(db, "Matches", matchesEnabled, matchEventHandler(db)))
	r.HandleFunc("/match/{matchId}/events", matchEventListHandler(db))
	r.HandleFunc("/season/update/{updateType}", seasonBulkUpdateHandler(db))

	r.HandleFunc("/teams", teamHandler(db))
	r.HandleFunc("/teams/activate", requireGoogleAdmin(db, teamActivateHandler(db)))
	r.HandleFunc("/teams/sporty-sync", requireGoogleAdmin(db, sportySyncHandler(db)))
	r.HandleFunc("/teams/court-session", requireTeamFeature(db, "Court", courtEnabled, teamCourtHandler(db)))
	r.HandleFunc("/teams/active-match", requireGoogleAdmin(db, teamActiveMatchHandler(db)))

	r.HandleFunc("/season", seasonHandler(db))
	r.HandleFunc("/season/{seasonId}", seasonSpecificHandler(db))

	return r
}

func homeHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("homeHandler")
		if handleHomeTeamKey(w, r, db) {
			return
		}
		ctx := GetContext(r, db)
		teamId := getTeamId(ctx)
		activeSeason, _ := GetActiveSeason(db)
		seasonId := uint(0)
		if activeSeason != nil {
			seasonId = activeSeason.ID
		}

		decoder := schema.NewDecoder()
		queryParams := new(HomeQueryParams)
		warnStr := ""
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("presetFineMasterHandler - Error decoding query params: %v", err)
			http.Error(w, "Bad Request - home Decode", http.StatusBadRequest)
		}

		if seasonId == 0 || teamId == 0 {
			warnStr = warnStr + "No active season/team - showing empty fines\n"
		}

		playersWithFines, err := GetPlayersWithFines(db, seasonId, teamId, []uint64{})
		if err != nil {
			log.Printf("Error fetching players with fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		pFines, err := GetPresetFines(db, false)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var approvedPFines = []PresetFine{}
		var pendingPFines = []PresetFine{}

		for _, f := range pFines {
			if f.Approved {
				approvedPFines = append(approvedPFines, f)
			} else {
				pendingPFines = append(pendingPFines, f)
			}
		}

		fineWithPlayers, err := GetFineWithPlayers(db, seasonId, teamId, 0, 999999)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
		}

		mst, err := GetMatchSeasonTeam(db)
		if err != nil {
			warning := warning(fmt.Sprintf("Error fetching GetMatchSeasonTeam: %v", err))
			warning.Render(ctx, w)
			return
		}

		if mst.Team == nil {
			http.Redirect(w, r, "/teams?includeFrame=true", http.StatusSeeOther)
			return
		}
		// Ensure the rendered "active team" matches the user's admin selection (cookie).
		if teamId > 0 && (mst.Team == nil || mst.Team.ID != teamId) {
			if t, err := GetTeam(db, teamId); err == nil && t != nil {
				mst.Team = t
			}
		}
		matchDayLineupID := uint(0)
		if mst.Team != nil {
			if upcomingMatch, err := upcomingMatchForTeam(db, mst.Team.ID); err == nil {
				mst.Match = upcomingMatch
				matchDayLineupID = matchDayLineupIDForTeam(db, mst.Team.ID)
				if matchDayLineupID > 0 {
					mst.Match.LineupID = matchDayLineupID
				}
			} else {
				mst.Match = nil
			}
		}

		teams, err := GetTeams(db, 9999, 0)
		if err != nil {
			log.Printf("Error fetching teams: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		matches, err := GetMatches(db, seasonId, 0, 9999)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		previewPassword := ""
		if config.UsePreviewPassword {
			previewPassword = os.Getenv("PASS")
		}

		hasAdminToken := false
		selectedTeamID := uint(0)
		allowedTeams := []Team{}
		if googleAuthEnabled() {
			if adminUser, tid, ok := currentAdminUser(r, db); ok {
				hasAdminToken = true
				selectedTeamID = tid
				allowedTeams = FilterTeamsForAdminUser(db, adminUser.ID, teams)
			}
		} else if c, err := r.Cookie(adminTokenCookieName); err == nil {
			if tid, ok := parseAdminToken(c.Value, adminTokenSecret()); ok && tid > 0 {
				hasAdminToken = true
				selectedTeamID = tid
			}
		}
		if selectedTeamID == 0 && mst.Team != nil {
			selectedTeamID = mst.Team.ID
		}

		home := home(playersWithFines, approvedPFines, pendingPFines, fineWithPlayers, *queryParams, matches, mst, warnStr, previewPassword, teams, hasAdminToken, selectedTeamID, allowedTeams, matchDayLineupID, os.Getenv("GOOGLE_CLIENT_ID"))
		home.Render(ctx, w)
	}
}

func handleHomeTeamKey(w http.ResponseWriter, r *http.Request, db *gorm.DB) bool {
	teamKey := strings.TrimSpace(r.URL.Query().Get("teamKey"))
	if teamKey == "" {
		return false
	}
	team, err := GetTeamByKey(db, teamKey)
	if err != nil || team == nil {
		warning(fmt.Sprintf("Team not found for key %q", teamKey)).Render(GetContext(r, db), w)
		return true
	}
	if googleAuthEnabled() {
		if user, _, ok := currentAdminUser(r, db); ok {
			if canAdminAccessTeam(db, user.ID, team.ID) {
				if err := activateAdminTeam(w, *user, *team); err != nil {
					warning(fmt.Sprintf("Error saving active team: %v", err)).Render(GetContext(r, db), w)
					return true
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return true
			}
			teamKeyAccessPage(*team, os.Getenv("GOOGLE_CLIENT_ID"), true).Render(GetContext(r, db), w)
			return true
		}
		teamKeyAccessPage(*team, os.Getenv("GOOGLE_CLIENT_ID"), false).Render(GetContext(r, db), w)
		return true
	}
	session, _ := store.Get(r, "session-name")
	setTeamSessionValues(*team, session)
	if err := session.Save(r, w); err != nil {
		warning(fmt.Sprintf("Error saving active team: %v", err)).Render(GetContext(r, db), w)
		return true
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return true
}

// startServer starts the HTTP server with graceful shutdown.
func startServer(r *chi.Mux, portStr string) {
	server := &http.Server{
		Addr:    portStr,
		Handler: r,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s", portStr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", portStr, err)
		}
	}()

	// Block until we receive a signal
	<-stop

	// Create a context with a timeout to allow graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func GenerateURLWithParams(baseURL string, params *HomeQueryParams) (string, error) {
	// Parse the base URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Initialize the encoder
	encoder := schema.NewEncoder()

	// Prepare a map to encode the parameters into
	values := url.Values{}

	// Encode the struct into the map
	err = encoder.Encode(params, values)
	if err != nil {
		return "", err
	}

	// Append the query parameters to the URL
	u.RawQuery = values.Encode()

	return u.String(), nil
}
