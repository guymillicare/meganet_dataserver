package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
)

// FetchAllGames retrieves all games from the database
func CreateCompetitor(prematch *proto.Prematch) (*types.CompetitorItem, error) {
	competitor := &types.CompetitorItem{}
	homeCompetitor, _ := CompetitorFindByName(prematch.HomeTeam)
	if homeCompetitor == nil {
		competitor.Name = prematch.HomeTeam
		competitor.SportId, _ = GetSportRefId(prematch.Sport)
		if err := database.DB.Table("competitors").Create(&competitor).Error; err != nil {
			return competitor, fmt.Errorf("CompetitorCreate: %v", err)
		}
	}
	competitor = &types.CompetitorItem{}
	awayCompetitor, _ := CompetitorFindByName(prematch.HomeTeam)
	if awayCompetitor == nil {
		competitor.Name = prematch.AwayTeam
		competitor.SportId, _ = GetSportRefId(prematch.Sport)
		if err := database.DB.Table("competitors").Create(&competitor).Error; err != nil {
			return competitor, fmt.Errorf("CompetitorCreate: %v", err)
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
