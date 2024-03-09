package main

import (
	"sort"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Player represents a player in the social football team
type Player struct {
    gorm.Model
    Name   string `gorm:"unique" form:"name"`
    Fines  []Fine
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
	
    db, err := gorm.Open(sqlite.Open("production.sqlite3.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Migrate the schema
    err = db.AutoMigrate(&Player{}, &Fine{}, &PresetFine{})
    if err != nil {
        return nil, err
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
    }).Find(&players)

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

// SaveFine adds a new fine or updates an existing fine in the database
func SaveFine(db *gorm.DB, fine *Fine) error {
    // Create or update fine
    if err := db.Save(fine).Error; err != nil {
        return err
    }
    return nil
}

func ApproveFine(db *gorm.DB, id uint) error {
    // Find and update the fine's Approved field to true
    result := db.Model(&Fine{}).Where("id = ?", id).Update("approved", true)

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