package repositories

import (
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
)

func MarketOutcomeFindByMarketAndOutcome(market string, outcome string) (*types.MarketOutcomeItem, error) {
	var marketOutcome *types.MarketOutcomeItem
	if err := database.DB.Table("market_outcomes").Where("market_description=? and outcome_name=?", market, outcome).First(&marketOutcome).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return marketOutcome, nil
	}
	return marketOutcome, nil
}
