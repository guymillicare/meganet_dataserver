package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"gorm.io/gorm/clause"
)

// CreateCompetitorsBatch inserts competitors into the database in a single batch.
func CreateCompetitorsBatch(competitors []*types.CompetitorItem) error {
	if len(competitors) > 0 {
		err := database.DB.Table("competitors").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&competitors).Error

		if err != nil {
			return fmt.Errorf("CompetitorCreateBatch: %v", err)
		}
	}

	return nil
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

// PrepareCompetitors collects competitors from a prematch
func PrepareCompetitors(prematch *proto.Prematch) []*types.CompetitorItem {
	sport, err := GetSportFromRedis(prematch.Sport)
	if err != nil {
		fmt.Println("GetSportFromRedis:", err)
		return nil
	}

	teams := []string{prematch.HomeTeam, prematch.AwayTeam}
	competitors := make([]*types.CompetitorItem, len(teams))

	for i, team := range teams {
		competitors[i] = &types.CompetitorItem{
			Name:    team,
			SportId: sport.ReferenceId,
		}
	}

	return competitors
}
