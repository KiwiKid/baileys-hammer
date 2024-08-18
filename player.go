package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func playerNamesHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			pwfs, err := GetPlayers(db, 0, 999)
			if err != nil {
				http.Error(w, "Could not get matches", http.StatusNotFound)
				return
			}
			matchComp := playerNames(pwfs)
			matchComp.Render(GetContext(r), w)
		}
	}
}

type PlayerPaymentsWithTotals struct {
	PlayerPayments []PlayerPayment
	Total          float64
	PlayerName     string
}

func playerPayments(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			log.Printf("PlayerPayments POST entry")
			//err := warning("hmmm")
			//err.Render(GetContext(r), w)

			if err := r.ParseForm(); err != nil {
				http.Error(w, "P - Invalid form data", http.StatusBadRequest)
				return
			}
			playerIDs := r.Form["players[]"]
			//if err != nil {
			//	http.Error(w, "Could not parse playerIDs", http.StatusBadRequest)
			//	return
			//}
			amountStr := r.FormValue("amount")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				http.Error(w, "Could not parse amount", http.StatusBadRequest)
				return
			}
			seasonIdParam := chi.URLParam(r, "seasonId")
			seasonId, err := strconv.ParseUint(seasonIdParam, 10, 64)
			if err != nil {
				http.Error(w, "Could not parse seasonId", http.StatusBadRequest)
				return
			}
			var totalAdded float64
			var totalRecords int
			for _, playerIDStr := range playerIDs {
				playerID, err := strconv.ParseUint(playerIDStr, 10, 64)

				err = CreatePlayerPayment(db, uint(playerID), amount, uint(seasonId))
				if err != nil {
					http.Error(w, "Could not add payment", http.StatusBadRequest)
					return
				} else {
					totalAdded += amount
					totalRecords++
				}
			}
			success := success(fmt.Sprintf("%d Payment  added for $%v and %d player(s) ", totalRecords, totalAdded, len(playerIDs)))
			success.Render(GetContext(r), w)
			return

		case "GET":
			activeSeason, err := GetActiveSeason(db)
			if err != nil {
				http.Error(w, "Failed when getting active season", http.StatusNotFound)
				return
			}

			displayType := r.URL.Query().Get("displayType")
			if displayType == "button" {

				btn := playerPaymentsButton("Open Player Payments", "table", activeSeason)
				btn.Render(GetContext(r), w)
				return
			}

			seasonIdParam := chi.URLParam(r, "seasonId")
			seasonId, err := strconv.ParseUint(seasonIdParam, 10, 64)

			pays, err := GetPlayerPayments(db, uint(seasonId))
			if err != nil {
				http.Error(w, "Could not get matches", http.StatusNotFound)
				return
			}
			players, err := GetPlayers(db, 0, 999)

			// Create a map of player id to player name with total via PlayerPaymentsWithTotals
			playerTotals := []PlayerPaymentsWithTotals{}
			for _, player := range players {
				total := 0.0
				playerPayments := []PlayerPayment{}
				for _, pay := range pays {
					if pay.PlayerID == player.ID {
						playerPayments = append(playerPayments, pay)
						total += pay.Amount
					}
				}
				//	if total > 0 {
				playerTotals = append(playerTotals, PlayerPaymentsWithTotals{
					PlayerPayments: playerPayments,
					Total:          total,
					PlayerName:     player.Name,
				})
				//	}
			}
			if displayType == "" {
				displayType = "table"
			}
			switch displayType {
			case "table":
				paymentsComp := viewPlayerPayments(playerTotals, players, activeSeason)
				paymentsComp.Render(GetContext(r), w)
				return
			default:
				http.Error(w, "Invalid display type", http.StatusBadRequest)
				return
			}
		}
		http.Error(w, "Invalid http method  type", http.StatusBadRequest)
		return

	}
}
