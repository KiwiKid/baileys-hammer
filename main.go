package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var decoder = schema.NewDecoder()

type Context struct {
    Title   string
	UseRoles bool
	UseMatchEventTracker bool
}


func playerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch(r.Method){
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
					playersWithFines, err := GetPlayersWithFines(db, playerIds)
					if err != nil || len(playersWithFines) == 0 {
						log.Printf("Error fetching players with fines: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}

					playerList := playerRoleSelector(playersWithFines[0], config, "")
					playerList.Render(GetContext(r), w)
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

						playerList := playerInputSelector(players, playerId, inputType)
						playerList.Render(GetContext(r), w)
				}

				
				
				

			case "POST": {
				if err := r.ParseForm(); err != nil {
					log.Printf("Error parsing form data: %v", err)
					http.Error(w, "Bad Request savePlayerHandler ParseForm", http.StatusBadRequest)
					return
				}
	
				var player Player
				// Use the decoder to populate the player struct
				if err := decoder.Decode(&player, r.PostForm); err != nil {
					log.Printf("Error decoding form into player struct: %v", err)
					http.Error(w, "Bad Request savePlayerHandler Decode", http.StatusBadRequest)
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
				playersWithFines, err := GetPlayersWithFines(db, playerIds)
				if err != nil || len(playersWithFines) == 0 {
					log.Printf("Error fetching players with fines: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				playerList := playerRoleSelector(playersWithFines[0], config, fmt.Sprintf("Updated player"))
				playerList.Render(GetContext(r), w)
				return
	
			}
		case "DELETE": {
			log.Println("PlayerHanlder")
			if err := r.ParseForm(); err != nil {
				log.Printf("Error parsing form data: %v", err)
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
		success.Render(GetContext(r), w)
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
		var matchId uint64 = 0;
		var fineAt *time.Time;
		if matchIdStr != "NA" {
			matchId, err = strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid matchId ID - %s", matchIdStr), http.StatusBadRequest)
				return
			}

			match, getMatchErr := GetMatch(db, matchId);
			if getMatchErr != nil {
				http.Error(w, fmt.Sprintf("Invalid matchId ID - %s", matchIdStr), http.StatusBadRequest)
				return
			} else {
				fineAt = match.StartTime
			}
		}

		context := r.FormValue("context")

		
		

		err = UpdateFineContextByID(db, uint(fineID), uint(matchId), context, fineAt)
		if err != nil {
			http.Error(w, "Invalid UpdateFineContextByID ID", http.StatusBadRequest)
			return
		}
		
		success := contextSuccess(matchId, context, fineAt)
		success.Render(GetContext(r), w)
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
				match, err = GetMatch(db, uint64(fine.MatchId))
				if err != nil {
					http.Error(w, fmt.Sprintf("GetMatch Match not found - %d", fine.MatchId), http.StatusNotFound)
					return
				}
			}

            // Prepare the data for rendering
            fineWithPlayer := FineWithPlayer{
                Fine:   *fine,
                Player: *player,
				Match:  *match,
            }

			isEdit := r.URL.Query().Get("isEdit")
			isContest := r.URL.Query().Get("isContest")
			isContext := r.URL.Query().Get("isContext")

			isFineMasterStr := r.URL.Query().Get("isFineMaster")

			var isFineMaster bool = isFineMasterStr == "true"
			if(isFineMaster) {
				log.Printf("permision Error fetching fineEditHandler")
				http.Error(w, "Not this time mate.", http.StatusInternalServerError)
				return
			}


			if (isEdit == "true") {
				fineEditRow := fineEditRow(fineWithPlayer, isFineMaster)
				fineEditRow.Render(GetContext(r), w)
			} else if(isEdit == "form") {
				fineFormRow := fineEditForm(fineWithPlayer, isFineMaster)
				fineFormRow.Render(GetContext(r), w)
			} else if(isContest == "true") {
				fineContestRow := fineContestRow(fineWithPlayer)
				fineContestRow.Render(GetContext(r), w)
			} else if(isContext == "true") {
				matches, err := GetMatches(db, 1, 0, 9999)
				if err != nil {
					http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
					return
				}
				fineContestRow := fineContextRow(fineWithPlayer, matches)
				fineContestRow.Render(GetContext(r), w)
			} else {
				fineRowComp := fineRow(isFineMaster, fineWithPlayer)
				fineRowComp.Render(GetContext(r), w)
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
			if(len(approvedStr) == 0){
				approved = config.DefaultToApproved
			} else {
				approved = approvedStr == "true" //todo: on?
			}



            playerId, err := strconv.ParseUint(r.FormValue("playerId"), 10, 64)
            if err != nil {
				errComp := errMsg(F("Invalid playerId", r.FormValue("playerId")))
				errComp.Render(GetContext(r), w)
            }

			matchId, err := strconv.ParseUint(r.FormValue("matchId"), 10, 64)
            if err != nil {
				errComp := errMsg(F("Invalid matchId", r.FormValue("matchId")))
				errComp.Render(GetContext(r), w)
            }

			context := r.FormValue("context")

			// Update the fine in the database
			fine := Fine{
				Model:    gorm.Model{ID: uint(fineID)},
				PlayerID: uint(playerId),
				Amount:   amount,
				Reason:   reason,
				Context: context,
				Approved: approved,
			}

			
			if matchId != 0 {
				match, err := GetMatch(db, matchId)
				if err != nil {
					errComp := errMsg("Cannot select get the active match?")
					errComp.Render(GetContext(r), w)
				} else if match != nil {
					fine.FineAt = *match.StartTime
				}
			} else {
				activeMatch, actMatchERr := GetActiveMatch(db)
				if actMatchERr != nil {
					errComp := errMsg("Cannot select get the active match?")
					errComp.Render(GetContext(r), w)
				}else if activeMatch != nil {
					fine.FineAt = *activeMatch.StartTime
				}
			}
			
            if err := SaveFine(db, &fine); err != nil {
                http.Error(w, "Failed to update fine", http.StatusInternalServerError)
                return
            }else {
				log.Printf("\n\nSAve fine %+v\n\n", fine)
			}

			player, err := GetPlayerByID(db, fine.PlayerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
				return
			}

			match := &Match{}
			if fine.MatchId > 0 {
				match, _ = GetMatch(db, uint64(fine.MatchId))
			}

			// Prepare the data for rendering
			fineWithPlayer := FineWithPlayer{
				Fine:   fine,
				Player: *player,
				Match: *match,
			}
				
			fineRowComp := fineRow(true, fineWithPlayer)
			fineRowComp.Render(GetContext(r), w)

        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }
}


type FineWithPlayer struct {
	Fine Fine
	Player Player
	Match Match
}

func fineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch(r.Method){
		case "GET": {
			var pageId = 0
			pageStr := r.URL.Query().Get("page")
			if len(pageStr) == 0 {
				pageId = 0
			}else {
				pageIdUint, err := strconv.ParseInt(pageStr, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
					return 
				}
				pageId = int(pageIdUint)
			}

			limitStr := r.URL.Query().Get("limit")
			var limit = 999999
			if len(limitStr) == 0{
				limit = 999999
			}else {
				limit, err := strconv.ParseInt(limitStr, 10, 64)
				if err != nil || limit == 0 {
					http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
					return 
				}
			}
		
		finemasterPage := false
		splitUrl := strings.Split(r.Header.Get("Referer"), "/")
		for _, urlBit := range splitUrl {
			if(urlBit == "finemaster"){
				finemasterPage = true
			}
		}

			fineWithPlayers, err := GetFineWithPlayers(db, 0, limit)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
				return 
			}

			fineList := fineList(fineWithPlayers, pageId, 0, finemasterPage, false)
			fineList.Render(GetContext(r), w)
		}
			case "POST": {

				if err := r.ParseForm(); err != nil {
					http.Error(w, "fineHandler - Invalid form data", http.StatusBadRequest)
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

				if(len(playerIdStr) == 0 || r.FormValue("fineOption") == "applyAgain") {

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
							Reason: reason,
							Approved: false,
							Amount: amount,
							Context: context,
							DisplayOrder: displayOrder,
							IsKudos: isKudos,
							Icon: icon,
						}
						err := SavePresetFine(db, suggestedPFine)
						if err != nil {
							// Handle the error if the conversion fails
							http.Error(w, "SavePresetFine failed", http.StatusBadRequest)
							return
						}else{
							createdPFines = append(createdPFines, *suggestedPFine)
						}
					}
				}
				if (len(playerIdStr) > 0) {
				

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
							errComp.Render(GetContext(r), w)
						}
						if pfIdStr == "-1" && len(presetFineIds) > 1 {
							errComp := errMsg("Cannot select \"Fine is not listed here\" with others")
							errComp.Render(GetContext(r), w)
						}

						var fineAt = time.Now()
						activeMatch, actMatchERr := GetActiveMatch(db)
						if actMatchERr != nil {
							errComp := errMsg("Cannot select get the active match?")
							errComp.Render(GetContext(r), w)
						}else if activeMatch != nil {
							fineAt = *activeMatch.StartTime
						}
	
						if pfIdStr == "-1" {

							fine := Fine{
								Amount: amount,
								Reason: reason,
								Context: context,
								PlayerID: uint(playerId),
								FineAt: fineAt,
								Approved: approved,
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
						
							
							if(len(context) == 0) {
								context = presetFine.Context
							}
							
							fine := Fine{
								Amount: presetFine.Amount,
								Reason: presetFine.Reason,
								FineAt: time.Now(),
								Context: context,
								PlayerID: uint(playerId),
								Approved: approved,
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
				success.Render(GetContext(r), w)
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
		var newFines []uint
        if(r.Method == "POST"){
		
			var input struct {
				Reason string `json:"reason"`
			}
			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var suggestedPFine = &PresetFine{
				Reason: input.Reason,
				Approved: true,
				Amount: 2,
			}
			err := SavePresetFine(db, suggestedPFine)
			if err != nil {
				http.Error(w, "SavePresetFine failed", http.StatusBadRequest)
				return
			}
			newFines = append(newFines, suggestedPFine.ID)

			//success := success(fmt.Sprintf("Added Fine - %s", input.Reason))
			//success.Render(GetContext(r), w)
		}



		playersWithFines, err := GetPlayersWithFines(db, []uint64{})
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
		fsComp.Render(GetContext(r), w)
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

		log.Printf("pfineIDs: %+v %d", len(pfineIDs), len(playerIDs))
		


        for _, pfineIDStr := range pfineIDs {
			

                for _, playerIDStr := range playerIDs {
					playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
					if err != nil {
						errComp := errMsg(F("Invalid player ID: [%s]", playerIDStr))
						errComp.Render(GetContext(r), w)
						return
					}
					fine, err :=  GetFineFromPreset(db, pfineIDStr)
					if err != nil {
						errComp := errMsg(F("Could not GetFineFromPreset: [%s]", pfineIDStr))
						errComp.Render(GetContext(r), w)
						return
					}
					fine.PlayerID = uint(playerID)
					fine.Context = context

					activeMatch, err := GetActiveMatch(db)
					if err != nil {
						errComp := errMsg(F("Could not get active match %v", err))
						errComp.Render(GetContext(r), w)
					}

					if(activeMatch != nil){
						fine.MatchId = activeMatch.ID
						if(activeMatch.StartTime != nil){
							fine.FineAt = *activeMatch.StartTime
						}else {
							fine.FineAt = time.Now()
						}
					}

					fine.CreatedAt = time.Now()
					err = SaveFine(db, fine)
					if err != nil {
						http.Error(w, "Invalid player ID", http.StatusBadRequest)
						errComp := errMsg(F("Could not Save Fine %+v", fine))
						errComp.Render(GetContext(r), w)
						return
					} else {
						savedFines = append(savedFines, *fine)
					}
            	}
        }

		playersWithFines, err := GetPlayersWithFines(db, []uint64{})
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

		res := fineSuperSelectResults(playersWithFines, pFines, savedFines)
		res.Render(GetContext(r), w)
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
				success.Render(GetContext(r), w)	
			} else {
				success := warning("Declined")
				success.Render(GetContext(r), w)	
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
		switch(r.Method){
			case "POST": {

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
			success.Render(GetContext(r), w)
			return
			}
			default: 
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
	}
	}
}

func GenerateUrl(baseURL string, queryParams *HomeQueryParams) (*string, error, ) {
	fullURL, err := GenerateURLWithParams(baseURL, queryParams)
    if err != nil {
        return nil, err
    }
	return &fullURL, nil
}

type HomeQueryParams struct {
    FinesOpen bool `schema:"f"`
	FineListOpen bool `schema:"fl"`
	PlayerOpen bool `schema:"p"`
	PresetFinesOpen bool `schema:"pf"`
	MatchesOpen bool `schema:"m"`
}

func presetFineHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
			pass := r.FormValue("pass")
			realPass := os.Getenv("PASS")
			if(pass != realPass) {
				log.Printf("key Error fetching presetFineHandler")
				http.Error(w, "Not this time mate.", http.StatusInternalServerError)
				return
			}
			
			if err := r.ParseForm(); err != nil {
				http.Error(w, "presetFineHandler - Invalid form data", http.StatusBadRequest)
				return
			}

			

			// Manual assignment of form values to struct
			presetFine := PresetFine{
				Reason: r.FormValue("reason"),
				Context: r.FormValue("context"),
				Icon: r.FormValue("icon"),
				IsKudos: r.FormValue("isKudos") == "on",
				Approved: true,
			}

			pfIdStr := r.FormValue("pfid")
			pfID, err := strconv.ParseUint(pfIdStr, 10, 64)
			if err != nil || pfID > 0 {
				presetFine.ID = uint(pfID);
			}

			displayOrderStr := r.FormValue("displayOrder")
			var displayOrder uint64 = 0;
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
			editComp.Render(GetContext(r), w)
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
	FinesOpen bool `schema:"f"`
	PlayerOpen bool `schema:"p"`
	PresetFinesOpen bool `schema:"pf"`
	PresetFineUpdated uint `schema:"upf"`
	FineList bool `schema:"fl"`
	MatchesOpen bool `schema:"m"`
}

func presetFineMasterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := chi.URLParam(r, "pass")

		/*realPass := os.Getenv("PASS")
		if(pass != realPass) {
			log.Printf("Error fetching presetFineMasterHandler - key miss match %s %s", pass, r.Referer())
			http.Error(w, "Not this time mate.", http.StatusInternalServerError)
			return
		}*/

		decoder := schema.NewDecoder()
		queryParams := new(FineMasterQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("presetFineMasterHandler - Error decoding query params: %v", err)
			http.Error(w, "Bad Request presetFineMasterHandler", http.StatusBadRequest)
			return
		}

		playersWithFines, err := GetPlayersWithFines(db, []uint64{})
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

		matches, err := GetMatches(db, 1, 0, 9999)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fineWithPlayers, err := GetFineWithPlayers(db, 0, 999999)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
			return 
		}

		finemaster := finemaster(pass, playersWithFines, fineWithPlayers, pFines, matches, *queryParams)
		finemaster.Render(GetContext(r), w)
	}
}



func main() {

	
	log.Printf("Started")
	db, err := DBInit()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	} else {
		log.Printf("DB connected.")
	}

	r := chi.NewRouter()

	r.HandleFunc("/players", playerHandler(db))
	//r.HandleFunc("/players/{playerId}", playerHandler(db))
	r.HandleFunc("/fines", fineHandler(db))
	r.HandleFunc("/fines/add", fineAddHandler(db))
	r.HandleFunc("/fines-multi", fineMultiHandler(db))
	r.HandleFunc("/fines/approve", fineApproveHandler(db))
	r.HandleFunc("/fines/edit/{fid}", fineEditHandler(db))
	r.HandleFunc("/fines/contest", fineContestHandler(db))
	r.HandleFunc("/fines/context", fineContextHandler(db))
	r.HandleFunc("/preset-fines", presetFineHandler(db))
	r.HandleFunc("/preset-fines/approve", presetFineApproveHandler(db))
	r.HandleFunc("/preset-fines/{showOrHide}", fineQuickHideHandler(db))
	r.HandleFunc("/finemaster/{pass}", presetFineMasterHandler(db))
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := schema.NewDecoder()
		queryParams := new(HomeQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("presetFineMasterHandler - Error decoding query params: %v", err)
			http.Error(w, "Bad Request - home Decode", http.StatusBadRequest)
			return
		}

		playersWithFines, err := GetPlayersWithFines(db, []uint64{})
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
			}else{
				pendingPFines = append(pendingPFines, f)
			}
		}

		fineWithPlayers, err := GetFineWithPlayers(db, 0, 999999)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
			return 
		}

		activeMatch, err := GetActiveMatch(db)
		if err != nil {
			errComp := errMsg(F("Could not get active match %v", err))
			errComp.Render(GetContext(r), w)
		}


		matches, err := GetMatches(db, 1, 0, 9999)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		home := home(playersWithFines, approvedPFines, pendingPFines, fineWithPlayers, *queryParams, matches, activeMatch)
		home.Render(GetContext(r), w)
	})

	log.Printf("Router 1/2")
	r.HandleFunc("/match-list", matchListHandler(db))
	r.HandleFunc("/match/{matchId}", matchHandler(db))
	r.HandleFunc("/match", matchHandler(db))

	r.HandleFunc("/playersName", playerNamesHandler(db))

	r.HandleFunc("/match/{matchId}/event", matchEventHandler(db))
	r.HandleFunc("/match/{matchId}/event/{eventId}", matchEventHandler(db))
	r.HandleFunc("/match/{matchId}/events", matchEventListHandler(db))
	log.Printf("Router 2/2")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Router 3/4")

		log.Panic("Server error: %v", err)
	}else {
		log.Printf("Router 3/3")

		log.Printf("\n=================\n\n    Listening on\n\n================\n%d", 8080)	
	}
	log.Printf("Router 4/4")
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