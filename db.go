package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Player represents a player in the social football team
type Player struct {
	gorm.Model
	Name                   string `gorm:"unique" form:"name"`
	Fines                  []Fine
	Active                 bool
	Role                   string
	RoleDescription        string
	LatestSeasonId         int     `form:"seasonId" json:"seasonId"`
	SubsOutstandingAmount  float64 `form:"subsOutstandingAmount"`
	FinesOutstandingAmount float64 `form:"finesOutstandingAmount"`
}

type PlayerPayment struct {
	gorm.Model
	PlayerID        uint
	SeasonID        uint
	Amount          float64
	PaymentLoggedAt time.Time
}

type ImageArray struct {
	Images []string
}

// Value implements the driver.Valuer interface, converting LatLngArray to a JSON-encoded byte array.
func (l ImageArray) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *ImageArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, l)
}

type FineImage struct {
	gorm.Model
	Filename string
	FineID   uint
	Image    string
	Data     []byte
}

// Fine represents a fine assigned to a player
type Fine struct {
	gorm.Model
	PlayerID          uint
	SeasonID          uint
	CourtSessionOrder uint
	CourtSessionNote  string
	MatchId           uint
	FineAt            time.Time
	Reason            string
	Amount            float64
	Approved          bool
	Context           string
	Contest           string
	Images            ImageArray
}

type Team struct {
	gorm.Model
	ID                    uint
	TeamName              string
	TeamKey               string
	TeamAdminPass         string
	TeamMemberPass        string
	ShowFineAddOnHomePage bool
	ActiveMatchIDOverride uint
	ShowCourtTotals       bool
	CourtNotes            string
}

// DBInit initializes the database and creates the tables
func DBInit() (*gorm.DB, error) {

	dbUrl := os.Getenv("DATABASE_URL")
	if len(dbUrl) == 0 {
		//log.Panic("No DATABASE_URL set")
	}

	log.Printf("connecting to db: %s", dbUrl)
	db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Printf("Could not access :%s - %v", dbUrl, err)
		return nil, err
	} else {
		log.Printf("Connected to db at \"%s\"", dbUrl)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Player{}, &Fine{}, &PresetFine{}, &Match{}, &MatchEvent{}, &FineImage{}, &Season{}, &PlayerPayment{}, &Team{})
	if err != nil {
		return nil, err
	}

	//db.Migrator().DropTable(&Match{})
	if (!db.Migrator().HasColumn(&Match{}, "SeasonId")) {
		db.Migrator().AddColumn(&Match{}, "SeasonId")
	}

	if (!db.Migrator().HasColumn(&Match{}, "PlayerOfTheDay")) {
		db.Migrator().AddColumn(&Match{}, "PlayerOfTheDay")
	}

	if (!db.Migrator().HasColumn(&Match{}, "DudOfTheDay")) {
		db.Migrator().AddColumn(&Match{}, "DudOfTheDay")
	}

	if (!db.Migrator().HasColumn(&Match{}, "MatchLat")) {
		db.Migrator().AddColumn(&Match{}, "MatchLat")
	}

	if (!db.Migrator().HasColumn(&Match{}, "MatchLng")) {
		db.Migrator().AddColumn(&Match{}, "MatchLng")
	}

	if (!db.Migrator().HasColumn(&Match{}, "MatchRectangle")) {
		db.Migrator().AddColumn(&Match{}, "MatchRectangle")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "Context")) {
		db.Migrator().AddColumn(&Fine{}, "Context")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "CourtSessionOrder")) {
		db.Migrator().AddColumn(&Fine{}, "CourtSessionOrder")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "CourtSessionNote")) {
		db.Migrator().AddColumn(&Fine{}, "CourtSessionNote")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "MatchId")) {
		db.Migrator().AddColumn(&Fine{}, "MatchId")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "Contest")) {
		db.Migrator().AddColumn(&Fine{}, "Contest")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "FineAt")) {
		db.Migrator().AddColumn(&Fine{}, "FineAt")
	}

	if (!db.Migrator().HasColumn(&Fine{}, "SeasonId")) {
		db.Migrator().AddColumn(&Fine{}, "SeasonId")
	}

	if (!db.Migrator().HasColumn(&Player{}, "RoleDescription")) {
		db.Migrator().AddColumn(&Player{}, "RoleDescription")
	}

	if (!db.Migrator().HasColumn(&Player{}, "Role")) {
		db.Migrator().AddColumn(&Player{}, "Role")
	}

	if (!db.Migrator().HasColumn(&Player{}, "SubsOutstandingAmount")) {
		db.Migrator().AddColumn(&Player{}, "SubsOutstandingAmount")
	}

	if (!db.Migrator().HasColumn(&Player{}, "FinesOutstandingAmount")) {
		db.Migrator().AddColumn(&Player{}, "FinesOutstandingAmount")
	}

	if (!db.Migrator().HasColumn(&PresetFine{}, "Context")) {
		db.Migrator().AddColumn(&PresetFine{}, "Context")
	}

	if (!db.Migrator().HasColumn(&PresetFine{}, "NotQuickFine")) {
		db.Migrator().AddColumn(&PresetFine{}, "NotQuickFine")
	}

	if (!db.Migrator().HasColumn(&PresetFine{}, "DisplayOrder")) {
		db.Migrator().AddColumn(&PresetFine{}, "DisplayOrder")
	}

	if (!db.Migrator().HasColumn(&PresetFine{}, "Icon")) {
		db.Migrator().AddColumn(&PresetFine{}, "Icon")
	}

	if (!db.Migrator().HasColumn(&PresetFine{}, "IsKudos")) {
		db.Migrator().AddColumn(&PresetFine{}, "IsKudos")
	}

	result := db.Model(&Fine{}).Where("fine_at IS NULL").Updates(map[string]interface{}{
		"fine_at": gorm.Expr("created_at"),
	})

	if result.Error != nil {
		panic(result.Error)
	}

	return db, nil
}

