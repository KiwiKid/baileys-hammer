package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

				
				fineList := playerName(*player)
				fineList.Render(r.Context(), w)

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

type FineWithPlayer struct {
	Fine Fine
	Player Player
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
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			players, getPlayerErr := FetchPlayers(db, 0, 100)
			if getPlayerErr != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			finemasterPage := false
			splitUrl := strings.Split(r.Header.Get("Referer"), "/")
			for _, urlBit := range splitUrl {
				if(urlBit == "finemaster"){
					finemasterPage = true
				}
			}
			w.WriteHeader(http.StatusOK)

			
			var fineWithPlayers = []FineWithPlayer{}
			for _, fine := range fines {
				for _, player := range players {
					if(fine.PlayerID == player.ID){
						fineWithPlayers = append(fineWithPlayers, FineWithPlayer{ 
							Fine: fine,
							Player: player,
						})
						continue;
					}
				}
			}

			fineList := fineList(fineWithPlayers, pageId, finemasterPage)
			fineList.Render(r.Context(), w)
		}
			case "POST": {

				if err := r.ParseForm(); err != nil {
					http.Error(w, "Invalid form data", http.StatusBadRequest)
					return
				}

				playerIdStr := r.FormValue("playerId")
				playerId, err := strconv.ParseUint(playerIdStr, 10, 64)
				if err != nil || playerId == 0 {
					http.Error(w, fmt.Sprintf("Error parsing playerId %v", err), http.StatusBadRequest)
					return 
				}
				
				var presetFine *PresetFine
				presetFineStr := r.FormValue("presetFineId")
				if(len(presetFineStr) > 0 && presetFineStr != "-1"){
					pfId, err := strconv.ParseUint(r.FormValue("presetFineId"), 10, 64)
					if err != nil {
						http.Error(w, "Invalid presetFineId", http.StatusBadRequest)
						return
					}
					presetFine, err = GetPresetFine(db, pfId)
					if err != nil {
						http.Error(w, "Could not GetPresetFine", http.StatusBadRequest)
						return
					}


				}

				var fine Fine
				if presetFine != nil {
					fine = Fine{
						Amount: presetFine.Amount,
						Reason: presetFine.Reason,
					}
				} else {
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


					reason := r.FormValue("reason")

					if r.FormValue("should_be_a_rule") == "on" {
						var suggestedPFine = &PresetFine{
							Reason: reason,
							Approved: false,
							Amount: amount,
						}
						SavePresetFine(db, suggestedPFine)
					}

					fine = Fine{
						Amount: amount,
						Reason: reason,
					}
				}

				fine.PlayerID = uint(playerId)
				fine.Approved = r.FormValue("approved") == "on"
				
				if err := SaveFine(db, &fine); err != nil {
					log.Printf("Error saving fine: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))

				w.WriteHeader(http.StatusOK)
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
			w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
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
			w.Header().Set("HX-Redirect", r.Header.Get("Referrer"))
			w.WriteHeader(http.StatusOK)
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
    FinesOpen bool `schema:"f"` // Assuming `pfid` is the query param name
	PlayerOpen bool `schema:"p"` // Assuming `pfid` is the query param name
	PresetFinesOpen bool `schema:"pf"` // Assuming `pfid` is the query param name
}

func presetFineHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
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
	FinesOpen bool `schema:"f"` // Assuming `pfid` is the query param name
	PlayerOpen bool `schema:"p"` // Assuming `pfid` is the query param name
	PresetFinesOpen bool `schema:"pf"` // Assuming `pfid` is the query param name
	// Pass *string `schema:"pass"`
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

		playersWithFines, err := FetchPlayersWithFines(db)
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

		finemaster := finemaster(pass, playersWithFines, pFines, *queryParams)
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
	r.HandleFunc("/fines", fineHandler(db))
	r.HandleFunc("/fines/approve", fineApproveHandler(db))
	r.HandleFunc("/preset-fines", presetFineHandler(db))
	r.HandleFunc("/preset-fines/approve", presetFineApproveHandler(db))
	r.HandleFunc("/finemaster/{pass}", presetFineMasterHandler(db))
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := schema.NewDecoder()
		queryParams := new(HomeQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("Error decoding query params: %v", err)
			http.Error(w, "Bad Request - home Decode", http.StatusBadRequest)
			return
		}

		playersWithFines, err := FetchPlayersWithFines(db)
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

    // Start the HTTP server.
    http.ListenAndServe(":8080", r)
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