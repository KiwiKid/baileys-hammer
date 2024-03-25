package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MatchPageData struct {
	Match Match 
}

type NewMatchForm struct {
    Location  string
    StartTime string // Using string here for simplicity; parsing is needed
    Opponent  string
    Subtitle  string
}

func matchHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			matchIdStr := r.URL.Query().Get("match_id")
			if matchIdStr == "" {
				matchComp := createMatch()
				matchComp.Render(r.Context(), w)
				return
			}

			matchId, err := strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
				return
			}

			match, err := GetMatchWithEvents(db, matchId)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error retrieving match: %v", err), http.StatusInternalServerError)
				return
			}

			// Prepare the data for the template
			matchData := MatchPageData{
				Match:  *match,
			}

			matchComp := matchPage(matchData)
			matchComp.Render(r.Context(), w)

		case "POST":
			// Simplified example: Assume the response after a POST is a success message or error
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
	
			form := NewMatchForm{
				Location:  r.FormValue("location"),
				StartTime: r.FormValue("startTime"),
				Opponent:  r.FormValue("opponent"),
				Subtitle:  r.FormValue("subtitle"),
			}
	
			// Parse start time
			startTime, err := time.Parse("2006-01-02T15:04", form.StartTime)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid start time format %s", form.StartTime), http.StatusBadRequest)
				return
			}
	
			// Create a new match based on the form data
			match := Match{
				Location:  form.Location,
				StartTime: &startTime,
				Opponent:  form.Opponent,
				Subtitle:  form.Subtitle,
			}

			id, err := SaveMatch(db, &match)
			if err != nil {
				http.Error(w, "Invalid start time format", http.StatusBadRequest)
				return
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", id))

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Mock function for rendering templates. Replace with your actual implementation.
func RenderTemplate(w http.ResponseWriter, r *http.Request, data interface{}) {
	// Your template rendering logic here
	fmt.Fprintf(w, "Template rendering with data: %+v\n", data) // Placeholder implementation
}
