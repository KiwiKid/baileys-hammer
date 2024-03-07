package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

var decoder = schema.NewDecoder()

func savePlayerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch(r.Method){
			case "POST": {
				if err := r.ParseForm(); err != nil {
					log.Printf("Error parsing form data: %v", err)
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
	
				var player Player
				// Use the decoder to populate the player struct
				if err := decoder.Decode(&player, r.PostForm); err != nil {
					log.Printf("Error decoding form into player struct: %v", err)
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
	
				// Now, `player` is populated with values from the form
				// Save the player using your existing logic
				if err := SavePlayer(db, &player); err != nil {
					log.Printf("Error saving player: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
	
			}
		}

	}
}

func saveFineHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch(r.Method){
		case "POST": {
			var fine Fine
			if err := json.NewDecoder(r.Body).Decode(&fine); err != nil {
				log.Printf("Error decoding fine data: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			if err := SaveFine(db, &fine); err != nil {
				log.Printf("Error saving fine: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	}
}

func GenerateUrl(baseURL string, queryParams *HomeQueryParams) (*string, error, ) {
	queryParams.IsFineMaster = true
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
	IsFineMaster bool `schema:"fm"`
	OpenPlayers []uint `schema:"op"`
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
				Reason: r.FormValue("Reason"),
			}
	
			// Parse Amount as float64 from form value
			if amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64); err == nil {
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

            w.WriteHeader(http.StatusCreated)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }
    }
}

func main() {
	log.Printf("Start")
	db, err := DBInit()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	http.HandleFunc("/players", savePlayerHandler(db))
	http.HandleFunc("/fines", saveFineHandler(db))
	http.HandleFunc("/preset-fines", presetFineHandler(db))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := schema.NewDecoder()
		queryParams := new(HomeQueryParams)
		if err := decoder.Decode(queryParams, r.URL.Query()); err != nil {
			log.Printf("Error decoding query params: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}else{
			log.Printf("GOT QUERY:\n%+v", queryParams)
		}

		playersWithFines, err := FetchPlayersWithFines(db)
		if err != nil {
			log.Printf("Error fetching players with fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		pFines, err := GetPresetFines(db)
		if err != nil {
			log.Printf("Error retrieving preset fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		home := home(playersWithFines, pFines, *queryParams)
		home.Render(r.Context(), w)
	})

    // Start the HTTP server.
    http.ListenAndServe(":8080", nil)
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