func GetTeam(db *gorm.DB, id uint) (*Team, error) {
	var team Team
	result := db.First(&team, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func GetTeams(db *gorm.DB, limit int, offset int) ([]Team, error) {
	var teams []Team
	result := db.Limit(limit).Offset(offset).Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func SaveTeam(db *gorm.DB, team *Team) (Team, error) {
	if err := db.Save(team).Error; err != nil {
		return Team{}, err
	}
	return *team, nil
}

func DeleteTeam(db *gorm.DB, id uint) error {
	err := db.Delete(&Team{}, id)
	if err.Error != nil {
		return err.Error
	}
	if err.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetTeamByKeyAndPassword authenticates a team by key and admin password
func GetTeamByKeyAndPassword(db *gorm.DB, teamKey, adminPassword string) (*Team, error) {
	var team Team
	result := db.Where("team_key = ? AND team_admin_pass = ?", teamKey, adminPassword).First(&team)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

// PlayerWithFines represents a player along with their fines
type PlayerWithFines struct {
	ID                     uint
	Name                   string
	TotalFineCount         int
	TotalFines             int
	Role                   string
	RoleDescription        string
	Fines                  []Fine
	PendingFines           []Fine
	PendingFineSum         int
	PendingTotalCount      int
	FinesOutstandingAmount float64
	SubsOutstandingAmount  float64
}

func GetPlayersWithFines(db *gorm.DB, seasonId uint64, playerIds []uint64) ([]PlayerWithFines, error) {
	var playersWithFines []PlayerWithFines

	var players []Player

	if len(playerIds) > 0 {
		db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
			return db.Order("fines.fine_at DESC")
		}).Where("id IN ?", playerIds).Find(&players).Order("players.name")

	} else {
		if seasonId == 0 {
			db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
				return db.Order("fines.fine_at DESC")
			}).Find(&players).Where("season_ids IN ?", seasonId).Order("players.name")
		} else {
			db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
				return db.Order("fines.fine_at DESC")
			}).Find(&players).Order("players.name")
		}

	}
	//.Where("active = true")

	// Construct the PlayerWithFines slice
	for _, player := range players {

		var approvedFines = []Fine{}
		var pendingFines = []Fine{}
		var fineSum = 0
		var pendingSum = 0
		for _, f := range player.Fines {
			if f.Approved {
				approvedFines = append(approvedFines, f)
				fineSum = fineSum + int(f.Amount)
			} else {
				pendingFines = append(pendingFines, f)
				pendingSum = pendingSum + int(f.Amount)
			}
		}
		pwf := PlayerWithFines{
			ID:                     player.ID,
			Name:                   player.Name,
			TotalFineCount:         len(approvedFines),
			TotalFines:             fineSum,
			Role:                   player.Role,
			RoleDescription:        player.RoleDescription,
			Fines:                  approvedFines,
			PendingFines:           pendingFines,
			PendingFineSum:         pendingSum,
			PendingTotalCount:      len(pendingFines),
			FinesOutstandingAmount: player.FinesOutstandingAmount,
			SubsOutstandingAmount:  player.SubsOutstandingAmount,
		}
		playersWithFines = append(playersWithFines, pwf)
	}

	sort.Slice(playersWithFines, func(i, j int) bool {
		return playersWithFines[i].TotalFines > playersWithFines[j].TotalFines
	})

	return playersWithFines, nil
}

