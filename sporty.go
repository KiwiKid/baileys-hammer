package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const sportyUpcomingFixturesURL = "https://www.mainlandfootball.co.nz/api/v2/competition/widget/fixture/UpcomingFixtures"
const sportyRecentResultsURL = "https://www.mainlandfootball.co.nz/api/v2/competition/widget/fixture/RecentResults"
const sportyAutoSyncCooldown = 6 * time.Hour

type sportyFixturesRequest struct {
	CompIDs  []uint `json:"CompIds"`
	OrgIDs   []uint `json:"OrgIds"`
	GradeIDs []uint `json:"GradeIds"`
}

type sportyFixturesResponse struct {
	Fixtures []sportyFixture `json:"Fixtures"`
}

type sportyFixture struct {
	ID           uint    `json:"Id"`
	From         string  `json:"From"`
	HomeTeamID   uint    `json:"HomeTeamId"`
	HomeTeamName string  `json:"HomeTeamName"`
	AwayTeamID   uint    `json:"AwayTeamId"`
	AwayTeamName string  `json:"AwayTeamName"`
	VenueName    string  `json:"VenueName"`
	VenueAbbr    string  `json:"VenueAbbr"`
	LocationLat  float64 `json:"LocationLat"`
	LocationLng  float64 `json:"LocationLng"`
	GradeName    string  `json:"GradeName"`
	RoundName    string  `json:"RoundName"`
	StatusName   string  `json:"StatusName"`
}

type SportySyncResult struct {
	Fetched        int
	Matched        int
	Created        int
	Updated        int
	Skipped        int
	SkippedRecords []SportySyncSkippedRecord
	SyncedAt       time.Time
}

type SportySyncSkippedRecord struct {
	Fixture sportyFixture
	Reason  string
}

type SportySyncCandidate struct {
	Fixture sportyFixture
	Match   Match
	IsNew   bool
}

type SportySyncPreview struct {
	Fetched        int
	Matched        int
	Skipped        int
	IncludePast    bool
	SkippedRecords []SportySyncSkippedRecord
	Candidates     []SportySyncCandidate
}

