package repositories

import (
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"strings"
	"time"
)

func CreateOutcome(prematch *proto.Prematch, sportEvent *types.SportEventItem) (*types.MarketOutcomeItem, error) {
	sportRefId, _ := GetSportRefId(prematch.Sport)
	sportId, _ := GetSportId(prematch.Sport)
	newMarketOutcome := &types.MarketOutcomeItem{}
	for _, odds := range prematch.Odds {
		oddsName := odds.Name
		oddsName = strings.Replace(oddsName, prematch.HomeTeam, "Home", -1)
		oddsName = strings.Replace(oddsName, prematch.AwayTeam, "Away", -1)
		marketConstant, _ := GetMarketConstant(odds.MarketName)

		marketOutcome, _ := MarketOutcomeFindByMarketAndOutcome(odds.MarketName, oddsName)
		if marketOutcome == nil {
			newMarketOutcome = &types.MarketOutcomeItem{}
			newOutcomeConstant := &types.OutcomeConstantItem{}
			newOutcomeConstant.ReferenceId = odds.MarketName + ":" + oddsName
			newOutcomeConstant.Name = oddsName
			if err := database.DB.Table("outcome_constants").Create(&newOutcomeConstant).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("OutcomeConstantsCreate: %v", err)
			}

			newMarketOutcome.MarketRefId = marketConstant.ReferenceId
			newMarketOutcome.MarketDescription = odds.MarketName
			newMarketOutcome.OutcomeRefId = newOutcomeConstant.ReferenceId
			newMarketOutcome.OutcomeName = oddsName
			newMarketOutcome.SportRefId = sportRefId
			if err := database.DB.Table("market_outcomes").Create(&newMarketOutcome).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("MarketOutcomeCreate: %v", err)
			}
		}

		sportMarketGroup, _ := SportMarketGroupFindBy(sportId, marketConstant.Id)
		if sportMarketGroup == nil {
			newSportMarketGroup := &types.SportMarketGroupItem{}
			newSportMarketGroup.SportId = sportId
			newSportMarketGroup.SportName = prematch.Sport
			newSportMarketGroup.MarketId = marketConstant.Id
			newSportMarketGroup.MarketName = odds.MarketName
			if err := database.DB.Table("sport_market_groups").Create(&newSportMarketGroup).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("SportMarketGroupCreate: %v", err)
			}
		}

		outcome, _ := OutcomeFind(sportEvent.Id, marketConstant.Id, oddsName)
		outcomeConstant, _ := OutcomeConstantFind(odds.MarketName + ":" + oddsName)
		if outcome == nil {
			outcome = &types.OutcomeItem{}
			outcome.ReferenceId = outcomeConstant.ReferenceId
			outcome.EventId = sportEvent.Id
			outcome.MarketId = marketConstant.Id
			outcome.Name = oddsName
			outcome.Odds = odds.Price
			outcome.CreatedAt = time.Now().UTC()
			outcome.UpdatedAt = time.Now().UTC()
			if err := database.DB.Table("outcomes").Create(&outcome).Error; err != nil {
				return newMarketOutcome, fmt.Errorf("OutcomeCreate: %v", err)
			}
		} else {
			if outcome.Odds != odds.Price {
				outcome.Odds = odds.Price
				outcome.UpdatedAt = time.Now().UTC()
				if err := database.DB.Table("outcomes").Save(&outcome).Error; err != nil {
					return newMarketOutcome, fmt.Errorf("OutcomeSave: %v", err)
				}
			}
		}
	}
	return newMarketOutcome, nil
}

func OutcomeFind(eventId int32, marketId int32, name string) (*types.OutcomeItem, error) {
	var outcome *types.OutcomeItem
	if err := database.DB.Table("outcomes").Where("event_id =? and market_id=? and name=?", eventId, marketId, name).First(&outcome).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcome, fmt.Errorf("OutcomeFind: %v", err)
	}
	return outcome, nil
}

func OutcomeConstantFind(reference_id string) (*types.OutcomeConstantItem, error) {
	var outcomeConstant *types.OutcomeConstantItem
	if err := database.DB.Table("outcome_constants").Where("reference_id =?", reference_id).First(&outcomeConstant).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcomeConstant, fmt.Errorf("OutcomeConstantFind: %v", err)
	}
	return outcomeConstant, nil
}
