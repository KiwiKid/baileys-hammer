package main

import (
	"log"
	"os"
	"sort"
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
}

// Fine represents a fine assigned to a player
type Fine struct {
    gorm.Model
    PlayerID uint
    Reason   string
    Amount   float64
    Approved bool
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
    if(!db.Migrator().HasColumn(&Match{}, "seasonId")){
        db.Migrator().AddColumn(&Match{}, "seasonId")
    }

    return db, nil
}

// PlayerWithFines represents a player along with their fines
type PlayerWithFines struct {
    PlayerID    uint 
    Name        string
    TotalFineCount int
    TotalFines  int
    Fines []Fine
    PendingFines []Fine
    PendingFineSum int
    PendingTotalCount int
}

func FetchPlayersWithFines(db *gorm.DB) ([]PlayerWithFines, error) {
    var playersWithFines []PlayerWithFines

    // Query all players
    var players []Player
    db.Preload("Fines", func(db *gorm.DB) *gorm.DB {
        return db.Order("fines.created_at DESC")
    }).Find(&players).Where("active = true")

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
            PlayerID:    player.ID,
            Name:        player.Name,
            TotalFineCount: len(approvedFines),
            TotalFines:  fineSum,
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




func GetPlayerByID(db *gorm.DB, playerID uint) (*Player, error) {
    var player Player
    result := db.First(&player, playerID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &player, nil
}


func SavePlayer(db *gorm.DB, player *Player) error {
    // Create or update player
    if err := db.Save(player).Error; err != nil {
        return err
    }
    return nil
}

// FetchLatestFines fetches a paginated list of the latest fines.
func FetchLatestFines(db *gorm.DB, page int, pageSize int) ([]Fine, error) {
    var fines []Fine
    offset := (page - 1) * pageSize

    // Query the latest fines with pagination
    result := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&fines)
    if result.Error != nil {
        return nil, result.Error
    }

    return fines, nil
}

// FetchLatestFines fetches a paginated list of the latest fines.
func FetchPlayers(db *gorm.DB, page int, pageSize int) ([]Player, error) {
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

func ApproveFine(db *gorm.DB, id uint, amount float64) error {
    // Find and update the fine's Approved field to true
    updates := map[string]interface{}{
        "approved": true,
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
    if err := db.Find(&presetFines).Where("approved = ?", includeUnapproved).Error; err != nil {
        return nil, err
    }
    return presetFines, nil
}


func GetPresetFine(db *gorm.DB, id uint64) (*PresetFine, error) {
    var presetFines PresetFine
    if err := db.Where("id = ?", id).First(&presetFines).Error; err != nil {
        return nil, err
    }
    return &presetFines, nil
}





type Match struct {
    gorm.Model
    Location   string
    StartTime  *time.Time `json:"timestamp" gorm:"type:datetime"`
    Opponent   string
    Subtitle   string
    Events     []MatchEvent `gorm:"foreignKey:MatchId"`
    SeasonId    uint64
}

type MatchEvent struct {
    gorm.Model
    MatchId  uint64
    EventName string
    EventType string // 'subbed-off' / 'subbed-on' / 'goal' / 'assist' / 'own-goal'
    EventTime *time.Time `json:"timestamp" gorm:"type:datetime"`
}


// FetchLatestFines fetches a paginated list of the latest fines.
func GetMatches(db *gorm.DB, season uint64, page int, pageSize int) ([]Match, error) {
    var matches []Match
    offset := (page - 1) * pageSize

    result := db.Order("created_at DESC").Where("season_id = ?", season).Offset(offset).Limit(pageSize).Find(&matches)
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

func SaveMatchEvent(db *gorm.DB, fine *MatchEvent) error {
    if err := db.Save(fine).Error; err != nil {
        return err
    }
    return nil
}

func DeleteMatchEvent(db *gorm.DB, fine *Match) error {
    if err := db.Delete(fine).Error; err != nil {
        return err
    }
    return nil
}

