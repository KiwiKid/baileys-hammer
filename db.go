package main

import (
	"sort"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Player represents a player in the social football team
type Player struct {
    gorm.Model
    Name   string `form:"name"`
    Number int	 	`form:"number"`
    Fines  []Fine
}

// Fine represents a fine assigned to a player
type Fine struct {
    gorm.Model
    PlayerID uint
    Reason   string
    Amount   float64
    Status   string // "proposed", "approved", "paid"
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
    Number      int
    TotalFineCount int
    TotalFines  int
    Fines []Fine
}

func FetchPlayersWithFines(db *gorm.DB) ([]PlayerWithFines, error) {
    var playersWithFines []PlayerWithFines

    // Query all players
    var players []Player
    db.Preload("Fines").Find(&players)


    // Construct the PlayerWithFines slice
    for _, player := range players {

        var fineSum = 0
        for _, f := range player.Fines {
            fineSum = fineSum + int(f.Amount)
        }
    

        pwf := PlayerWithFines{
            PlayerID:    player.ID,
            Name:        player.Name,
            Number:      player.Number,
            TotalFineCount: len(player.Fines),
            TotalFines:  fineSum,
            Fines: player.Fines,
        }
        playersWithFines = append(playersWithFines, pwf)
    }

    sort.Slice(playersWithFines, func(i, j int) bool {
        return playersWithFines[i].TotalFines > playersWithFines[j].TotalFines
    })


    return playersWithFines, nil
}


func SavePlayer(db *gorm.DB, player *Player) error {
    // Create or update player
    if err := db.Save(player).Error; err != nil {
        return err
    }
    return nil
}

// SaveFine adds a new fine or updates an existing fine in the database
func SaveFine(db *gorm.DB, fine *Fine) error {
    // Create or update fine
    if err := db.Save(fine).Error; err != nil {
        return err
    }
    return nil
}


type PresetFine struct {
    gorm.Model
    Reason string
    Amount float64
}


func SavePresetFine(db *gorm.DB, presetFine *PresetFine) error {
    if err := db.Save(presetFine).Error; err != nil {
        return err
    }
    return nil
}

func GetPresetFines(db *gorm.DB) ([]PresetFine, error) {
    var presetFines []PresetFine
    if err := db.Find(&presetFines).Error; err != nil {
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