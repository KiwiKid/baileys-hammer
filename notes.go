package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const (
	noteTargetMatch      = "match"
	noteTargetOpposition = "opposition"
	noteTargetField      = "field"
	noteTargetPlayer     = "player"
	noteTargetFormation  = "formation"
	noteTargetLineup     = "lineup"

	noteChannelNote     = "note"
	noteChannelFeedback = "feedback"

	noteTypeFocus = "focus"
)

var noteTargetKinds = []string{
	noteTargetMatch,
	noteTargetOpposition,
	noteTargetField,
	noteTargetPlayer,
	noteTargetFormation,
	noteTargetLineup,
}

type NoteTargetOption struct {
	Kind  string
	ID    uint
	Label string
}

type NotesPageData struct {
	Notes             []MLNote
	Feedback          []MLNote
	Targets           map[string][]NoteTargetOption
	Selected          MLNote
	Message           string
	ShowAdminNotes    bool
	IsEditing         bool
	ReturnTo          string
	LockedTargetKind  string
	LockedTargetID    uint
	LockedTargetLabel string
}

type FeedbackPageData struct {
	Matches   []NoteTargetOption
	Selected  MLNote
	Message   string
	IsError   bool
	ActionURL string
	LockMatch bool
}

func noteTargetKindLabel(kind string) string {
	switch kind {
	case noteTargetMatch:
		return "Match"
	case noteTargetOpposition:
		return "Opposition"
	case noteTargetField:
		return "Field"
	case noteTargetPlayer:
		return "Player"
	case noteTargetFormation:
		return "Formation"
	case noteTargetLineup:
		return "Line-up"
	case "general":
		return "General"
	default:
		return kind
	}
}

func noteTargetLabel(note MLNote) string {
	if strings.TrimSpace(note.TargetLabel) != "" {
		return note.TargetLabel
	}
	if note.TargetID > 0 {
		return fmt.Sprintf("#%d", note.TargetID)
	}
	return "General"
}

func isValidNoteTargetKind(kind string) bool {
	for _, candidate := range noteTargetKinds {
		if kind == candidate {
			return true
		}
	}
	return false
}

func noteTargets(db *gorm.DB, teamID uint) (map[string][]NoteTargetOption, error) {
	targets := map[string][]NoteTargetOption{}

	matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
	if err != nil {
		return nil, err
	}
	opponents := map[string]bool{}
	fields := map[string]bool{}
	for _, match := range matches {
		labelParts := []string{}
		if strings.TrimSpace(match.Opponent) != "" {
			labelParts = append(labelParts, "vs "+strings.TrimSpace(match.Opponent))
			opponents[strings.TrimSpace(match.Opponent)] = true
		}
		if strings.TrimSpace(match.Location) != "" {
			labelParts = append(labelParts, strings.TrimSpace(match.Location))
			fields[strings.TrimSpace(match.Location)] = true
		}
		if match.StartTime != nil {
			labelParts = append(labelParts, match.StartTime.Format("2 Jan 2006"))
		}
		if len(labelParts) == 0 {
			labelParts = append(labelParts, fmt.Sprintf("Match #%d", match.ID))
		}
		targets[noteTargetMatch] = append(targets[noteTargetMatch], NoteTargetOption{Kind: noteTargetMatch, ID: match.ID, Label: strings.Join(labelParts, " - ")})
	}

	for opponent := range opponents {
		targets[noteTargetOpposition] = append(targets[noteTargetOpposition], NoteTargetOption{Kind: noteTargetOpposition, Label: opponent})
	}
	sort.Slice(targets[noteTargetOpposition], func(i, j int) bool {
		return targets[noteTargetOpposition][i].Label < targets[noteTargetOpposition][j].Label
	})

	for field := range fields {
		targets[noteTargetField] = append(targets[noteTargetField], NoteTargetOption{Kind: noteTargetField, Label: field})
	}
	sort.Slice(targets[noteTargetField], func(i, j int) bool {
		return targets[noteTargetField][i].Label < targets[noteTargetField][j].Label
	})

	players, err := GetPlayers(db, 0, 999)
	if err != nil {
		return nil, err
	}
	sort.Slice(players, func(i, j int) bool { return players[i].Name < players[j].Name })
	for _, player := range players {
		targets[noteTargetPlayer] = append(targets[noteTargetPlayer], NoteTargetOption{Kind: noteTargetPlayer, ID: player.ID, Label: player.Name})
	}

	actor := LineupActor{TeamID: teamID, IsAdmin: true}
	formations, err := visibleFormations(db, actor)
	if err != nil {
		return nil, err
	}
	sort.Slice(formations, func(i, j int) bool { return formations[i].Name < formations[j].Name })
	for _, formation := range formations {
		targets[noteTargetFormation] = append(targets[noteTargetFormation], NoteTargetOption{Kind: noteTargetFormation, ID: formation.ID, Label: formation.Name})
	}

	lineups, err := visibleLineups(db, actor)
	if err != nil {
		return nil, err
	}
	sort.Slice(lineups, func(i, j int) bool { return lineups[i].Name < lineups[j].Name })
	for _, lineup := range lineups {
		targets[noteTargetLineup] = append(targets[noteTargetLineup], NoteTargetOption{Kind: noteTargetLineup, ID: lineup.ID, Label: lineup.Name})
	}

	return targets, nil
}

