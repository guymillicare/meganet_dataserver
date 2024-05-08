package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func TournamentsFindAll() ([]*types.TournamentItem, error) {
	var tournaments []*types.TournamentItem
	if err := database.DB.Table("tournaments").Find(&tournaments).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return tournaments, fmt.Errorf("TournamentsFindAll: %v", err)
	}
	return tournaments, nil
}

func TournamentsPreload() {
	database.TOURNAMENTS, _ = TournamentsFindAll()
}

func GetTournamentId(name string) (string, error) {
	for _, tournament := range database.TOURNAMENTS {
		if tournament.Name == name {
			return tournament.ReferenceId, nil
		}
	}
	return "0", fmt.Errorf("Tournament not found")
}
