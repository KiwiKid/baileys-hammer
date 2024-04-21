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

var EVENT_TYPES = []string{"subbed-off", "subbed-on", "goal", "assist", "own-goal", "opponent-goal", "attended-training", "injury", "pregame-injury"}

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

type MatchMetaGeneral struct {
    Match Match
    Players []Player
    PlayerOfTheDay *Player
    DudOfTheDay *Player
}


func matchEventHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request){

	return func(w http.ResponseWriter, r *http.Request) {
		matchIdStr := chi.URLParam(r, "matchId")
		eventIdStr := chi.URLParam(r, "eventId")
		var isOpen bool = false

		matchId, err := strconv.ParseUint(matchIdStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("matchEventHandler - Error parsing match ID: %v", err), http.StatusBadRequest)
			return
		}

        players, err := FetchActivePlayers(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("matchEventHandler - Error FetchActivePlayers: %v", err), http.StatusBadRequest)
			return
		} else{
            log.Printf("GOT %d active players", len(players))
        }
        
        match, err := GetMatch(db, matchId)
        if err != nil {
			http.Error(w, fmt.Sprintf("matchEventHandler - Error GetMatch: %v", err), http.StatusBadRequest)
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

				matchComp := createNewEvent(matchId)
				matchComp.Render(GetContext(r), w)
				return
			} else {
                match, err := GetMatchEvent(db, matchId)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusBadRequest)
                    return
                }

				editEventComp := editMatchEvent(*meta, *match,  isOpen, matchId)
				editEventComp.Render(GetContext(r), w)
				return
			}
		case "POST":
            log.Printf("matchEventHandler %+v", r)
			form, err := parseForm(r, *match)
            if err != nil {
                var errMessage = fmt.Sprintf("matchEventHandler - POST parseForm failed - %v \n\n%+v", match, err)
                errComp := errMsg(errMessage)
				errComp.Render(GetContext(r), w)
                return
            }

            if form.ID == 0 {
                // Create new event
                handleCreateMatchEvent(db, form, w)
            } else {
                // Edit existing event
                handleEditMatchEvent(db, form, *match, w)
            }

            matchState, matchEvents, getMatchErr := GetMatchAndEvents(db, matchId)
			if getMatchErr != nil {
				errComp := errMsg(fmt.Sprintf("Invalid amount %v", getMatchErr))
				errComp.Render(GetContext(r), w)
			}

			fineList := listMatchEvents(*matchState, matchEvents)
			fineList.Render(GetContext(r), w)
        case "DELETE":
            if eventIdStr == "" {
                errComp := errMsg("Invalid eventId param")
				errComp.Render(GetContext(r), w)
            }
            eventId, err := strconv.ParseUint(eventIdStr, 10, 64)
            if err != nil {
                errComp := errMsg(fmt.Sprintf("matchEventHandler - Error parsing match ID: %v", err))
				errComp.Render(GetContext(r), w)
                return
            }
            delErr := DeleteMatchEvent(db, uint(eventId))
            if delErr != nil {
                errComp := errMsg(fmt.Sprintf("matchEventHandler - Error DeleteMatchEvent: %v", delErr))
				errComp.Render(GetContext(r), w)
                return
            }

            successComp := success(F("deleted event %d", eventId))
            successComp.Render(GetContext(r), w)
            return
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
    PlayerId uint64
    EventTimeOffset  int
}



func parseForm(r *http.Request, match Match) (NewMatchEventForm, error) {
    log.Printf("parseForm")
    if err := r.ParseForm(); err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid form data - %+v", err)
    }
    log.Printf("parseForm1")

   /* startTime, err := time.Parse("2006-01-02T15:04", r.FormValue("eventTime"))
    if err != nil {
        return NewMatchEventForm{}, fmt.Errorf("invalid start time format: %s", r.FormValue("eventTime"))
    }*/

    var eventTimeOffsetStr = r.FormValue("eventTimeOffset")
    var eventTimeStr = r.FormValue("eventTime")
    log.Printf("parseForm2")

    var eventTime time.Time
    if(eventTimeOffsetStr == "now"){
        if eventTimeStr != "" {
            newStartTime, startErr := time.Parse("2006-01-02T15:04", eventTimeStr)
            if startErr != nil {
                return NewMatchEventForm{}, fmt.Errorf("invalid eventTimeStr data - %s- %+v", eventTimeStr, startErr)
            }
            eventTime = newStartTime
        } else {
            eventTime = time.Now()
        }
    } else {
        eventTimeOffset, err := strconv.ParseUint(eventTimeOffsetStr, 10, 64)
        if err != nil {
            return NewMatchEventForm{}, fmt.Errorf("invalid eventTimeOffsetStr data - %s %+v", eventTimeOffsetStr, err)
        }
        eventTime = eventTime.Add(-time.Minute * time.Duration(eventTimeOffset))
    }
    log.Printf("parseForm3")

    matchIdStr := r.FormValue("matchId")
    matchId, err := strconv.ParseUint(matchIdStr, 10, 64) // Convert string to int
    if err != nil {
        // Handle the error if the conversion fails
        return NewMatchEventForm{}, fmt.Errorf("invalid MatchId: %s %+v",matchIdStr,  err)
    }
    log.Printf("parseForm4")
     /*eventMinStr := r.FormValue("eventMinute")
  eventMin, err := strconv.Atoi(eventMinStr) // Convert string to int
    if err != nil {
        // Handle the error if the conversion fails
        return NewMatchEventForm{}, fmt.Errorf("invalid eventMinute: %s %+v", eventMinStr,  err)
    }*/

    return NewMatchEventForm{
        EventName:  r.FormValue("eventName"),
        EventType: r.FormValue("eventType"),
        MatchId: matchId,
        EventTime: &eventTime,
        EventMinute: 0,
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
        PlayerId: uint(form.PlayerId),
    }

    // Save the match
    if err := db.Create(&matchEvt).Error; err != nil {
        http.Error(w, fmt.Sprintf("Error saving match: %v", err), http.StatusInternalServerError)
        return
    }else{
        log.Printf("Created Match event ✨ \n%+v", matchEvt)
    }

    
    // Redirect to the new match
   // w.Header().Set("HX-Redirect", fmt.Sprintf("/match/%d", matchEvt.MatchId))
   // w.WriteHeader(http.StatusOK);
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
                errComp := errMsg(fmt.Sprintf("Error parsing matchIdStr %v", matchIdStr))
                errComp.Render(GetContext(r), w)
				return 
			} else {
				matchId, err = strconv.ParseUint(matchIdStr, 10, 64)
				if err != nil || matchId == 0 {
                    errComp := errMsg(fmt.Sprintf("Error parsing matchIdStr %v", matchIdStr))
				    errComp.Render(GetContext(r), w)
                    return 
				}
			}

            log.Printf("matchEventListHandler - GetMatchEvents(%d)", matchId)


			matchState, matchEvents, getMatchErr := GetMatchAndEvents(db, matchId)
			if getMatchErr != nil {
				errComp := errMsg(fmt.Sprintf("Invalid amount %v", getMatchErr))
				errComp.Render(GetContext(r), w)
			}

			fineList := listMatchEvents(*matchState, matchEvents)
			fineList.Render(GetContext(r), w)
		}
        default: 
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }
    }
}