func getNotes(db *gorm.DB, teamID uint) ([]MLNote, error) {
	var notes []MLNote
	q := db.Order("priority DESC, updated_at DESC").Where("channel = ? OR channel = '' OR channel IS NULL", noteChannelNote)
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if err := q.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func getNotesForTarget(db *gorm.DB, teamID uint, targetKind string, targetID uint, targetLabel string) ([]MLNote, error) {
	var notes []MLNote
	q := db.Order("priority DESC, updated_at DESC").Where("(channel = ? OR channel = '' OR channel IS NULL) AND target_kind = ?", noteChannelNote, targetKind)
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if targetID > 0 {
		q = q.Where("target_id = ?", targetID)
	} else {
		q = q.Where("target_label = ?", targetLabel)
	}
	if err := q.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func getNotesForTargetByType(db *gorm.DB, teamID uint, targetKind string, targetID uint, targetLabel string, noteType string) ([]MLNote, error) {
	var notes []MLNote
	q := db.Order("priority DESC, updated_at DESC").
		Where("(channel = ? OR channel = '' OR channel IS NULL) AND target_kind = ? AND type = ?", noteChannelNote, targetKind, noteType)
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if targetID > 0 {
		q = q.Where("target_id = ?", targetID)
	} else {
		q = q.Where("target_label = ?", targetLabel)
	}
	if err := q.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func getFeedback(db *gorm.DB, teamID uint) ([]MLNote, error) {
	var feedback []MLNote
	q := db.Order("updated_at DESC").Where("channel = ?", noteChannelFeedback)
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if err := q.Find(&feedback).Error; err != nil {
		return nil, err
	}
	return feedback, nil
}

func noteReturnTo(r *http.Request) string {
	returnTo := strings.TrimSpace(r.URL.Query().Get("returnTo"))
	if strings.HasPrefix(returnTo, "/") && !strings.HasPrefix(returnTo, "//") {
		return returnTo
	}
	return ""
}

func noteActionURL(path string, returnTo string) string {
	if strings.TrimSpace(returnTo) == "" {
		return path
	}
	return path + "?returnTo=" + url.QueryEscape(returnTo)
}

func noteContextURL(kind string, targetID uint, targetLabel string, returnTo string) string {
	values := url.Values{}
	values.Set("targetKind", kind)
	if targetID > 0 {
		values.Set("targetId", fmt.Sprintf("%d", targetID))
	}
	if strings.TrimSpace(targetLabel) != "" {
		values.Set("targetLabel", targetLabel)
	}
	if strings.TrimSpace(returnTo) != "" {
		values.Set("returnTo", returnTo)
	}
	return "/notes/context?" + values.Encode()
}

func contextualNotesData(db *gorm.DB, teamID uint, kind string, targetID uint, targetLabel string, returnTo string) (NotesPageData, error) {
	notes, err := getNotesForTarget(db, teamID, kind, targetID, targetLabel)
	if err != nil {
		return NotesPageData{}, err
	}
	return NotesPageData{
		Notes:             notes,
		Selected:          MLNote{TeamID: teamID, Channel: noteChannelNote, TargetKind: kind, TargetID: targetID, TargetLabel: targetLabel, Priority: 3},
		ReturnTo:          returnTo,
		LockedTargetKind:  kind,
		LockedTargetID:    targetID,
		LockedTargetLabel: targetLabel,
	}, nil
}

func matchDayFocusPointsData(db *gorm.DB, teamID uint, match Match, returnTo string, showAdminNotes bool) (NotesPageData, error) {
	targetLabel := "Match"
	if strings.TrimSpace(match.Opponent) != "" {
		targetLabel = "vs " + strings.TrimSpace(match.Opponent)
	}
	notes, err := getNotesForTargetByType(db, teamID, noteTargetMatch, match.ID, targetLabel, noteTypeFocus)
	if err != nil {
		return NotesPageData{}, err
	}
	return NotesPageData{
		Notes:             notes,
		Selected:          MLNote{TeamID: teamID, Channel: noteChannelNote, TargetKind: noteTargetMatch, TargetID: match.ID, TargetLabel: targetLabel, Priority: 3, Type: noteTypeFocus},
		ReturnTo:          returnTo,
		ShowAdminNotes:    showAdminNotes,
		LockedTargetKind:  noteTargetMatch,
		LockedTargetID:    match.ID,
		LockedTargetLabel: targetLabel,
	}, nil
}

func getNote(db *gorm.DB, noteID uint, teamID uint) (*MLNote, error) {
	var note MLNote
	q := db.Where("id = ? AND (channel = ? OR channel = '' OR channel IS NULL)", noteID, noteChannelNote)
	if teamID > 0 {
		q = q.Where("team_id = ?", teamID)
	}
	if err := q.First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func noteFromForm(r *http.Request, teamID uint) (MLNote, error) {
	priority, err := strconv.Atoi(strings.TrimSpace(r.FormValue("priority")))
	if err != nil {
		return MLNote{}, errors.New("priority must be a number")
	}
	if priority < 0 {
		priority = 0
	}
	kind := strings.TrimSpace(r.FormValue("targetKind"))
	if !isValidNoteTargetKind(kind) {
		return MLNote{}, errors.New("target type is required")
	}
	targetID := uint(0)
	targetLabel := strings.TrimSpace(r.FormValue(kind + "Label"))
	if kind == noteTargetMatch || kind == noteTargetPlayer || kind == noteTargetFormation || kind == noteTargetLineup {
		parsedID, err := parseUintFormValue(r, kind+"Id")
		if err != nil || parsedID == 0 {
			return MLNote{}, errors.New("target is required")
		}
		targetID = parsedID
		targetLabel = strings.TrimSpace(r.FormValue(kind + "Label_" + fmt.Sprintf("%d", targetID)))
	} else if targetLabel == "" {
		return MLNote{}, errors.New("target is required")
	}
	noteType := strings.TrimSpace(r.FormValue("type"))
	if noteType == "" {
		return MLNote{}, errors.New("type is required")
	}
	note := strings.TrimSpace(r.FormValue("note"))
	if note == "" {
		return MLNote{}, errors.New("note is required")
	}
	creator := strings.TrimSpace(r.FormValue("creator"))
	if creator == "" {
		creator = "admin"
	}
	return MLNote{
		TeamID:      teamID,
		Channel:     noteChannelNote,
		TargetKind:  kind,
		TargetID:    targetID,
		TargetLabel: targetLabel,
		Priority:    priority,
		Creator:     creator,
		Type:        noteType,
		Note:        note,
	}, nil
}

func feedbackMatches(db *gorm.DB, teamID uint) ([]NoteTargetOption, error) {
	seasonID := activeSeasonID(db)
	matches := []Match{}
	query := db.Where("season_id = ?", seasonID)
	if teamID > 0 {
		query = query.Where("team_id = ?", teamID)
	}
	if err := query.Order("start_time DESC").Find(&matches).Error; err != nil {
		return nil, err
	}
	options := []NoteTargetOption{}
	now := time.Now()
	for _, match := range matches {
		options = append(options, NoteTargetOption{Kind: noteTargetMatch, ID: match.ID, Label: feedbackMatchLabel(match, now)})
	}
	return options, nil
}

func feedbackMatchLabel(match Match, now time.Time) string {
	labelParts := []string{}
	if strings.TrimSpace(match.Opponent) != "" {
		labelParts = append(labelParts, "vs "+strings.TrimSpace(match.Opponent))
	}
	if strings.TrimSpace(match.Location) != "" {
		labelParts = append(labelParts, strings.TrimSpace(match.Location))
	}
	if match.StartTime != nil {
		dateLabel := match.StartTime.Format("2 Jan 2006")
		if match.StartTime.After(now) {
			dateLabel += " (upcoming)"
		}
		labelParts = append(labelParts, dateLabel)
	}
	if len(labelParts) == 0 {
		labelParts = append(labelParts, fmt.Sprintf("Match #%d", match.ID))
	}
	return strings.Join(labelParts, " - ")
}

func feedbackFromForm(db *gorm.DB, r *http.Request, teamID uint) (MLNote, error) {
	scope := strings.TrimSpace(r.FormValue("scope"))
	if scope == "" {
		scope = "general"
	}
	if scope != "general" && scope != noteTargetMatch {
		return MLNote{}, errors.New("feedback scope is required")
	}
	feedbackType := strings.TrimSpace(r.FormValue("type"))
	if feedbackType != "general" && feedbackType != "specific" {
		return MLNote{}, errors.New("feedback type is required")
	}
	note := strings.TrimSpace(r.FormValue("note"))
	if note == "" {
		return MLNote{}, errors.New("feedback is required")
	}
	creator := strings.TrimSpace(r.FormValue("creator"))
	if r.FormValue("anonymous") == "on" || creator == "" {
		creator = "Anonymous"
	}

	feedback := MLNote{
		TeamID:      teamID,
		Channel:     noteChannelFeedback,
		TargetKind:  "general",
		TargetLabel: "General",
		Priority:    0,
		Creator:     creator,
		Type:        feedbackType,
		Note:        note,
	}
	if scope == noteTargetMatch {
		feedback.TargetKind = noteTargetMatch
		matchID, err := parseUintFormValue(r, "matchId")
		if err != nil || matchID == 0 {
			return feedback, errors.New("match is required")
		}
		match, err := GetMatch(db, matchID)
		if err != nil {
			return feedback, errors.New("match not found")
		}
		if teamID > 0 && match.TeamID != teamID {
			return feedback, errors.New("match not found")
		}
		feedback.TargetKind = noteTargetMatch
		feedback.TargetID = match.ID
		feedback.TargetLabel = feedbackMatchLabel(*match, time.Now())
	}
	return feedback, nil
}

func validateNoteTarget(db *gorm.DB, note MLNote) error {
	switch note.TargetKind {
	case noteTargetMatch:
		_, err := GetMatch(db, note.TargetID)
		return err
	case noteTargetPlayer:
		_, err := GetPlayerByID(db, note.TargetID)
		return err
	case noteTargetFormation:
		formation, err := getFormation(db, note.TargetID)
		if err != nil {
			return err
		}
		if note.TeamID > 0 && formation.TeamID != note.TeamID {
			return errors.New("formation is not for this team")
		}
	case noteTargetLineup:
		lineup, err := getLineup(db, note.TargetID)
		if err != nil {
			return err
		}
		if note.TeamID > 0 && lineup.TeamID != note.TeamID {
			return errors.New("line-up is not for this team")
		}
	case noteTargetOpposition, noteTargetField:
		if strings.TrimSpace(note.TargetLabel) == "" {
			return errors.New("target is required")
		}
	default:
		return errors.New("invalid target type")
	}
	return nil
}

func notesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		teamID := actor.TeamID
		if teamID == 0 || !actor.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		renderPage := func(selected MLNote, isEditing bool, msg string) {
			notes, err := getNotes(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load notes: %v", err)).Render(GetContext(r, db), w)
				return
			}
			feedback, err := getFeedback(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load feedback: %v", err)).Render(GetContext(r, db), w)
				return
			}
			targets, err := noteTargets(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load note targets: %v", err)).Render(GetContext(r, db), w)
				return
			}
			notesPage(NotesPageData{Notes: notes, Feedback: feedback, Targets: targets, Selected: selected, Message: msg, ShowAdminNotes: true, IsEditing: isEditing}).Render(GetContext(r, db), w)
		}

		switch r.Method {
		case "GET":
			renderPage(MLNote{Channel: noteChannelNote, Priority: 3, TargetKind: noteTargetMatch}, false, "")
		case "POST":
			if err := r.ParseForm(); err != nil {
				renderPage(MLNote{Channel: noteChannelNote, Priority: 3, TargetKind: noteTargetMatch}, false, "Invalid form data")
				return
			}
			note, err := noteFromForm(r, teamID)
			if err != nil {
				renderPage(note, false, err.Error())
				return
			}
			if err := validateNoteTarget(db, note); err != nil {
				renderPage(note, false, err.Error())
				return
			}
			if err := db.Create(&note).Error; err != nil {
				renderPage(note, false, fmt.Sprintf("Could not save note: %v", err))
				return
			}
			if returnTo := noteReturnTo(r); returnTo != "" {
				http.Redirect(w, r, returnTo, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/notes", http.StatusSeeOther)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func notesContextHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		teamID := actor.TeamID
		if teamID == 0 || !actor.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		kind := strings.TrimSpace(r.URL.Query().Get("targetKind"))
		if !isValidNoteTargetKind(kind) {
			http.Error(w, "Invalid note target", http.StatusBadRequest)
			return
		}
		targetID := uint(0)
		if targetIDStr := strings.TrimSpace(r.URL.Query().Get("targetId")); targetIDStr != "" {
			parsed, err := strconv.ParseUint(targetIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid note target", http.StatusBadRequest)
				return
			}
			targetID = uint(parsed)
		}
		targetLabel := strings.TrimSpace(r.URL.Query().Get("targetLabel"))
		notesData, err := contextualNotesData(db, teamID, kind, targetID, targetLabel, noteReturnTo(r))
		if err != nil {
			warning(fmt.Sprintf("Could not load model notes: %v", err)).Render(GetContext(r, db), w)
			return
		}
		notesData.ShowAdminNotes = true
		contextualNotes(notesData).Render(GetContext(r, db), w)
	}
}

func noteDetailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actor := getLineupActor(r, db)
		teamID := actor.TeamID
		if teamID == 0 || !actor.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		noteID64, err := strconv.ParseUint(chi.URLParam(r, "noteId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid note", http.StatusBadRequest)
			return
		}
		noteID := uint(noteID64)
		renderPage := func(selected MLNote, msg string) {
			notes, err := getNotes(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load notes: %v", err)).Render(GetContext(r, db), w)
				return
			}
			feedback, err := getFeedback(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load feedback: %v", err)).Render(GetContext(r, db), w)
				return
			}
			targets, err := noteTargets(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load note targets: %v", err)).Render(GetContext(r, db), w)
				return
			}
			notesPage(NotesPageData{Notes: notes, Feedback: feedback, Targets: targets, Selected: selected, Message: msg, ShowAdminNotes: true, IsEditing: true}).Render(GetContext(r, db), w)
		}

		existing, err := getNote(db, noteID, teamID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Note not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Could not load note", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case "GET":
			renderPage(*existing, "")
		case "POST":
			if err := r.ParseForm(); err != nil {
				renderPage(*existing, "Invalid form data")
				return
			}
			note, err := noteFromForm(r, teamID)
			if err != nil {
				note.ID = existing.ID
				renderPage(note, err.Error())
				return
			}
			if err := validateNoteTarget(db, note); err != nil {
				note.ID = existing.ID
				renderPage(note, err.Error())
				return
			}
			if err := db.Model(existing).Updates(map[string]interface{}{
				"channel":      noteChannelNote,
				"target_kind":  note.TargetKind,
				"target_id":    note.TargetID,
				"target_label": note.TargetLabel,
				"priority":     note.Priority,
				"creator":      note.Creator,
				"type":         note.Type,
				"note":         note.Note,
			}).Error; err != nil {
				renderPage(note, fmt.Sprintf("Could not save note: %v", err))
				return
			}
			if returnTo := noteReturnTo(r); returnTo != "" {
				http.Redirect(w, r, returnTo, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/notes", http.StatusSeeOther)
		case "DELETE":
			if err := db.Delete(existing).Error; err != nil {
				warning(fmt.Sprintf("Could not delete note: %v", err)).Render(GetContext(r, db), w)
				return
			}
			if returnTo := noteReturnTo(r); returnTo != "" {
				w.Header().Set("HX-Redirect", returnTo)
			} else {
				w.Header().Set("HX-Redirect", "/notes")
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func feedbackHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := getTeamId(GetContext(r, db))
		if teamID == 0 {
			http.Error(w, "Team not found", http.StatusNotFound)
			return
		}
		team, err := GetTeam(db, teamID)
		if err != nil || team == nil {
			http.Error(w, "Team not found", http.StatusNotFound)
			return
		}
		if !team.EnablePublicFeedbackForm {
			http.Error(w, "Feedback form is disabled", http.StatusForbidden)
			return
		}
		renderPage := func(selected MLNote, msg string, isError bool) {
			matches, err := feedbackMatches(db, teamID)
			if err != nil {
				warning(fmt.Sprintf("Could not load matches: %v", err)).Render(GetContext(r, db), w)
				return
			}
			feedbackPage(FeedbackPageData{Matches: matches, Selected: selected, Message: msg, IsError: isError}).Render(GetContext(r, db), w)
		}

		switch r.Method {
		case "GET":
			renderPage(MLNote{Channel: noteChannelFeedback, TargetKind: "general", TargetLabel: "General", Type: "general"}, "", false)
		case "POST":
			if err := r.ParseForm(); err != nil {
				renderPage(MLNote{Channel: noteChannelFeedback, TargetKind: "general", TargetLabel: "General", Type: "general"}, "Invalid form data", true)
				return
			}
			feedback, err := feedbackFromForm(db, r, teamID)
			if err != nil {
				renderPage(feedback, err.Error(), true)
				return
			}
			if err := db.Create(&feedback).Error; err != nil {
				renderPage(feedback, fmt.Sprintf("Could not save feedback: %v", err), true)
				return
			}
			renderPage(MLNote{Channel: noteChannelFeedback, TargetKind: "general", TargetLabel: "General", Type: "general"}, "Thanks, feedback saved.", false)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func publicMatchFeedbackHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimSpace(chi.URLParam(r, "matchURLSlug"))
		match, err := GetMatchByURLSlug(db, slug)
		if err != nil {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}
		team, err := GetTeam(db, match.TeamID)
		if err != nil || team == nil {
			http.Error(w, "Team not found", http.StatusNotFound)
			return
		}
		if !team.EnablePublicFeedbackForm {
			http.Error(w, "Feedback form is disabled", http.StatusForbidden)
			return
		}

		matchLabel := feedbackMatchLabel(*match, time.Now())
		selected := MLNote{
			TeamID:      match.TeamID,
			Channel:     noteChannelFeedback,
			TargetKind:  noteTargetMatch,
			TargetID:    match.ID,
			TargetLabel: matchLabel,
			Type:        "general",
		}
		actionURL := publicMatchFeedbackURLPath(*match)
		renderPage := func(selected MLNote, msg string, isError bool) {
			ctx := context.WithValue(GetContext(r, db), "team_id", match.TeamID)
			ctx = context.WithValue(ctx, teamKey, *team)
			feedbackPage(FeedbackPageData{
				Matches:   []NoteTargetOption{{Kind: noteTargetMatch, ID: match.ID, Label: matchLabel}},
				Selected:  selected,
				Message:   msg,
				IsError:   isError,
				ActionURL: actionURL,
				LockMatch: true,
			}).Render(ctx, w)
		}

		switch r.Method {
		case http.MethodGet:
			renderPage(selected, "", false)
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				renderPage(selected, "Invalid form data", true)
				return
			}
			feedback, err := feedbackFromForm(db, r, match.TeamID)
			if err != nil {
				renderPage(feedback, err.Error(), true)
				return
			}
			if feedback.TargetKind != noteTargetMatch || feedback.TargetID != match.ID {
				renderPage(feedback, "match is required", true)
				return
			}
			if err := db.Create(&feedback).Error; err != nil {
				renderPage(feedback, fmt.Sprintf("Could not save feedback: %v", err), true)
				return
			}
			renderPage(selected, "Thanks, feedback saved.", false)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
