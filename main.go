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

func playerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch(r.Method){
			case "GET":
				playerIdStr := r.URL.Query().Get("playerId")
				playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
				if err != nil || playerId == 0 {
					http.Error(w, fmt.Sprintf("Error playerHandler playerId %v", err), http.StatusBadRequest)
					return 
				}

				player, err := GetPlayerByID(db, uint(playerId))
				if err != nil {
					http.Error(w, fmt.Sprintf("Error playerHandler GetPlayerByID %v", err), http.StatusBadRequest)
					return 
				}

				
				playerList := playerName(*player)
				playerList.Render(r.Context(), w)

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

				w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))

				// Optionally, you can set the status code to 200 OK or any appropriate status
				w.WriteHeader(http.StatusOK)
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
		success.Render(r.Context(), w)
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
		matchId, err := strconv.ParseUint(matchIdStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid matchId ID - %s", matchIdStr), http.StatusBadRequest)
			return
		}

		context := r.FormValue("context")

		err = UpdateFineContextByID(db, uint(fineID), uint(matchId), context)
		if err != nil {
			http.Error(w, "Invalid UpdateFineContextByID ID", http.StatusBadRequest)
			return
		}

		success := success(fmt.Sprintf("Added Context - m:%d - %s", matchId, context))
		success.Render(r.Context(), w)
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
            } else { 
				log.Printf("GOT FINE: \n\n %+v", fine)
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

			


			if (isEdit == "true") {
				fineEditRow := fineEditRow(fineWithPlayer)
				fineEditRow.Render(r.Context(), w)
			} else if(isContest == "true") {
				fineContestRow := fineContestRow(fineWithPlayer)
				fineContestRow.Render(r.Context(), w)
			} else if(isContext == "true") {
				matches, err := GetMatches(db, 1, 0, 9999)
				if err != nil {
					http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
					return
				}
				fineContestRow := fineContextRow(fineWithPlayer, matches)
				fineContestRow.Render(r.Context(), w)
			} else {
				fineRowComp := fineRow(true, fineWithPlayer)
				fineRowComp.Render(r.Context(), w)
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
            approved := r.FormValue("approved") == "true"
            playerId, err := strconv.ParseUint(r.FormValue("playerId"), 10, 64)
            if err != nil {
                http.Error(w, "Invalid playerId ID", http.StatusBadRequest)
                return
            }

			fineAtStr := r.FormValue("fineAt")
			fineAt, err := time.Parse("2006-01-02", fineAtStr)
			if err != nil {
				// Handle parsing error, perhaps set a default value or return an error response
				http.Error(w, "Failed to parse fineAt time", http.StatusBadRequest)
				return
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
				FineAt: fineAt,
            }

            if err := SaveFine(db, &fine); err != nil {
                http.Error(w, "Failed to update fine", http.StatusInternalServerError)
                return
            }

			player, err := GetPlayerByID(db, fine.PlayerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Player not found - %d", fine.PlayerID), http.StatusNotFound)
				return
			}

			match, _ := GetMatch(db, uint64(fine.MatchId))
			

			// Prepare the data for rendering
			fineWithPlayer := FineWithPlayer{
				Fine:   fine,
				Player: *player,
				Match: *match,
			}
				
			fineRowComp := fineRow(true, fineWithPlayer)
			fineRowComp.Render(r.Context(), w)

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
			var limit = 50
			if len(limitStr) == 0{
				limit = 50
			}else {
				limit, err := strconv.ParseInt(limitStr, 10, 64)
				if err != nil || limit == 0 {
					http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
					return 
				}
			}
			fines, getFErr := FetchLatestFines(db, int(pageId), int(limit))
			if getFErr != nil {
				http.Error(w, "fineHandler - FetchLatestFines", http.StatusBadRequest)
				return
			}

			players, getPlayerErr := GetPlayers(db, 0, 1000)
			if getPlayerErr != nil {
				http.Error(w, "fineHandler - GetPlayers error", http.StatusBadRequest)
				return
			}

			// Get all relevant matches
			matches, getMatchErr := GetMatches(db, 0, 0, 1000)
			if getMatchErr != nil {
				http.Error(w, "fineHandler - GetMatches error", http.StatusBadRequest)
				return
			}

			matchMap := make(map[uint]Match)
			for _, match := range matches {
				matchMap[match.ID] = match
			}

			finemasterPage := false
			splitUrl := strings.Split(r.Header.Get("Referer"), "/")
			for _, urlBit := range splitUrl {
				if(urlBit == "finemaster"){
					finemasterPage = true
				}
			}
			w.WriteHeader(http.StatusOK)

			var fineWithPlayers []FineWithPlayer

			for _, fine := range fines {
				var matchedPlayer Player
				for _, player := range players {
					if fine.PlayerID == player.ID {
						matchedPlayer = player 
						break
					}
				}

				matchedMatch := matchMap[fine.MatchId]
				
				fineWithPlayers = append(fineWithPlayers, FineWithPlayer{
					Fine:   fine,
					Player: matchedPlayer,
					Match: matchedMatch,
				})
			}
			

			fineList := fineList(fineWithPlayers, pageId, finemasterPage)
			fineList.Render(r.Context(), w)
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


				log.Printf("WOAH WOAH WOAH WOAH playerIdStr %d fineOption: %s", len(playerIdStr), r.FormValue("fineOption") )

				if(len(playerIdStr) == 0 || r.FormValue("fineOption") == "applyAgain") {

					amountStr := r.FormValue("amount")
					if len(amountStr) == 0 {
						amountStr = "1"
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
					log.Printf("WOAH WOAH WOAH WOAH %s", playerIdStr)
				

					playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
					if err != nil {
						log.Printf("Error get playerIdStr fine: %v", err)
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
						if len(amountStr) == 0 {
							amountStr = "1"
						}
						amount, err := strconv.ParseFloat(amountStr, 64) // 64 specifies the bit size of the float type
						if err != nil {
							// Handle the error if the conversion fails
							errComp := errMsg("Invalid amount")
							errComp.Render(r.Context(), w)
						}
						if pfIdStr == "-1" && len(presetFineIds) > 1 {
							errComp := errMsg("Cannot select \"Fine is not listed here\" with others")
							errComp.Render(r.Context(), w)
						}
	
						if pfIdStr == "-1" {

							fine := Fine{
								Amount: amount,
								Reason: reason,
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
						} else {

							// Parse the string ID to an unsigned integer
							pfId, err := strconv.ParseUint(pfIdStr, 10, 64)
							if err != nil {
								http.Error(w, "Invalid strconv.ParseUint(pfIdStr", http.StatusBadRequest)
								return
							}
						
							// Assuming a function GetPresetFineByID that returns a *PresetFine struct for a given ID
							presetFine, err := GetPresetFine(db, pfId)
							if err != nil {
								http.Error(w, "Invalid GetPresetFine data", http.StatusBadRequest)
								return
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
				
				log.Printf("fineAddRes fineAddRes fineAddRes WOAH")

				success := fineAddRes(createdFines, createdPFines)
				success.Render(r.Context(), w)
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
        switch(r.Method) {
		case "POST": {
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
			}
			err := SavePresetFine(db, suggestedPFine)
			if err != nil {
				http.Error(w, "SavePresetFine failed", http.StatusBadRequest)
				return
			}

			// Assuming you have a function to write JSON response
			success := success(fmt.Sprintf("Added Fine - %s", input.Reason))
			success.Render(r.Context(), w)
		}
	case "GET": {

		playersWithFines, err := GetPlayersWithFines(db)
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

		fsComp := fineSuperSelect(playersWithFines, pFines)
		fsComp.Render(r.Context(), w)

	}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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
		savedFines := []Fine{}

		log.Printf("pfineIDs: %+v %d", len(pfineIDs), len(playerIDs))
		


        for _, pfineIDStr := range pfineIDs {
			

                for _, playerIDStr := range playerIDs {
					playerID, err := strconv.ParseUint(playerIDStr, 10, 64)
					if err != nil {
						http.Error(w, "Invalid player ID", http.StatusBadRequest)
						return
					}
					fine, err :=  GetFineFromPreset(db, pfineIDStr)
					if err != nil {
						http.Error(w, "Invalid player ID", http.StatusBadRequest)
						return
					}
					fine.PlayerID = uint(playerID)


					err = SaveFine(db, fine)
					if err != nil {
						http.Error(w, "Invalid player ID", http.StatusBadRequest)
						return
					} else {
						savedFines = append(savedFines, *fine)
						log.Printf("FINE\n\n%+v", fine)
					}
            	}
        }

		playersWithFines, err := GetPlayersWithFines(db)
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
		res.Render(r.Context(), w)
    }
}


func fineApproveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// Access form values
			fIDStr := r.FormValue("fid") // Now getting 'fid' from the form data
			amountStr := r.FormValue("amount")

			// Convert 'fid' to uint64
			fID, err := strconv.ParseUint(fIDStr, 10, 64)
			if err != nil || fID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine ID: %v", err), http.StatusBadRequest)
				return
			}

			// Convert 'amount' to float64
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing amount: %v", err), http.StatusBadRequest)
				return
			}
			log.Printf("Approve fine %d %f", fID, amount)
			// Call ApproveFine with parsed 'fid' and 'amount'
			if err := ApproveFine(db, uint(fID), amount); err != nil {
				log.Printf("Error approving fine with specified amount: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Redirect or handle response as needed
			//w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
			//w.WriteHeader(http.StatusOK)

			success := success("Approved")
			success.Render(r.Context(), w)
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
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			pfIDStr := r.FormValue("pfid")

			// Convert 'fid' to uint64
			pfID, err := strconv.ParseUint(pfIDStr, 10, 64)
			if err != nil || pfID == 0 {
				http.Error(w, fmt.Sprintf("Error parsing fine ID: %v", err), http.StatusBadRequest)
				return
			}

			if err := QuickHideFine(db, uint(pfID)); err != nil {
				log.Printf("Error approving fine with specified amount: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Redirect or handle response as needed
			w.Header().Set("HX-Redirect", r.Header.Get("Referrer")+"#preset-fine")
			w.Header().Set("HX-Reload", "true")
			w.WriteHeader(http.StatusOK)
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
			success.Render(r.Context(), w)
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
			if err := r.ParseForm(); err != nil {
				http.Error(w, "presetFineHandler - Invalid form data", http.StatusBadRequest)
				return
			}
	
			// Manual assignment of form values to struct
			presetFine := PresetFine{
				Reason: r.FormValue("reason"),
				Approved: true,
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

			// Save PresetFine to the database
			if err := SavePresetFine(db, &presetFine); err != nil {
				log.Printf("Error saving preset fine: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}


			w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))

			// Optionally, you can set the status code to 200 OK or any appropriate status
			w.WriteHeader(http.StatusOK)
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
	FineList bool `schema:"fl"`
	MatchesOpen bool `schema:"m"`
}

func presetFineMasterHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := chi.URLParam(r, "pass")

		realPass := os.Getenv("PASS")
		log.Printf("presetFineMasterHandler: %s (%s)", pass, realPass)

		if(pass != realPass) {
			log.Printf("Error fetching presetFineMasterHandler")
			http.Error(w, "Not this time mate.", http.StatusInternalServerError)
			return
		}

		decoder := schema.NewDecoder()
		queryParams := new(FineMasterQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("Error decoding query params: %v", err)
			http.Error(w, "Bad Request presetFineMasterHandler", http.StatusBadRequest)
			return
		}

		playersWithFines, err := GetPlayersWithFines(db)
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
		finemaster := finemaster(pass, playersWithFines, pFines, matches, *queryParams)
		finemaster.Render(r.Context(), w)
	}
}

func main() {
	log.Printf("Started")
	db, err := DBInit()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
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
	r.HandleFunc("/preset-fines/hide", fineQuickHideHandler(db))
	r.HandleFunc("/finemaster/{pass}", presetFineMasterHandler(db))
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := schema.NewDecoder()
		queryParams := new(HomeQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("Error decoding query params: %v", err)
			http.Error(w, "Bad Request - home Decode", http.StatusBadRequest)
			return
		}

		playersWithFines, err := GetPlayersWithFines(db)
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

		home := home(playersWithFines, approvedPFines, pendingPFines, *queryParams)
		home.Render(r.Context(), w)
	})

	r.HandleFunc("/match-list", matchListHandler(db))
	r.HandleFunc("/match/{matchId}", matchHandler(db))
	r.HandleFunc("/match", matchHandler(db))

	r.HandleFunc("/match/{matchId}/event", matchEventHandler(db))
	r.HandleFunc("/match/{matchId}/events", matchEventListHandler(db))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Printf("Server error: %v", err)
	}
	log.Printf("Listening on %d", 8080)
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