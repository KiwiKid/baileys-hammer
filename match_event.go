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

var EVENT_TYPES = []string{"subbed-off", "subbed-on"}

type TimeOpt struct {
	Name string
	Value string
}

var TIME_OPTS = []TimeOpt{
    {
        Name: "Now",
        Value: "now",
    },
    {
        Name: "+1 min",
        Value: "1",
    },
    {
        Name: "+1 min",
        Value: "2",
    },
    {
        Name: "+1 min",
        Value: "3",
    },
}

type MatchMeta struct {
    CurrentMatchMinute int
    EventTypes []string
    Players []Player
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

        players, err := FetchActivePlayers(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error FetchActivePlayers: %v", err), http.StatusBadRequest)
			return
		}else{
            log.Printf("GOT %d active players", len(players))
        }
        
        match, err := GetMatch(db, matchId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        startTime := time.Date(2024, 3, 28, 9, 0, 0, 0, time.UTC)
        now := time.Now()
        duration := now.Sub(startTime)

        var meta = &MatchMeta{
            CurrentMatchMinute: int(duration.Minutes()),
            EventTypes: EVENT_TYPES,
            Players: players,
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



            log.Println("matchEventHandler")
			form, err := parseForm(r, *match)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            if form.ID == 0 {
                // Create new event
                handleCreateMatchEvent(db, form, w)
            } else {
                // Edit existing event
                handleEditMatchEvent(db, form, *match, w)
            }

            log.Println("matchEventHandlerend")

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

        http.Error(w, "(fell out)", http.StatusMethodNotAllowed)

	}
}

type NewMatchEventForm struct {
    gorm.Model
    MatchId  uint64
    EventName string
    EventType string // 'subbed-off' / 'subbed-on' / 'goal' / 'assist' / 'own-goal'
    EventMinute int
    EventTime *time.Time `json:"timestamp" gorm:"type:datetime"`
    EventTimeOffset  int
}



func parseForm(r *http.Request, match Match) (NewMatchEventForm, error) {
    if err := r.ParseForm(); err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid form data")
    }

   /* startTime, err := time.Parse("2006-01-02T15:04", r.FormValue("eventTime"))
    if err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid start time format: %s", r.FormValue("eventTime"))
    }*/

    var eventTimeOffsetStr = r.FormValue("eventTimeOffset")
    var eventTimeStr = r.FormValue("eventTime")

    var eventTime time.Time
    if(eventTimeOffsetStr == "now"){
        if eventTimeStr != "" {
            newStartTime, startErr := time.Parse("2006-01-02T15:04", eventTimeStr)
            if startErr != nil {
                return NewMatchEventForm{}, startErr
            }
            eventTime = newStartTime
        } else {
            eventTime = time.Now()
        }
    } else {
        eventTimeOffset, err := strconv.ParseUint(eventTimeOffsetStr, 10, 64)
        if err != nil {
            return NewMatchEventForm{}, err
        }
        eventTime = eventTime.Add(-time.Minute * time.Duration(eventTimeOffset))
    }

    matchIdStr := r.FormValue("matchId")
    matchId, err := strconv.ParseUint(matchIdStr, 10, 64) // Convert string to int
    if err != nil {
        // Handle the error if the conversion fails
        return NewMatchEventForm{}, fmt.Errorf("invalid MatchId: %s %+v",matchIdStr,  err)
    }

    eventMinStr := r.FormValue("eventMinute")
    eventMin, err := strconv.Atoi(eventMinStr) // Convert string to int
    if err != nil {
        // Handle the error if the conversion fails
        return NewMatchEventForm{}, fmt.Errorf("invalid eventMinute: %s %+v", eventMinStr,  err)
    }

    return NewMatchEventForm{
        EventName:  r.FormValue("location"),
        EventType: r.FormValue("eventType"),
        MatchId: matchId,
        EventTime: &eventTime,
        EventMinute: eventMin,
    }, nil
}

func handleCreateMatchEvent(db *gorm.DB, form NewMatchEventForm, w http.ResponseWriter) {
    log.Println("handleCreateMatchEvent")
    // Convert startTime back to time.Time for saving
    matchEvt := MatchEvent{
        EventName: form.EventName,
        EventType: form.EventType,
        MatchId: form.MatchId,
        EventTime: form.EventTime,

    }

    // Save the match
    if err := db.Create(&matchEvt).Error; err != nil {
        http.Error(w, fmt.Sprintf("Error saving match: %v", err), http.StatusInternalServerError)
        return
    }else{
        log.Printf("Created Match event ✨ \n%+v", matchEvt)
    }

    // Redirect to the new match
    w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchEvt.MatchId))
    w.WriteHeader(http.StatusOK);
}


func handleEditMatchEvent(db *gorm.DB, form NewMatchEventForm, match Match, w http.ResponseWriter) {
    log.Println("handleEditMatchEvent")
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
	w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", match.ID))
    w.WriteHeader(http.StatusOK)
}

func matchEventListHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
        switch(r.Method){
        case "GET": {
            matchIdStr := chi.URLParam(r, "matchId")

            log.Printf("matchEventListHandler - %s", matchIdStr)

			var matchId uint64
            var err error
			if len(matchIdStr) == 0{
				http.Error(w, "Error parsing matchId", http.StatusBadRequest)
					return 
			} else {
				matchId, err = strconv.ParseUint(matchIdStr, 10, 64)
				if err != nil || matchId == 0 {
					http.Error(w, fmt.Sprintf("Error parsing limitStr %v", err), http.StatusBadRequest)
					return 
				}
			}

            log.Printf("matchEventListHandler - GetMatchEvents(%d)", matchId)


			matchEvents, getFErr := GetMatchEvents(db, matchId)
			if getFErr != nil {
				http.Error(w, "Invalid GetMatchEvents data", http.StatusBadRequest)
				return
			}

            players, err := GetPlayers(db, 1, 9999)
			if err != nil {
				http.Error(w, "Invalid GetPlayers", http.StatusBadRequest)
				return
			}

            initialState := MatchState{PlayersOn: []PlayerState{}, ScoreFor: 0, ScoreAgainst: 0}
            matchState := ConstructMatchState(matchEvents, initialState, 90, players)

			fineList := listMatchEvents(matchState, matchEvents)
			fineList.Render(r.Context(), w)
		}
        default: 
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }
    }
}