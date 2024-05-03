package main

import (
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
    Name   string `gorm:"unique" form:"name"`
    Fines  []Fine
    Active bool
    Role string
    RoleDescription string
}

// Fine represents a fine assigned to a player
type Fine struct {
    gorm.Model
    PlayerID uint
    MatchId uint
    FineAt time.Time
    Reason   string
    Amount   float64
    Approved bool
    Context string
    Contest string
}


// DBInit initializes the database and creates the tables
func DBInit() (*gorm.DB, error) {

    dbUrl := os.Getenv("DATABASE_URL")
    if len(dbUrl) == 0 {
        log.Panic("No DATABASE_URL set")
    }
    db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
    if err != nil {
        return nil, err
    }else{
        log.Printf("Connected to db at \"%s\"", dbUrl)
    }

    // Migrate the schema
    err = db.AutoMigrate(&Player{}, &Fine{}, &PresetFine{}, &Match{}, &MatchEvent{})
    if err != nil {
        return nil, err
    }

    //db.Migrator().DropTable(&Match{})
    if(!db.Migrator().HasColumn(&Match{}, "SeasonId")){
        db.Migrator().AddColumn(&Match{}, "SeasonId")
    }

    if(!db.Migrator().HasColumn(&Match{}, "PlayerOfTheDay")){
        db.Migrator().AddColumn(&Match{}, "PlayerOfTheDay")
    }

    if(!db.Migrator().HasColumn(&Match{}, "DudOfTheDay")){
        db.Migrator().AddColumn(&Match{}, "DudOfTheDay")
    }

    if(!db.Migrator().HasColumn(&Fine{}, "Context")){
        db.Migrator().AddColumn(&Fine{}, "Context")
    }

    if(!db.Migrator().HasColumn(&Fine{}, "MatchId")){
        db.Migrator().AddColumn(&Fine{}, "MatchId")
    }

    if(!db.Migrator().HasColumn(&Fine{}, "Contest")){
        db.Migrator().AddColumn(&Fine{}, "Contest")
    }

    if(!db.Migrator().HasColumn(&Fine{}, "FineAt")){
        db.Migrator().AddColumn(&Fine{}, "FineAt")
    }

    if(!db.Migrator().HasColumn(&Player{}, "RoleDescription")){
        db.Migrator().AddColumn(&Player{}, "RoleDescription")
    }

    if(!db.Migrator().HasColumn(&Player{}, "Role")){
        db.Migrator().AddColumn(&Player{}, "Role")
    }


    
    if(!db.Migrator().HasColumn(&PresetFine{}, "Context")){
        db.Migrator().AddColumn(&PresetFine{}, "Context")
    }

    if(!db.Migrator().HasColumn(&PresetFine{}, "NotQuickFine")){
        db.Migrator().AddColumn(&PresetFine{}, "NotQuickFine")
    }

    if(!db.Migrator().HasColumn(&PresetFine{}, "DisplayOrder")){
        db.Migrator().AddColumn(&PresetFine{}, "DisplayOrder")
    }

    if(!db.Migrator().HasColumn(&PresetFine{}, "Icon")){
        db.Migrator().AddColumn(&PresetFine{}, "Icon")
    }

    if(!db.Migrator().HasColumn(&PresetFine{}, "IsKudos")){
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

// PlayerWithFines represents a player along with their fines
type PlayerWithFines struct {
    ID    uint
    Name        string
    TotalFineCount int
    TotalFines  int
    Role string
    RoleDescription string
    Fines []Fine
    PendingFines []Fine
    PendingFineSum int
    PendingTotalCount int
}

func GetPlayersWithFines(db *gorm.DB, playerIds []uint64) ([]PlayerWithFines, error) {
    var playersWithFines []PlayerWithFines

    var players []Player

    if len(playerIds) > 0 {
        db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
            return db.Order("fines.fine_at DESC")
        }).Where("id IN ?", playerIds).Find(&players)

    }else{
        db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
            return db.Order("fines.fine_at DESC")
        }).Find(&players).Order("players.name")
    }
    //.Where("active = true")
    

    // Construct the PlayerWithFines slice
    for _, player := range players {

        var approvedFines = []Fine{}
        var pendingFines = []Fine{}
        var fineSum = 0
        var pendingSum = 0
        for _, f := range player.Fines {
            if(f.Approved){
                approvedFines = append(approvedFines, f)
                fineSum = fineSum + int(f.Amount)
            } else {
                pendingFines = append(pendingFines, f)
                pendingSum = pendingSum + int(f.Amount)
            }
        } 
        pwf := PlayerWithFines{
            ID:    player.ID,
            Name:        player.Name,
            TotalFineCount: len(approvedFines),
            TotalFines:  fineSum,
            Role: player.Role,
            RoleDescription: player.RoleDescription,
            Fines: approvedFines,
            PendingFines: pendingFines,
            PendingFineSum: pendingSum,
            PendingTotalCount: len(pendingFines),
        }
        playersWithFines = append(playersWithFines, pwf)
    }

    sort.Slice(playersWithFines, func(i, j int) bool {
        return playersWithFines[i].TotalFines > playersWithFines[j].TotalFines
    })


    return playersWithFines, nil
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

