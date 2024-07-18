package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"gorm.io/gorm/clause"
)

func SportMarketGroupFindBy(sportId int32, marketId int32) (*types.SportMarketGroupItem, error) {
	var sportMarketGroup *types.SportMarketGroupItem
	if err := database.DB.Table("sport_market_groups").Where("sport_id =? and market_id=?", sportId, marketId).First(&sportMarketGroup).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return sportMarketGroup, fmt.Errorf("SportMarketGroupFindBy: %v", err)
	}
	return sportMarketGroup, nil
}

func createOrUpdateSportMarketGroup(sportId int32, sportName string, marketConstant *types.MarketConstantItem, odds *proto.Odds) error {
	sportMarketGroup, _ := SportMarketGroupFindBy(sportId, marketConstant.Id)
	if sportMarketGroup == nil {
		newSportMarketGroup := &types.SportMarketGroupItem{
			SportId:    sportId,
			SportName:  sportName,
			MarketId:   marketConstant.Id,
			MarketName: odds.MarketName,
		}
		if err := database.DB.Table("sport_market_groups").Create(newSportMarketGroup).Error; err != nil {
			return fmt.Errorf("SportMarketGroupCreate: %v", err)
		}
	}
	return nil
}

func UpdateSportMarketGroup(sportMarketGroups []*types.SportMarketGroupItem) {
	if err := database.DB.Table("sport_market_groups").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "market_id"}},
		DoNothing: true,
		// DoUpdates: clause.AssignmentColumns([]string{"group_name", "market_name"}),
	}).Create(&sportMarketGroups).Error; err != nil {
		fmt.Printf("UpdateSportMarketGroup: %v\n", err)
	}
}
