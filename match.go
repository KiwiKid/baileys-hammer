package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type MatchPageData struct {
	Match Match
	isOpen bool
}

type NewMatchForm struct {
    Location  string
    StartTime string // Using string here for simplicity; parsing is needed
    Opponent  string
    Subtitle  string
}


const seasonId = 2024

func matchListHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var isOpen bool = false

		switch r.Method {
			case "GET":
				isOpen = r.URL.Query().Get("isOpen") == "true"

				var season uint64
				seasonStr := r.URL.Query().Get("season")
				if len(seasonStr) == 0 {
					season = 0
				}else {
					seasonUint, err := strconv.ParseInt(seasonStr, 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
						return 
					}
					season = uint64(seasonUint)
				}

				var page = 0
				pageStr := r.URL.Query().Get("page")
				if len(pageStr) == 0 {
					page = 0
				}else {
					pageIdUint, err := strconv.ParseInt(pageStr, 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("Error parsing page %v", err), http.StatusBadRequest)
						return 
					}
					page = int(pageIdUint)
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

				log.Printf("GetMatches(season%+v, page%+v, limit: %+v)", season, page, limit)
			match, err := GetMatches(db, season, page, limit)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error retrieving match: %v", err), http.StatusInternalServerError)
				return
			}

			matchComp := matchListPage(match, isOpen)
			matchComp.Render(r.Context(), w)
		
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func matchHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("matchHandler")
	var matchId uint
	var isOpen bool = false
	var err error
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			isOpen = r.URL.Query().Get("isOpen") == "true"

			matchIdStr := chi.URLParam(r, "matchId")
			fmt.Printf("matchHandler %s", matchIdStr)

			if matchIdStr == "" {
				matchComp := createMatch(isOpen)
				matchComp.Render(r.Context(), w)
				return
			}
			var matchId64 uint64
			matchId64, err = strconv.ParseUint(matchIdStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
				return
			}

			matchId = uint(matchId64)

			match, err := GetMatchWithEvents(db, matchId)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error retrieving match: %v", err), http.StatusInternalServerError)
				return
			}
	
			// Prepare the data for the template
			matchData := MatchPageData{
				Match:  *match,
				isOpen: isOpen,
			}
	
			matchComp := matchPage(matchData)
			matchComp.Render(r.Context(), w)
			return
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
				SeasonId:  seasonId,
			}

			matchId, err = SaveMatch(db, &match)
			if err != nil {
				http.Error(w, "Invalid start time format", http.StatusBadRequest)
				return
			}

			w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchId))
			return
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}



// Mock function for rendering templates. Replace with your actual implementation.
func RenderTemplate(w http.ResponseWriter, r *http.Request, data interface{}) {
	// Your template rendering logic here
	fmt.Fprintf(w, "Template rendering with data: %+v\n", data) // Placeholder implementation
}


func ConstructMatchState(events []MatchEvent, currentState MatchState, currentTime int, allPlayers []Player) MatchState {
    if len(events) == 0 {
        return currentState
    }

    event := events[0]
	var player Player 
	if event.PlayerId > 0 { 

	
		for _, p := range allPlayers {
			if(p.ID == uint(event.PlayerId)){
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
    }

    // Update players' time played if they are on the field
    for i, player := range currentState.PlayersOn {
        if uint64(player.PlayerId) == event.PlayerId {
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