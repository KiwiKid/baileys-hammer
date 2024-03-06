package main

import (
	"encoding/json"
	"log"
	"net/http"

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
			case "GET": {
				player := player()
				player.Render(r.Context(), w)
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
		case "GET": {
			fines := fines()
			fines.Render(r.Context(), w)
		}

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

	http.HandleFunc("/player", savePlayerHandler(db))
	http.HandleFunc("/fine", saveFineHandler(db))
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		playersWithFines, err := FetchPlayersWithFines(db)
		if err != nil {
			log.Printf("Error fetching players with fines: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("response")


		home := home(playersWithFines, true)
		home.Render(r.Context(), w)
	})

    // Start the HTTP server.
    http.ListenAndServe(":8080", nil)
	log.Printf("Listening on %d", 8080)
}