func SetAllFineAmounts(db *gorm.DB, amount float64) (int64, int64, error) {
	// Create a map with the fields you want to update
	updates := map[string]interface{}{
		"Amount": amount,
	}

	// Find the Fine by ID and update the Contest field
	result := db.Model(&Fine{}).Where("amount < ? OR amount is null", amount).Updates(updates)

	// Check if the update operation resulted in an error
	if result.Error != nil {
		return 0, 0, result.Error
	}

	result2 := db.Model(&PresetFine{}).Where("amount < ? OR amount is null", amount).Updates(updates)
	// Check if the update operation resulted in an error
	if result2.Error != nil {
		return result.RowsAffected, 0, result2.Error
	}

	return result.RowsAffected, result2.RowsAffected, nil
}

func SetSeasonId(db *gorm.DB, days int) (int64, int64, error) {
	activeSeason, err := GetActiveSeason(db)
	if err != nil {
		return 0, 0, err
	}
	if activeSeason == nil {
		return 0, 0, errors.New("no active season found")
	}

	playerUpdates := map[string]interface{}{
		"LatestSeasonId": activeSeason.ID,
	}

	playerResult := db.Model(&Player{}).Where("created_at > ?", time.Now().AddDate(0, 0, -days)).Updates(playerUpdates)

	if playerResult.Error != nil {
		return 0, 0, playerResult.Error
	}

	updates := map[string]interface{}{
		"SeasonID": activeSeason.ID,
	}

	result := db.Model(&Fine{}).Where("created_at > ? AND (season_id != ? OR season_id is null)", time.Now().AddDate(0, 0, -days), activeSeason.ID).Updates(updates)
	if result.Error != nil {
		return 0, playerResult.RowsAffected, result.Error
	}

	return result.RowsAffected, playerResult.RowsAffected, nil
}

