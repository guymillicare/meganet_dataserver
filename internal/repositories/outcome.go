package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"strings"
	"time"

	"gorm.io/gorm/clause"
)

func CreateOutcome(prematch *proto.Prematch, sportEvent *types.SportEventItem) (*types.MarketOutcomeItem, error) {
	sport, _ := GetSportFromRedis(prematch.Sport)
	if sport == nil {
		return nil, nil
	}
	sportRefId := sport.ReferenceId
	// sportId := sport.Id
	newMarketOutcome := &types.MarketOutcomeItem{}
	const (
		homeReplacement = "Home"
		awayReplacement = "Away"
	)

	var newMarketOutcomes []*types.MarketOutcomeItem
	var newOutcomeConstants []*types.OutcomeConstantItem

	for _, odds := range prematch.Odds {
		oddsName := odds.Name
		oddsName = strings.Replace(oddsName, prematch.HomeTeam, homeReplacement, -1)
		oddsName = strings.Replace(oddsName, prematch.AwayTeam, awayReplacement, -1)
		marketConstant, _ := GetMarketConstantFromRedis(odds.MarketName)
		if marketConstant == nil {
			continue
		}

		newOutcomeConstant := &types.OutcomeConstantItem{
			ReferenceId: odds.MarketName + ":" + oddsName,
			Name:        oddsName,
		}
		newOutcomeConstants = append(newOutcomeConstants, newOutcomeConstant)

		newMarketOutcome := &types.MarketOutcomeItem{
			MarketRefId:       marketConstant.ReferenceId,
			MarketDescription: odds.MarketName,
			OutcomeRefId:      newOutcomeConstant.ReferenceId,
			OutcomeName:       oddsName,
			SportRefId:        sportRefId,
		}
		newMarketOutcomes = append(newMarketOutcomes, newMarketOutcome)

		// createOrUpdateMarketOutcome(newMarketOutcome, sportRefId, marketConstant, oddsName, odds.MarketName)
		// createOrUpdateSportMarketGroup(sportId, prematch.Sport, marketConstant, odds)
	}

	if len(newOutcomeConstants) > 0 {
		// if err := database.DB.Table("outcome_constants").Create(&newOutcomeConstants).Error; err != nil {
		// 	fmt.Printf("OutcomeConstantsCreate: %v\n", err)
		// 	return nil, err
		// }
		if err := database.DB.Table("outcome_constants").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "reference_id"}},
			DoNothing: true,
		}).Create(&newOutcomeConstants).Error; err != nil {
			return nil, err
		}
	}

	if len(newMarketOutcomes) > 0 {
		// if err := database.DB.Table("market_outcomes").Create(&newMarketOutcomes).Error; err != nil {
		// 	fmt.Printf("MarketOutcomeCreate: %v\n", err)
		// 	return nil, err
		// }
		if err := database.DB.Table("market_outcomes").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "market_ref_id"}, {Name: "outcome_ref_id"}},
			DoNothing: true,
		}).Create(&newMarketOutcomes).Error; err != nil {
			return nil, err
		}
	}

	for _, odds := range prematch.Odds {
		oddsName := odds.Name
		oddsName = strings.Replace(oddsName, prematch.HomeTeam, homeReplacement, -1)
		oddsName = strings.Replace(oddsName, prematch.AwayTeam, awayReplacement, -1)
		marketConstant, _ := GetMarketConstantFromRedis(odds.MarketName)
		if marketConstant == nil {
			continue
		}
		createOrUpdateOutcome(odds, sportEvent, marketConstant, oddsName)
	}

	return newMarketOutcome, nil
}

func createOrUpdateOutcome(odds *proto.Odds, sportEvent *types.SportEventItem, marketConstant *types.MarketConstantItem, oddsName string) error {
	// outcome, _ := OutcomeFind(sportEvent.Id, marketConstant.Id, oddsName)
	outcomeConstant, err := OutcomeConstantFind(odds.MarketName + ":" + oddsName)
	if err != nil {
		fmt.Printf("Error finding outcome constant: %v\n", err)
		return err
	}
	// if outcome == nil {
	outcome := &types.OutcomeItem{
		ReferenceId: outcomeConstant.ReferenceId,
		EventId:     sportEvent.Id,
		MarketId:    marketConstant.Id,
		Name:        oddsName,
		Odds:        odds.Price,
		Active:      true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	// 	if err := database.DB.Table("outcomes").Create(outcome).Error; err != nil {
	// 		return fmt.Errorf("OutcomeCreate: %v", err)
	// 	}
	// } else if outcome.Odds != odds.Price {
	// 	outcome.Odds = odds.Price

	// 	outcome.UpdatedAt = time.Now().UTC()
	// 	if err := database.DB.Table("outcomes").Save(outcome).Error; err != nil {
	// 		return fmt.Errorf("OutcomeSave: %v", err)
	// 	}
	// }
	if err := SaveOutcomeToRedis(outcome); err != nil {
		return err
	}
	return nil
}

func SaveOutcomeToRedis(outcome *types.OutcomeItem) error {
	ctx := context.Background()

	// Update the cache with the new outcome
	cacheKey := fmt.Sprintf("event:%d-outcomes", outcome.EventId)
	err := appendOutcomeToCache(ctx, cacheKey, outcome)
	if err != nil {
		fmt.Println("Error updating outcome cache:", err)
		return err
	}

	return nil
}

func appendOutcomeToCache(ctx context.Context, cacheKey string, outcome *types.OutcomeItem) error {
	// Attempt to fetch cached outcomes
	cachedOutcomesJSON, _ := database.RedisDB.Get(ctx, cacheKey).Bytes()
	if cachedOutcomesJSON == nil {
		// If cache miss, initialize empty array
		// return nil
		cachedOutcomesJSON = []byte("[]")
	}

	// Deserialize cached outcomes
	var cachedOutcomes []*types.OutcomeItem
	if err := json.Unmarshal(cachedOutcomesJSON, &cachedOutcomes); err != nil {
		return err
	}

	// Append new outcome to cached outcomes
	cachedOutcomes = append(cachedOutcomes, outcome)

	// Serialize updated outcomes
	updatedOutcomesJSON, err := json.Marshal(cachedOutcomes)
	if err != nil {
		return err
	}

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	// Update cache with the updated outcomes
	if err := database.RedisDB.Set(ctx, cacheKey, updatedOutcomesJSON, expiration).Err(); err != nil {
		return err
	}

	return nil
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
