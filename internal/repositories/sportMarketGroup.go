package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
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
