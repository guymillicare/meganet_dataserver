package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
)

func CreateMarketOutcome(prematch *proto.Prematch) (*types.MarketOutcomeItem, error) {
	newMarketOutcome := &types.MarketOutcomeItem{}
	for _, odds := range prematch.Odds {
		oddsName := odds.Name
		if odds.Name == prematch.HomeTeam {
			oddsName = "Home"
		}
		if odds.Name == prematch.AwayTeam {
			oddsName = "Away"
		}
		marketOutcome, _ := MarketOutcomeFindByMarketAndOutcome(odds.MarketName, oddsName)
		if marketOutcome == nil {
			newOutcomeConstant := &types.OutcomeConstantItem{}
			newOutcomeConstant.ReferenceId = odds.MarketName + ":" + oddsName
			newOutcomeConstant.Name = oddsName
			if err := database.DB.Table("outcome_constants").Create(&newOutcomeConstant).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("MarketOutcomeCreate: %v", err)
			}

			newMarketOutcome.MarketRefId, _ = GetMarketReferenceId(odds.MarketName)
			newMarketOutcome.MarketDescription = odds.MarketName
			newMarketOutcome.OutcomeRefId = newOutcomeConstant.ReferenceId
			newMarketOutcome.OutcomeName = oddsName
			if err := database.DB.Table("market_outcomes").Create(&newMarketOutcome).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("MarketOutcomeCreate: %v", err)
			}
		}

	}
	return newMarketOutcome, nil
}

func MarketOutcomeFindByMarketAndOutcome(market string, outcome string) (*types.MarketOutcomeItem, error) {
	var marketOutcome *types.MarketOutcomeItem
	if err := database.DB.Table("market_outcomes").Where("market_description=? and outcome_name=?", market, outcome).First(&marketOutcome).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return marketOutcome, fmt.Errorf("FindById: %v", err)
	}
	return marketOutcome, nil
}