func GetMatchMetaGeneral(db *gorm.DB, matchId uint) (*MatchMetaGeneral, error){
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


func UpdateFineContextByID(db *gorm.DB, fineID uint, matchId uint, context string, fineAt time.Time) error {
    // Create a map with the fields you want to update
    updates := map[string]interface{}{
        "Context": context,
        "MatchId": matchId,
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



type PresetFine struct {
    gorm.Model
    Reason string
    Amount float64
    Approved bool
    NotQuickFine bool
    Context string
    DisplayOrder int64
    Icon string
    IsKudos bool
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
            Match: matchedMatch,
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
    var match Match
    now := time.Now()

    // Find the nearest upcoming match where the start time is in the future
    result := db.Where("start_time > ?", now).Order("start_time ASC").First(&match)

    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // No upcoming match found is not considered an error; return nil
            return nil, nil
        }
        // Some other error occurred; return the error
        return nil, result.Error
    }

    return &match, nil
}




type Match struct {
    gorm.Model
    Location   string
    StartTime  *time.Time `json:"timestamp" gorm:"type:datetime"`
    Opponent   string
    Subtitle   string
    Events     []MatchEvent `gorm:"foreignKey:MatchId"`
    SeasonId    uint64
    PlayerOfTheDay    uint64
    DudOfTheDay uint64
}

type MatchEvent struct {
    gorm.Model
    MatchId  uint64
    EventName string
    EventType string // 'subbed-off' / 'subbed-on' / 'goal' / 'assist' / 'own-goal' / 'concded-goal'
    EventTime *time.Time `json:"timestamp" gorm:"type:datetime"`
    EventMinute int
    PlayerId uint
}

type PlayerState struct {
    PlayerId  uint
    PlayerName string
    TimePlayed  int // in minutes
    Goals       int
    Assists     int
}

type MatchState struct {
    PlayersOn     []PlayerState
    MatchID uint64
    ScoreAgainst int
    TrainingTotalNumbers int
    ScoreFor int
    MatchDuration int            // Duration of the match in minutes so far
}


// FetchLatestFines fetches a paginated list of the latest fines.
func GetMatches(db *gorm.DB, season uint64, page int, pageSize int) ([]Match, error) {
    var matches []Match
    offset := (page - 1) * pageSize
//.Where("season_id = ?", season)
    result := db.Order("start_time DESC").Offset(offset).Limit(pageSize).Find(&matches)
    if result.Error != nil {
        return nil, result.Error
    }

    return matches, nil
}


func GetMatch(db *gorm.DB, id uint64) (*Match, error){ 
    var match Match
    if err := db.Where("id = ?", id).First(&match).Error; err != nil {
        return nil, err
    }
    return &match, nil
}

func WrapMatchWithMeta(db *gorm.DB, match Match) (*MatchMetaGeneral, error){
    players, err := GetPlayers(db, 0, 999)
    if err != nil {
        return nil, err
    }

    var playerOfTheDay *Player
    var dudOfTheDay *Player
    var goalScorers []Player = []Player{};
    var opponentGoalCount uint = 0
    for i, p := range players {
        if(p.ID == uint(match.PlayerOfTheDay)){
            playerOfTheDay = &p
        }
        if(p.ID == uint(match.DudOfTheDay)){
            dudOfTheDay = &p
        }
        for _, e := range match.Events {
            if(e.EventType == "goal" && p.ID == e.PlayerId){
                goalScorers = append(goalScorers, p)
            }
            if(e.EventType == "conceded-goal" && i == 0){
                opponentGoalCount = opponentGoalCount + 1
            }
        }
    
    }

    
    // log.Printf("WrapMatchWithMeta - %d:%s %d:%s", match.PlayerOfTheDay, playerOfTheDay.Name, match.DudOfTheDay, dudOfTheDay.Name)
    return &MatchMetaGeneral{
        Match: match,
        GoalScorers: goalScorers,
        OpponentGoalCount: opponentGoalCount,
        Players: players,
        PlayerOfTheDay: playerOfTheDay,
        DudOfTheDay: dudOfTheDay,
    }, nil
}

func GetMatchWithEvents(db *gorm.DB, id uint) (*Match, error) {
    var match Match
    if err := db.Preload("Events").Where("id = ?", id).First(&match).Error; err != nil {
        return nil, err
    }
    

    return &match, nil 
}

func SaveMatch(db *gorm.DB, match *Match) (uint, error) {
    if err := db.Save(match).Error; err != nil {
        return 0, err // Return 0 as the ID in case of error
    }
    return match.ID, nil // Return the new ID which should now be populated
}

func GetMatchEvents(db *gorm.DB, id uint64) ([]MatchEvent, error) {
    var events []MatchEvent
    if err := db.Where("match_id = ?", id).Find(&events).Error; err != nil {
        return nil, err
    }
    return events, nil
}


func GetMatchEvent(db *gorm.DB, id uint64) (*MatchEvent, error) {
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