func SetPlayerFinesOutstandingForActiveSeason(db *gorm.DB, seasonID uint) (int, error) {

	playerWithFines, err := GetPlayersWithFines(db, seasonId, []uint64{})
	if err != nil {
		return 0, err
	}

	var rows = 0
	for _, player := range playerWithFines {

		var total = 0.0
		for _, fine := range player.PendingFines {
			total = total + fine.Amount
		}

		updates := map[string]interface{}{
			"FinesOutstandingAmount": total,
		}

		result := db.Model(&Player{}).Where("id = ?", player.ID).Updates(updates)
		if result.Error != nil {
			return 0, result.Error
		} else {
			rows = rows + int(result.RowsAffected)

		}

	}

	return rows, nil
}
func SetPlayerSubsOutstandingForActiveSeason(db *gorm.DB, seasonID uint, seasonSubs int) (int, error) {

	updates := map[string]interface{}{
		"SubsOutstandingAmount": seasonSubs,
	}

	result := db.Model(&Player{}).Where("latest_season_id = ?", seasonID).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func CreatePlayerPayment(db *gorm.DB, playerId uint, amount float64, seasonId uint) error {
	payment := PlayerPayment{
		PlayerID:        playerId,
		Amount:          amount,
		PaymentLoggedAt: time.Now(),
		SeasonID:        seasonId,
	}

	if err := db.Save(&payment).Error; err != nil {
		return err
	}
	return nil
}

func GetPlayerPayments(db *gorm.DB, seasonId uint) ([]PlayerPayment, error) {
	var payments []PlayerPayment
	if err := db.Where("season_id = ?", seasonId).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func DeletePlayerPayment(db *gorm.DB, paymentId uint) error {
	result := db.Delete(&PlayerPayment{}, paymentId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func FetchActivePlayers(db *gorm.DB) ([]Player, error) {
	var players []Player

	// Query all active players
	result := db.Where("active = ?", true).Find(&players)
	if result.Error != nil {
		// Return an empty slice and the error
		return []Player{}, result.Error
	}

	// Return the slice of active players and nil error
	return players, nil
}

func GetFineByID(db *gorm.DB, fineID uint) (*Fine, error) {
	var fine Fine
	result := db.First(&fine, fineID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fine, nil
}

func GetPlayerByID(db *gorm.DB, playerID uint) (*Player, error) {
	var player Player
	result := db.First(&player, playerID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &player, nil
}

func GetMatchMetaGeneral(db *gorm.DB, matchId uint) (*MatchMetaGeneral, error) {
	match, err := GetMatchWithEvents(db, matchId)
	if err != nil {
		return nil, err
	}

	genMeta, err := WrapMatchWithMeta(db, *match)
	if err != nil {
		return nil, err
	}

	return genMeta, nil
}

func SavePlayer(db *gorm.DB, player *Player) error {
	// Create or update player
	if err := db.Save(player).Error; err != nil {
		return err
	}
	return nil
}

func DeletePlayer(db *gorm.DB, playerId uint) error {

	if err := db.Model(&Fine{}).Delete("playerId = ?", playerId).Error; err != nil {
		return err
	}

	if err := db.Model(&Player{}).Delete("id = ?", playerId).Error; err != nil {
		return err
	}
	return nil
}

func FetchLatestFines(db *gorm.DB, page int, pageSize int) ([]Fine, error) {
	var fines []Fine
	offset := (page - 1) * pageSize

	result := db.Order("fine_at DESC").Offset(offset).Limit(pageSize).Find(&fines)
	if result.Error != nil {
		return nil, result.Error
	}

	return fines, nil
}

// FetchLatestFines fetches a paginated list of the latest fines.
func GetPlayers(db *gorm.DB, page int, pageSize int) ([]Player, error) {
	var players []Player
	offset := (page - 1) * pageSize

	// Query the latest fines with pagination
	result := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&players)
	if result.Error != nil {
		return nil, result.Error
	}

	return players, nil
}

func UpdateFineCourtSessionOrderByID(db *gorm.DB, fineID uint, courtSessionOrder uint, courtSessionNotes string) error {
	// Create a map with the fields you want to update
	updates := map[string]interface{}{
		"CourtSessionOrder": courtSessionOrder,
		"CourtSessionNote":  courtSessionNotes,
	}
	result := db.Model(&Fine{}).Where("id = ?", fineID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	log.Printf("Fine %d updated with court session order %d", fineID, courtSessionOrder)
	return nil
}

func SaveFine(db *gorm.DB, fine *Fine) error {
	if err := db.Save(fine).Error; err != nil {
		return err
	}
	return nil
}

func UpdateFineContestByID(db *gorm.DB, fineID uint, contest string) error {
	// Create a map with the fields you want to update
	updates := map[string]interface{}{
		"Contest": contest,
	}

	// Find the Fine by ID and update the Contest field
	result := db.Model(&Fine{}).Where("id = ?", fineID).Updates(updates)

	// Check if the update operation resulted in an error
	if result.Error != nil {
		return result.Error
	}

	// Optionally, check if the record was found and updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("no fine found with ID %d", fineID)
	}

	return nil
}

func UpdateFineContextByID(db *gorm.DB, fineID uint, matchId uint, context string, fineAt *time.Time) error {
	// Create a map with the fields you want to update
	updates := map[string]interface{}{
		"Context": context,
		"MatchId": matchId,
	}

	if fineAt != nil {
		updates["FineAt"] = *fineAt
	}

	// Find the Fine by ID and update the Contest field
	result := db.Model(&Fine{}).Where("id = ?", fineID).Updates(updates)

	// Check if the update operation resulted in an error
	if result.Error != nil {
		return result.Error
	}

	// Optionally, check if the record was found and updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("no fine found with ID %d", fineID)
	}

	return nil
}

func ApproveFine(db *gorm.DB, id uint, amount float64, approved bool) error {
	// Find and update the fine's Approved field to true
	updates := map[string]interface{}{
		"approved": approved,
		"amount":   amount,
	}

	// Find and update the fine's Approved and Amount fields
	result := db.Model(&Fine{}).Where("id = ?", id).Updates(updates)

	// Check for errors during the operation
	if result.Error != nil {
		return result.Error
	}

	// Check if the record was found and updated
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	log.Printf("Fine %d approved! %f", id, amount)

	return nil
}

func QuickHideFine(db *gorm.DB, id uint, showOrHide bool) error {
	// Find and update the fine's Approved field to true
	result := db.Model(&PresetFine{}).Where("id = ?", id).Update("not_quick_fine", showOrHide)

	// Check for errors during the operation
	if result.Error != nil {
		return result.Error
	}

	// Check if the record was found and updated
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func ApprovePresetFine(db *gorm.DB, id uint) error {
	// Find and update the fine's Approved field to true
	result := db.Model(&PresetFine{}).Where("id = ?", id).Update("approved", true)

	// Check for errors during the operation
	if result.Error != nil {
		return result.Error
	}

	// Check if the record was found and updated
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteFineByID(db *gorm.DB, id uint) error {
	result := db.Delete(&Fine{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetFineImages(db *gorm.DB, fineId uint) ([]FineImage, error) {
	var images []FineImage
	if err := db.Where("fine_id = ?", fineId).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

type PresetFine struct {
	gorm.Model
	Reason       string
	Amount       float64
	Approved     bool
	NotQuickFine bool
	Context      string
	DisplayOrder int64
	Icon         string
	IsKudos      bool
}

func GetFineWithPlayers(db *gorm.DB, pageId uint64, limit int) ([]FineWithPlayer, error) {
	fines, getFErr := FetchLatestFines(db, int(pageId), int(limit))
	if getFErr != nil {
		return []FineWithPlayer{}, getFErr
	}

	players, getPlayerErr := GetPlayers(db, 0, 1000)
	if getPlayerErr != nil {
		return []FineWithPlayer{}, getPlayerErr
	}

	// Get all relevant matches
	matches, getMatchErr := GetMatches(db, 0, 0, 1000)
	if getMatchErr != nil {
		return []FineWithPlayer{}, getMatchErr
	}

	matchMap := make(map[uint]Match)
	for _, match := range matches {
		matchMap[match.ID] = match
	}

	var fineWithPlayers []FineWithPlayer

	// sort the fines by CourtSessionOrder, 0 at the bottom
	sort.Slice(fines, func(i, j int) bool {
		if fines[i].CourtSessionOrder == 0 {
			return false
		}
		if fines[j].CourtSessionOrder == 0 {
			return true
		}

		return fines[i].CourtSessionOrder < fines[j].CourtSessionOrder
	})

	for _, fine := range fines {
		var matchedPlayer Player
		for _, player := range players {
			if fine.PlayerID == player.ID {
				matchedPlayer = player
				break
			}
		}

		matchedMatch := matchMap[fine.MatchId]

		fineWithPlayers = append(fineWithPlayers, FineWithPlayer{
			Fine:   fine,
			Player: matchedPlayer,
			Match:  matchedMatch,
		})
	}

	return fineWithPlayers, nil
}

func SavePresetFine(db *gorm.DB, presetFine *PresetFine) error {
	if err := db.Save(presetFine).Error; err != nil {
		return err
	}
	return nil
}

func DeletePresetFineByID(db *gorm.DB, id uint) error {
	result := db.Delete(&PresetFine{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetPresetFines(db *gorm.DB, includeUnapproved bool) ([]PresetFine, error) {
	var presetFines []PresetFine
	if err := db.Find(&presetFines).Where("approved = ?", includeUnapproved).Order("display_order").Error; err != nil {
		return nil, err
	}
	return presetFines, nil
}

func GetPresetFine(db *gorm.DB, id uint64, reasonName string) (*PresetFine, error) {
	var presetFines PresetFine
	if err := db.Where("id = ?", id).Or("reason = ?", reasonName).First(&presetFines).Error; err != nil {
		return nil, err
	}
	return &presetFines, nil
}

func GetFineFromPreset(db *gorm.DB, pfIDOrReason string) (*Fine, error) {
	pfId, err := strconv.ParseUint(pfIDOrReason, 10, 64)
	var reasonStr = ""
	if err != nil {
		reasonStr = pfIDOrReason
	}
	presetFine, err := GetPresetFine(db, pfId, reasonStr)
	if err != nil {
		return nil, err
	}
	var fine Fine
	if presetFine != nil {
		fine = Fine{
			Amount: presetFine.Amount,
			Reason: presetFine.Reason,
		}
		return &fine, nil
	}

	return nil, errors.New("failed to tget fine from preset")
}

func GetActiveMatch(db *gorm.DB) (*Match, error) {

	activeOverrideId := os.Getenv("ACTIVE_MATCH_OVERRIDE")
	if len(activeOverrideId) > 0 {
		id, err := strconv.ParseUint(activeOverrideId, 10, 64)
		if err != nil {
			return nil, err
		}
		return GetMatch(db, uint(id))
	}

	var match Match
	now := time.Now().Add(-72 * time.Hour)

	// Find the nearest upcoming match where the start time is in the future
	result := db.Where("start_time > ?", now).Order("start_time ASC").First(&match)

	if result.Error != nil {
		/*if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No active match found - create one in the future to allow for it to be associated with a game ACTIVE_MATCH_OVERRIDE:%s", activeOverrideId)
			return nil, result.Error
		}*/
		// Some other error occurred; return the error
		return nil, result.Error
	}

	return &match, nil
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// LatLngArray is a slice of LatLng
type LatLngArray []LatLng

// Value implements the driver.Valuer interface, converting LatLngArray to a JSON-encoded byte array.
func (l LatLngArray) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *LatLngArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, l)
}

type Season struct {
	gorm.Model
	Title     string
	StartDate time.Time
	IsActive  bool
}

type Match struct {
	gorm.Model
	Location       string
	StartTime      *time.Time `json:"timestamp" gorm:"type:datetime"`
	Opponent       string
	Subtitle       string
	Events         []MatchEvent `gorm:"foreignKey:MatchId"`
	SeasonId       uint64
	PlayerOfTheDay uint
	DudOfTheDay    uint
	MatchLat       float64
	MatchLng       float64
	MatchPointList LatLngArray `json:"coords"`
}

type DrinkPayment struct {
	gorm.Model
	PlayerId uint
	Amount   float64
}

type DrinkPurchase struct {
	gorm.Model
	DrinkTypeId uint
	PlayerId    uint
	Amount      float64
}

type DrinkType struct {
	gorm.Model
	Name  string
	Price float64
}

type MatchEvent struct {
	gorm.Model
	MatchId     uint64
	EventName   string
	EventType   string     // 'subbed-off' / 'subbed-on' / 'goal' / 'assist' / 'own-goal' / 'concded-goal'
	EventTime   *time.Time `json:"timestamp" gorm:"type:datetime"`
	EventMinute int
	PlayerId    uint
}

type PlayerState struct {
	PlayerId   uint
	PlayerName string
	TimePlayed int // in minutes
	Goals      int
	Assists    int
}

type MatchState struct {
	PlayersOn            []PlayerState
	MatchID              uint
	ScoreAgainst         int
	TrainingTotalNumbers int
	ScoreFor             int
	MatchDuration        int // Duration of the match in minutes so far
}

// FetchLatestFines fetches a paginated list of the latest fines.
func GetMatches(db *gorm.DB, season uint, page int, pageSize int) ([]Match, error) {
	var matches []Match
	offset := (page - 1) * pageSize
	//.Where("season_id = ?", season)
	result := db.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&matches)
	if result.Error != nil {
		return nil, result.Error
	}

	return matches, nil
}

func GetMatch(db *gorm.DB, id uint) (*Match, error) {
	var match Match
	if err := db.Where("id = ?", id).First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func GetSeasons(db *gorm.DB) ([]Season, error) {
	var seasons []Season
	if err := db.Find(&seasons).Error; err != nil {
		return nil, err
	}
	return seasons, nil
}

func GetSeasonById(db *gorm.DB, id uint) (*Season, error) {
	var season Season
	if err := db.Where("id = ?", id).First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

func DeleteSeason(db *gorm.DB, seasonId uint) error {

	// update all fines to not have a season
	result := db.Model(&Fine{}).Where("season_id = ?", seasonId).Update("season_id", nil)
	if result.Error != nil {
		return result.Error
	}

	result2 := db.Delete(&Season{}, seasonId)
	if result2.Error != nil {
		return result2.Error
	}
	if result2.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type MatchSeasonTeam struct {
	Season *Season
	Team   *Team
	Match  *Match
}

func GetMatchSeasonTeam(db *gorm.DB) (*MatchSeasonTeam, error) {
	season, err := GetActiveSeason(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &MatchSeasonTeam{
				Season: nil,
				Team:   nil,
				Match:  nil,
			}, nil
		} else {
			return nil, err
		}
	}

	team, err := GetActiveTeam(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &MatchSeasonTeam{
				Season: season,
				Team:   nil,
				Match:  nil,
			}, nil
		} else {
			return nil, err
		}
	}

	/*var match *MatchMetaGeneral
	if team != nil {
		if team.ActiveMatchIDOverride > 0 {
			match, err = GetMatchMetaGeneral(db, uint(team.ActiveMatchIDOverride))
			if err != nil {
				return nil, err
			}

		}
	}*/

	match, err := GetActiveMatch(db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &MatchSeasonTeam{
				Season: season,
				Team:   team,
				Match:  nil,
			}, nil
		} else {
			return nil, err
		}
	}

	return &MatchSeasonTeam{
		Season: season,
		Team:   team,
		Match:  match,
	}, nil

}

func GetActiveSeason(db *gorm.DB) (*Season, error) {
	var season Season
	if err := db.Where("is_active = ?", true).Order("start_date desc").First(&season).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active season, return nil instead of an error
		}
		return nil, err // Other errors should be returned
	}
	return &season, nil
}

func GetActiveTeam(db *gorm.DB) (*Team, error) {
	var team Team
	if err := db.First(&team).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active team, return nil instead of an error
		}
		return nil, err // Other errors should be returned
	}
	return &team, nil
}

func WrapMatchWithMeta(db *gorm.DB, match Match) (*MatchMetaGeneral, error) {
	log.Printf("🎁 WrapMatchWithMeta START")
	players, err := GetPlayers(db, 0, 999)
	if err != nil {
		return nil, err
	}

	var playerOfTheDay *Player
	var dudOfTheDay *Player
	var goalScorers []Player = []Player{}
	var goalAssisters []Player = []Player{}
	var opponentGoalCount uint = 0
	for i := range players {
		p := &players[i]
		if p.ID == uint(match.PlayerOfTheDay) {
			playerOfTheDay = p
		}
		if p.ID == uint(match.DudOfTheDay) {
			dudOfTheDay = p
		}
		for _, e := range match.Events {
			if e.EventType == "goal" && p.ID == e.PlayerId {
				goalScorers = append(goalScorers, *p)
			}

			if e.EventType == "assist" && p.ID == e.PlayerId {
				goalAssisters = append(goalAssisters, *p)
			}
			if e.EventType == "conceded-goal" && i == 0 {
				opponentGoalCount = opponentGoalCount + 1
			}
		}
	}

	log.Printf("🎁 WrapMatchWithMeta END - %d:%v %d:%v", match.PlayerOfTheDay, playerOfTheDay, match.DudOfTheDay, dudOfTheDay)
	return &MatchMetaGeneral{
		Match:             match,
		GoalScorers:       goalScorers,
		GoalAssisters:     goalAssisters,
		OpponentGoalCount: opponentGoalCount,
		Players:           players,
		PlayerOfTheDay:    playerOfTheDay,
		DudOfTheDay:       dudOfTheDay,
	}, nil
}

func GetMatchWithEvents(db *gorm.DB, id uint) (*Match, error) {
	var match Match
	if err := db.Preload("Events").Where("id = ?", id).First(&match).Error; err != nil {
		return nil, err
	}

	return &match, nil
}

func GetFinesByMatchId(db *gorm.DB, matchId uint) ([]Fine, error) {
	var fines []Fine
	if err := db.Where("match_id = ?", matchId).Find(&fines).Error; err != nil {
		return nil, err
	}
	return fines, nil
}

func SetFineStartTime(db *gorm.DB, fineID uint, fineAt time.Time) error {
	// Create a map with the fields you want to update
	updates := map[string]interface{}{
		"FineAt": fineAt,
	}

	// Find the Fine by ID and update the Contest field
	result := db.Model(&Fine{}).Where("id = ?", fineID).Updates(updates)

	// Check if the update operation resulted in an error
	if result.Error != nil {
		return result.Error
	}

	// Optionally, check if the record was found and updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("no fine found with ID %d", fineID)
	}

	return nil
}

func SaveSeason(db *gorm.DB, season *Season) error {
	if err := db.Save(season).Error; err != nil {
		return err
	}
	return nil
}

func SaveMatch(db *gorm.DB, match *Match) (uint, error) {
	if err := db.Save(match).Error; err != nil {
		return 0, err // Return 0 as the ID in case of error
	}
	return match.ID, nil // Return the new ID which should now be populated
}

func GetMatchEvents(db *gorm.DB, id uint) ([]MatchEvent, error) {
	var events []MatchEvent
	if err := db.Where("match_id = ?", id).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func GetMatchEvent(db *gorm.DB, id uint) (*MatchEvent, error) {
	var event MatchEvent
	if err := db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func DeleteMatchEvent(db *gorm.DB, eventId uint) error {
	result := db.Delete(&MatchEvent{}, eventId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func SaveMatchEvent(db *gorm.DB, fine *MatchEvent) error {
	if err := db.Save(fine).Error; err != nil {
		return err
	}
	return nil
}

func DeleteMatch(db *gorm.DB, matchId uint) error {
	result := db.Delete(&Match{}, matchId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetDrinkTypes(db *gorm.DB) ([]DrinkType, error) {
	var drinkTypes []DrinkType
	if err := db.Find(&drinkTypes).Error; err != nil {
		return nil, err
	}
	return drinkTypes, nil
}

func GetDrinkPayments(db *gorm.DB) ([]DrinkPayment, error) {
	var drinkPayments []DrinkPayment
	if err := db.Find(&drinkPayments).Error; err != nil {
		return nil, err
	}
	return drinkPayments, nil
}

func AddDrinkPayment(db *gorm.DB, drinkPayment *DrinkPayment) error {
	if err := db.Create(drinkPayment).Error; err != nil {
		return err
	}
	return nil
}

func GetDrinkPurchases(db *gorm.DB) ([]DrinkPurchase, error) {
	var drinkPurchases []DrinkPurchase
	if err := db.Find(&drinkPurchases).Error; err != nil {
		return nil, err
	}
	return drinkPurchases, nil
}

func AddDrinkPurchase(db *gorm.DB, drinkPurchase *DrinkPurchase) error {
	if err := db.Create(drinkPurchase).Error; err != nil {
		return err
	}
	return nil
}

func AddDrinkType(db *gorm.DB, drinkType *DrinkType) error {
	if err := db.Create(drinkType).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	ID       uint     `gorm:"primaryKey"`
	Username string   `gorm:"uniqueIndex"`
	Roles    []string `gorm:"-"`
}

const userContextKey contextKey = "user"