func parseOptionalUint(value string) (uint, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func teamHasSportyConfig(team Team) bool {
	return team.SportyCompetitionID > 0 && team.SportyGradeID > 0 && team.SportyTeamID > 0
}

func sportyFixtureSummary(fixture sportyFixture) string {
	teams := strings.TrimSpace(fmt.Sprintf("%s vs %s", fixture.HomeTeamName, fixture.AwayTeamName))
	if teams == "vs" {
		teams = fmt.Sprintf("Fixture %d", fixture.ID)
	}
	if fixture.From != "" {
		return fmt.Sprintf("%s - %s", teams, fixture.From)
	}
	return teams
}

func (result *SportySyncResult) skipFixture(fixture sportyFixture, reason string) {
	if result == nil {
		return
	}
	result.Skipped++
	result.SkippedRecords = append(result.SkippedRecords, SportySyncSkippedRecord{
		Fixture: fixture,
		Reason:  reason,
	})
}

func AutoSyncSportyMatchesIfDue(db *gorm.DB, team *Team) (*SportySyncResult, error) {
	if team == nil || !team.SportyAutoMatchSync || !teamHasSportyConfig(*team) {
		return nil, nil
	}
	if team.SportyLastSyncedAt != nil && time.Since(*team.SportyLastSyncedAt) < sportyAutoSyncCooldown {
		return nil, nil
	}
	result, err := SyncSportyMatches(db, team)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func SyncSportyMatches(db *gorm.DB, team *Team) (*SportySyncResult, error) {
	return syncSportyMatches(db, team, nil, false)
}

func PreviewSportyMatches(db *gorm.DB, team *Team, includePast bool) (*SportySyncPreview, error) {
	if team == nil {
		return nil, errors.New("team is required")
	}
	if !teamHasSportyConfig(*team) {
		return nil, errors.New("sporty competition, grade, and team IDs are required")
	}

	fixtures, err := fetchSportyFixtures(team, includePast)
	if err != nil {
		return nil, err
	}

	preview := &SportySyncPreview{Fetched: len(fixtures), IncludePast: includePast}
	for _, fixture := range fixtures {
		if fixture.HomeTeamID != team.SportyTeamID && fixture.AwayTeamID != team.SportyTeamID {
			preview.Skipped++
			preview.SkippedRecords = append(preview.SkippedRecords, SportySyncSkippedRecord{
				Fixture: fixture,
				Reason:  fmt.Sprintf("Team ID %d was not listed as the home or away team.", team.SportyTeamID),
			})
			continue
		}
		preview.Matched++
		if reason, skip, err := sportyFixtureSeasonSkipReason(db, fixture); err != nil {
			return nil, err
		} else if skip {
			preview.Skipped++
			preview.SkippedRecords = append(preview.SkippedRecords, SportySyncSkippedRecord{
				Fixture: fixture,
				Reason:  reason,
			})
			continue
		}

		match, isNew, err := matchFromSportyFixture(db, *team, fixture)
		if err != nil {
			return nil, err
		}
		preview.Candidates = append(preview.Candidates, SportySyncCandidate{
			Fixture: fixture,
			Match:   *match,
			IsNew:   isNew,
		})
	}

	return preview, nil
}

func SyncSelectedSportyMatches(db *gorm.DB, team *Team, selectedFixtureIDs map[uint]bool, includePast bool) (*SportySyncResult, error) {
	return syncSportyMatches(db, team, selectedFixtureIDs, true, includePast)
}

func fetchSportyFixtures(team *Team, includePast bool) ([]sportyFixture, error) {
	fixtures, err := fetchSportyFixturesFromURL(team, sportyUpcomingFixturesURL)
	if err != nil {
		return nil, err
	}
	if !includePast {
		return fixtures, nil
	}

	recentFixtures, err := fetchSportyFixturesFromURL(team, sportyRecentResultsURL)
	if err != nil {
		return nil, err
	}

	seenFixtureIDs := map[uint]bool{}
	for _, fixture := range fixtures {
		seenFixtureIDs[fixture.ID] = true
	}
	for _, fixture := range recentFixtures {
		if seenFixtureIDs[fixture.ID] {
			continue
		}
		fixtures = append(fixtures, fixture)
		seenFixtureIDs[fixture.ID] = true
	}
	return fixtures, nil
}

func fetchSportyFixturesFromURL(team *Team, fixturesURL string) ([]sportyFixture, error) {
	payload := sportyFixturesRequest{
		CompIDs:  []uint{team.SportyCompetitionID},
		OrgIDs:   []uint{},
		GradeIDs: []uint{team.SportyGradeID},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fixturesURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.mainlandfootball.co.nz")
	req.Header.Set("Referer", "https://www.mainlandfootball.co.nz/")
	req.Header.Set("User-Agent", "baileys-hammer/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("sporty fixtures request failed with status %s", resp.Status)
	}

	var fixtures sportyFixturesResponse
	if err := json.NewDecoder(resp.Body).Decode(&fixtures); err != nil {
		return nil, err
	}

	return fixtures.Fixtures, nil
}

func syncSportyMatches(db *gorm.DB, team *Team, selectedFixtureIDs map[uint]bool, requireSelection bool, includePast ...bool) (*SportySyncResult, error) {
	if team == nil {
		return nil, errors.New("team is required")
	}
	if !teamHasSportyConfig(*team) {
		return nil, errors.New("sporty competition, grade, and team IDs are required")
	}

	shouldIncludePast := false
	if len(includePast) > 0 {
		shouldIncludePast = includePast[0]
	}
	fixtures, err := fetchSportyFixtures(team, shouldIncludePast)
	if err != nil {
		return nil, err
	}

	result := &SportySyncResult{Fetched: len(fixtures)}
	for _, fixture := range fixtures {
		if fixture.HomeTeamID != team.SportyTeamID && fixture.AwayTeamID != team.SportyTeamID {
			result.skipFixture(fixture, fmt.Sprintf("Team ID %d was not listed as the home or away team.", team.SportyTeamID))
			continue
		}
		result.Matched++
		if requireSelection && !selectedFixtureIDs[fixture.ID] {
			result.skipFixture(fixture, "Not selected in the sync preview.")
			continue
		}
		if reason, skip, err := sportyFixtureSeasonSkipReason(db, fixture); err != nil {
			return nil, err
		} else if skip {
			result.skipFixture(fixture, reason)
			continue
		}

		match, isNew, err := matchFromSportyFixture(db, *team, fixture)
		if err != nil {
			return nil, err
		}
		if err := db.Save(match).Error; err != nil {
			return nil, err
		}
		if isNew {
			result.Created++
		} else {
			result.Updated++
		}
	}

	syncedAt := time.Now()
	if err := db.Model(&Team{}).Where("id = ?", team.ID).Updates(map[string]interface{}{
		"sporty_last_synced_at": syncedAt,
	}).Error; err != nil {
		return nil, err
	}
	team.SportyLastSyncedAt = &syncedAt
	result.SyncedAt = syncedAt
	return result, nil
}

func sportyFixtureStartTime(fixture sportyFixture) (time.Time, error) {
	startTime, err := time.ParseInLocation("2006-01-02T15:04:05", fixture.From, time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse sporty fixture %d start time %q: %w", fixture.ID, fixture.From, err)
	}
	return startTime, nil
}

func sportyFixtureSeasonSkipReason(db *gorm.DB, fixture sportyFixture) (string, bool, error) {
	season, err := GetActiveSeason(db)
	if err != nil {
		return "", false, err
	}
	if season == nil {
		return "No active season is set, so the fixture season cannot be checked.", true, nil
	}
	startTime, err := sportyFixtureStartTime(fixture)
	if err != nil {
		return "", false, err
	}
	if startTime.Before(season.StartDate) {
		return fmt.Sprintf("Fixture date %s is before active season %q starts on %s.", startTime.Format("2006-01-02 15:04"), season.Title, season.StartDate.Format("2006-01-02 15:04")), true, nil
	}
	return "", false, nil
}

func matchFromSportyFixture(db *gorm.DB, team Team, fixture sportyFixture) (*Match, bool, error) {
	var match Match
	err := db.Where("team_id = ? AND sporty_fixture_id = ?", team.ID, fixture.ID).First(&match).Error
	isNew := false
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, err
		}
		isNew = true
		match = Match{
			TeamID:          team.ID,
			SportyFixtureID: fixture.ID,
		}
	}

	startTime, err := sportyFixtureStartTime(fixture)
	if err != nil {
		return nil, false, err
	}

	opponent := fixture.AwayTeamName
	if fixture.AwayTeamID == team.SportyTeamID {
		opponent = fixture.HomeTeamName
	}
	location := fixture.VenueName
	if location == "" {
		location = fixture.VenueAbbr
	}

	subtitleParts := []string{}
	for _, part := range []string{fixture.RoundName, fixture.GradeName, fixture.StatusName} {
		part = strings.TrimSpace(part)
		if part != "" {
			subtitleParts = append(subtitleParts, part)
		}
	}

	match.TeamID = team.ID
	if match.SeasonId == 0 {
		match.SeasonId = uint64(activeSeasonID(db))
	}
	match.SportyFixtureID = fixture.ID
	match.Opponent = strings.TrimSpace(opponent)
	match.Location = strings.TrimSpace(location)
	match.StartTime = &startTime
	match.MatchLat = fixture.LocationLat
	match.MatchLng = fixture.LocationLng
	match.Subtitle = strings.Join(subtitleParts, " - ")
	return &match, isNew, nil
}

func sportySyncSummary(result *SportySyncResult) string {
	if result == nil {
		return "Sporty sync skipped"
	}
	return fmt.Sprintf("Sporty sync complete: %d matched, %d created, %d updated, %d skipped", result.Matched, result.Created, result.Updated, result.Skipped)
}

func selectedSportyFixtureIDs(r *http.Request) (map[uint]bool, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	selected := map[uint]bool{}
	for _, value := range r.Form["sportyFixtureId"] {
		fixtureID, err := parseOptionalUint(value)
		if err != nil {
			return nil, err
		}
		if fixtureID > 0 {
			selected[fixtureID] = true
		}
	}
	return selected, nil
}

func includePastSportyMatches(r *http.Request) bool {
	value := strings.TrimSpace(r.FormValue("includePastSportyMatches"))
	return value == "on" || value == "1" || value == "true"
}

func sportySyncHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		teamID, err := parseOptionalUint(r.URL.Query().Get("teamId"))
		if err != nil || teamID == 0 {
			warning("Invalid team ID").Render(GetContext(r, db), w)
			return
		}

		team, err := GetTeam(db, teamID)
		if err != nil {
			warning(fmt.Sprintf("Error fetching team: %v", err)).Render(GetContext(r, db), w)
			return
		}
		includePast := includePastSportyMatches(r)

		if r.URL.Query().Get("confirm") == "1" {
			selectedFixtureIDs, err := selectedSportyFixtureIDs(r)
			if err != nil {
				matches, matchErr := GetMatches(db, activeSeasonID(db), 0, 999)
				if matchErr != nil {
					warning(fmt.Sprintf("Invalid Sporty selection: %v", err)).Render(GetContext(r, db), w)
					return
				}
				teamEditForm(*team, matches, warning(fmt.Sprintf("Invalid Sporty selection: %v", err))).Render(GetContext(r, db), w)
				return
			}
			result, err := SyncSelectedSportyMatches(db, team, selectedFixtureIDs, includePast)
			if err != nil {
				matches, matchErr := GetMatches(db, activeSeasonID(db), 0, 999)
				if matchErr != nil {
					warning(fmt.Sprintf("Sporty sync failed: %v", err)).Render(GetContext(r, db), w)
					return
				}
				teamEditForm(*team, matches, warning(fmt.Sprintf("Sporty sync failed: %v", err))).Render(GetContext(r, db), w)
				return
			}
			matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
			if err != nil {
				warning(fmt.Sprintf("Error fetching matches: %v", err)).Render(GetContext(r, db), w)
				return
			}
			teamEditForm(*team, matches, sportySyncResultMessage("", result, false)).Render(GetContext(r, db), w)
			return
		}

		matches, err := GetMatches(db, activeSeasonID(db), 0, 999)
		if err != nil {
			warning(fmt.Sprintf("Error fetching matches: %v", err)).Render(GetContext(r, db), w)
			return
		}

		preview, err := PreviewSportyMatches(db, team, includePast)
		if err != nil {
			teamEditForm(*team, matches, warning(fmt.Sprintf("Sporty sync preview failed: %v", err))).Render(GetContext(r, db), w)
			return
		}

		teamEditForm(*team, matches, sportySyncPreview(*team, preview)).Render(GetContext(r, db), w)
	}
}

func logAutoSportySync(db *gorm.DB, teamID uint) {
	if teamID == 0 {
		return
	}
	team, err := GetTeam(db, teamID)
	if err != nil {
		log.Printf("sporty auto-sync: failed to load team %d: %v", teamID, err)
		return
	}
	result, err := AutoSyncSportyMatchesIfDue(db, team)
	if err != nil {
		log.Printf("sporty auto-sync failed for team %d: %v", teamID, err)
		return
	}
	if result != nil {
		log.Printf("sporty auto-sync team %d: %s", teamID, sportySyncSummary(result))
	}
}
