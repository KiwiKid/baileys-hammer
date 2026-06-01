package main

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func GetUnavailablePlayerIDsForMatch(db *gorm.DB, teamID uint, matchID uint) (map[uint]bool, error) {
	unavailable := map[uint]bool{}
	if matchID == 0 {
		return unavailable, nil
	}
	var records []PlayerMatchUnavailability
	query := db.Where("match_id = ?", matchID)
	if teamID > 0 {
		query = query.Where("team_id = ?", teamID)
	}
	if err := query.Find(&records).Error; err != nil {
		return unavailable, err
	}
	for _, record := range records {
		if record.PlayerID > 0 {
			unavailable[record.PlayerID] = true
		}
	}
	return unavailable, nil
}

func SaveUnavailablePlayersForMatch(db *gorm.DB, teamID uint, matchID uint, playerIDs []uint) error {
	if matchID == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		query := tx.Where("match_id = ?", matchID)
		if teamID > 0 {
			query = query.Where("team_id = ?", teamID)
		}
		if err := query.Unscoped().Delete(&PlayerMatchUnavailability{}).Error; err != nil {
			return err
		}
		seen := map[uint]bool{}
		for _, playerID := range playerIDs {
			if playerID == 0 || seen[playerID] {
				continue
			}
			if err := tx.Create(&PlayerMatchUnavailability{
				TeamID:   teamID,
				MatchID:  matchID,
				PlayerID: playerID,
			}).Error; err != nil {
				return err
			}
			seen[playerID] = true
		}
		return nil
	})
}

func parseUnavailablePlayerIDs(r *http.Request) ([]uint, error) {
	values := r.Form["unavailablePlayerIds"]
	values = append(values, r.Form["unavailablePlayerIds[]"]...)
	playerIDs := []uint{}
	for _, value := range values {
		if value == "" {
			continue
		}
		playerID, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		playerIDs = append(playerIDs, uint(playerID))
	}
	return playerIDs, nil
}

func playerAvailabilityCounts(players []Player, unavailablePlayerIDs map[uint]bool) (int, int) {
	total := 0
	unavailable := 0
	for _, player := range players {
		if !player.Active {
			continue
		}
		total++
		if unavailablePlayerIDs[player.ID] {
			unavailable++
		}
	}
	return total - unavailable, total
}

func playerUnavailableForMatch(unavailablePlayerIDs map[uint]bool, playerID uint) bool {
	return unavailablePlayerIDs != nil && unavailablePlayerIDs[playerID]
}

func lineupPlayerOptionVisible(unavailablePlayerIDs map[uint]bool, selectedID uint, playerID uint) bool {
	return selectedID == playerID || !playerUnavailableForMatch(unavailablePlayerIDs, playerID)
}
