package main

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSyncLineupMatchAllowsSameLineupAcrossMatches(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &Lineup{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	oldMatch := Match{TeamID: team.ID}
	newMatch := Match{TeamID: team.ID}
	if err := db.Create(&oldMatch).Error; err != nil {
		t.Fatalf("create old match: %v", err)
	}
	if err := db.Create(&newMatch).Error; err != nil {
		t.Fatalf("create new match: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	if err := db.Model(&Match{}).Where("id = ?", oldMatch.ID).Update("lineup_id", lineup.ID).Error; err != nil {
		t.Fatalf("seed old match lineup: %v", err)
	}

	if err := syncLineupMatch(db, lineup, newMatch.ID); err != nil {
		t.Fatalf("sync lineup match: %v", err)
	}

	if err := db.First(&oldMatch, oldMatch.ID).Error; err != nil {
		t.Fatalf("reload old match: %v", err)
	}
	if oldMatch.LineupID != lineup.ID {
		t.Fatalf("expected old match to keep lineup %d, got %d", lineup.ID, oldMatch.LineupID)
	}
	if err := db.First(&newMatch, newMatch.ID).Error; err != nil {
		t.Fatalf("reload new match: %v", err)
	}
	if newMatch.LineupID != lineup.ID {
		t.Fatalf("expected new match lineup %d, got %d", lineup.ID, newMatch.LineupID)
	}
}

func TestGetLineupHydratesMatchFromMatchLineupID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &LineupUser{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	lineup := Lineup{TeamID: team.ID, Name: "Line-up"}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	match := Match{TeamID: team.ID, LineupID: lineup.ID, Opponent: "Rivals"}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}

	loaded, err := getLineup(db, lineup.ID)
	if err != nil {
		t.Fatalf("get lineup: %v", err)
	}
	if loaded.Match.ID != match.ID {
		t.Fatalf("expected hydrated match model %d, got %d", match.ID, loaded.Match.ID)
	}
}

func TestStartLineupMatchSetsMatchLineupID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	match := Match{}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	lineup := Lineup{Status: lineupStatusSelected}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	match.LineupID = lineup.ID
	if err := db.Model(&Match{}).Where("id = ?", match.ID).Update("lineup_id", lineup.ID).Error; err != nil {
		t.Fatalf("attach lineup: %v", err)
	}
	match.LineupID = lineup.ID
	lineup.Match = match

	if err := startLineupMatch(db, &lineup, time.Now()); err != nil {
		t.Fatalf("start lineup match: %v", err)
	}

	if err := db.First(&match, match.ID).Error; err != nil {
		t.Fatalf("reload match: %v", err)
	}
	if match.LineupID != lineup.ID {
		t.Fatalf("expected match lineup %d, got %d", lineup.ID, match.LineupID)
	}
}

func TestMatchDayInitSnapshotFreezesLineupForMatchDisplay(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Match{}, &MatchEvent{}, &Lineup{}, &LineupPlayer{}, &Formation{}, &FormationPosition{}, &Player{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	player := Player{Name: "Original player"}
	if err := db.Create(&player).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}
	replacement := Player{Name: "Edited player"}
	if err := db.Create(&replacement).Error; err != nil {
		t.Fatalf("create replacement: %v", err)
	}
	formation := Formation{Name: "Original shape", Positions: []FormationPosition{{IndexNumber: 1, PositionName: "GK", X: 50, Y: 90}}}
	if err := db.Create(&formation).Error; err != nil {
		t.Fatalf("create formation: %v", err)
	}
	lineup := Lineup{FormationID: formation.ID, Name: "Shared line-up", Status: lineupStatusSelected}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create lineup: %v", err)
	}
	if err := db.Create(&LineupPlayer{LineupID: lineup.ID, IndexNumber: 1, SlotOrder: 0, PlayerID: player.ID}).Error; err != nil {
		t.Fatalf("create lineup player: %v", err)
	}
	match := Match{LineupID: lineup.ID}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("create match: %v", err)
	}
	lineup.Match = match
	if err := startLineupMatch(db, &lineup, time.Now()); err != nil {
		t.Fatalf("start lineup match: %v", err)
	}
	if err := db.Model(&LineupPlayer{}).Where("lineup_id = ? AND index_number = ?", lineup.ID, 1).Update("player_id", replacement.ID).Error; err != nil {
		t.Fatalf("edit source lineup player: %v", err)
	}

	loaded, err := getLineup(db, lineup.ID)
	if err != nil {
		t.Fatalf("load lineup: %v", err)
	}
	loaded.Match = match
	if err := db.Preload("Events").First(&loaded.Match, match.ID).Error; err != nil {
		t.Fatalf("load match events: %v", err)
	}
	applyMatchDayInitSnapshot(loaded)

	if len(loaded.Players) != 1 || loaded.Players[0].PlayerID != player.ID {
		t.Fatalf("expected snapshot player %d, got %#v", player.ID, loaded.Players)
	}
}
