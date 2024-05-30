package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"gorm.io/gorm/clause"
)

// FetchAllGames retrieves all games from the database
func CreateCompetitor(prematch *proto.Prematch) (*types.CompetitorItem, error) {
	// Fetch sport details once
	sport, err := GetSportFromRedis(prematch.Sport)
	if err != nil {
		return nil, fmt.Errorf("GetSportFromRedis: %v", err)
	}

	teams := []string{prematch.HomeTeam, prematch.AwayTeam}
	competitor := &types.CompetitorItem{}
	for _, team := range teams {
		competitor = &types.CompetitorItem{
			Name:    team,
			SportId: sport.ReferenceId,
		}

		err := database.DB.Table("competitors").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&competitor).Error

		if err != nil {
			return nil, fmt.Errorf("CompetitorCreate: %v", err)
		}
	}

	return competitor, nil
}

func CompetitorFindByName(name string) (*types.CompetitorItem, error) {
	var competitor *types.CompetitorItem
	if err := database.DB.Table("competitors").Where("name =?", name).First(&competitor).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return competitor, fmt.Errorf("CompetitorFindOne: %v", err)
	}
	return competitor, nil
}
