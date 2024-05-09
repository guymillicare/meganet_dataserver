package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func GetMarketConstant(description string) (*types.MarketConstantItem, error) {
	var marketConstant *types.MarketConstantItem
	if err := database.DB.Table("market_constants").Where("description =?", description).First(&marketConstant).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return marketConstant, fmt.Errorf("FindById: %v", err)
	}
	return marketConstant, nil
}
