package repositories

import (
	"fmt"
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

func createOrUpdateMarketOutcome(newMarketOutcome *types.MarketOutcomeItem, sportRefId string, marketConstant *types.MarketConstantItem, oddsName, marketName string) error {
	marketOutcome, _ := MarketOutcomeFindByMarketAndOutcome(marketName, oddsName)
	if marketOutcome == nil {
		newOutcomeConstant := &types.OutcomeConstantItem{
			ReferenceId: marketName + ":" + oddsName,
			Name:        oddsName,
		}
		if err := database.DB.Table("outcome_constants").Create(newOutcomeConstant).Error; err != nil {
			return fmt.Errorf("OutcomeConstantsCreate: %v", err)
		}
		newMarketOutcome = &types.MarketOutcomeItem{
			MarketRefId:       marketConstant.ReferenceId,
			MarketDescription: marketName,
			OutcomeRefId:      newOutcomeConstant.ReferenceId,
			OutcomeName:       oddsName,
			SportRefId:        sportRefId,
		}
		if err := database.DB.Table("market_outcomes").Create(newMarketOutcome).Error; err != nil {
			return fmt.Errorf("MarketOutcomeCreate: %v", err)
		}
	}
	return nil
}
