package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

var EVENT_TYPES = []string{"subbed-off", "subbed-on"}

type TimeOpt struct {
	Name string
	Value string
}

var TIME_OPTS = []TimeOpt{
    TimeOpt{
        Name: "Now",
        Value: "now",
    },
    TimeOpt{
        Name: "+1 min",
        Value: "+1",
    },
    TimeOpt{
        Name: "+1 min",
        Value: "+2",
    },
    TimeOpt{
        Name: "+1 min",
        Value: "+3",
    },
}

type MatchMeta struct {
    TimeOpts []TimeOpt
    EventTypes []string
}

var meta = &MatchMeta{
    TimeOpts: TIME_OPTS,
    EventTypes: EVENT_TYPES,
}

func matchEventHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request){


	return func(w http.ResponseWriter, r *http.Request) {
		matchIdStr := chi.URLParam(r, "matchId")
		eventIdStr := chi.URLParam(r, "eventId")
		var isOpen bool = false

		fmt.Printf("matchHandler %s", matchIdStr)
		matchId, err := strconv.ParseUint(matchIdStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			isOpen = r.URL.Query().Get("isOpen") == "true"

			if eventIdStr == "" {
				matchComp := addMatchEvent(*meta, matchId, isOpen)
				matchComp.Render(r.Context(), w)
				return
			} else {
                match, err := GetMatchEvent(db, matchId)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusBadRequest)
                    return
                }

				editEventComp := editMatchEvent(*meta, *match,  isOpen, matchId)
				editEventComp.Render(r.Context(), w)
				return
			}
		case "POST":
			form, err := parseForm(r)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            if matchIdStr == "" {
                // Create new match
                handleCreateMatchEvent(db, form, w)
            } else {
                // Edit existing match
                handleEditMatchEvent(db, form, matchIdStr, w)
            }
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

type NewMatchEventForm struct {
    Location  string
    StartTime string // Using string here for simplicity; parsing is needed
    Opponent  string
    Subtitle  string
}


func parseForm(r *http.Request) (NewMatchEventForm, error) {
    if err := r.ParseForm(); err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid form data")
    }

    startTime, err := time.Parse("2006-01-02T15:04", r.FormValue("startTime"))
    if err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid start time format: %s", r.FormValue("startTime"))
    }

    return NewMatchEventForm{
        Location:  r.FormValue("location"),
        StartTime: startTime.Format("2006-01-02T15:04"),
        Opponent:  r.FormValue("opponent"),
        Subtitle:  r.FormValue("subtitle"),
    }, nil
}

func handleCreateMatchEvent(db *gorm.DB, form NewMatchEventForm, w http.ResponseWriter) {
    // Convert startTime back to time.Time for saving
    startTime, _ := time.Parse("2006-01-02T15:04", form.StartTime)
    match := Match{
        Location:  form.Location,
        StartTime: &startTime,
        Opponent:  form.Opponent,
        Subtitle:  form.Subtitle,
        // Assume seasonId is defined elsewhere
        SeasonId: seasonId,
    }

    // Save the match
    if err := db.Create(&match).Error; err != nil {
        http.Error(w, fmt.Sprintf("Error saving match: %v", err), http.StatusInternalServerError)
        return
    }

    // Redirect to the new match
    w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", match.ID))
}


func handleEditMatchEvent(db *gorm.DB, form NewMatchEventForm, matchIdStr string, w http.ResponseWriter) {
    matchId, err := strconv.ParseUint(matchIdStr, 10, 64)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error parsing match ID: %v", err), http.StatusBadRequest)
        return
    }
/*
    // Convert startTime back to time.Time for saving
    startTime, _ := time.Parse("2006-01-02T15:04", form.StartTime)

    // Update match
    if err := db.Model(&MatchEvent{}).Where("id = ?", matchId).Updates(MatchEvent{
		Location:  form.Location,
		StartTime: &startTime,
		Opponent:  form.Opponent,
		Subtitle:  form.Subtitle,
		// Assume SeasonId is appropriately handled elsewhere
	}).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error updating match: %v", err), http.StatusInternalServerError)
		return
	}
	*/
	// Redirect to the updated match
	w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d/not/done-yet", matchId))
}