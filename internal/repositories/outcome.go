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
)

func CreateOutcome(prematch *proto.Prematch, sportEvent *types.SportEventItem) (*types.MarketOutcomeItem, error) {
	sport, _ := GetSportFromRedis(prematch.Sport)
	if sport == nil {
		return nil, nil
	}
	sportRefId := sport.ReferenceId
	sportId := sport.Id
	newMarketOutcome := &types.MarketOutcomeItem{}
	const (
		homeReplacement = "Home"
		awayReplacement = "Away"
	)

	for _, odds := range prematch.Odds {
		oddsName := odds.Name
		oddsName = strings.Replace(oddsName, prematch.HomeTeam, homeReplacement, -1)
		oddsName = strings.Replace(oddsName, prematch.AwayTeam, awayReplacement, -1)
		marketConstant, _ := GetMarketConstantFromRedis(odds.MarketName)
		if marketConstant == nil {
			continue
		}

		createOrUpdateMarketOutcome(newMarketOutcome, sportRefId, marketConstant, oddsName, odds.MarketName)
		createOrUpdateSportMarketGroup(sportId, prematch.Sport, marketConstant, odds)
		createOrUpdateOutcome(odds, sportEvent, marketConstant, oddsName)
	}

	return newMarketOutcome, nil
}

func createOrUpdateOutcome(odds *proto.Odds, sportEvent *types.SportEventItem, marketConstant *types.MarketConstantItem, oddsName string) error {
	outcome, _ := OutcomeFind(sportEvent.Id, marketConstant.Id, oddsName)
	outcomeConstant, _ := OutcomeConstantFind(odds.MarketName + ":" + oddsName)
	if outcome == nil {
		outcome = &types.OutcomeItem{
			ReferenceId: outcomeConstant.ReferenceId,
			EventId:     sportEvent.Id,
			MarketId:    marketConstant.Id,
			Name:        oddsName,
			Odds:        odds.Price,
			Active:      true,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		if err := database.DB.Table("outcomes").Create(outcome).Error; err != nil {
			return fmt.Errorf("OutcomeCreate: %v", err)
		}
	} else if outcome.Odds != odds.Price {
		outcome.Odds = odds.Price

		outcome.UpdatedAt = time.Now().UTC()
		if err := database.DB.Table("outcomes").Save(outcome).Error; err != nil {
			return fmt.Errorf("OutcomeSave: %v", err)
		}
	}
	if err := SaveOutcomeToRedis(outcome); err != nil {
		return err
	}
	return nil
}

func SaveOutcomeToRedis(outcome *types.OutcomeItem) error {
	outcomeJSON, err := json.Marshal(outcome)
	if err != nil {
		fmt.Println("Error marshaling OutcomeItem:", err)
		return err
	}

	ctx := context.Background()
	key := fmt.Sprintf("event:%d-outcome:%s", outcome.EventId, outcome.ReferenceId)

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	// Save the outcome to Redis
	err = database.RedisDB.Set(ctx, key, outcomeJSON, expiration).Err()
	if err != nil {
		fmt.Println("Error saving OutcomeItem to Redis:", err)
		return err
	}

	// Update the cache with the new outcome
	cacheKey := fmt.Sprintf("event:%d-outcomes", outcome.EventId)
	err = appendOutcomeToCache(ctx, cacheKey, outcome)
	if err != nil {
		fmt.Println("Error updating outcome cache:", err)
		// Rollback the outcome from Redis
		rollbackErr := database.RedisDB.Del(ctx, key).Err()
		if rollbackErr != nil {
			fmt.Println("Error rolling back OutcomeItem from Redis:", rollbackErr)
		}
		return err
	}

	return nil
}

func appendOutcomeToCache(ctx context.Context, cacheKey string, outcome *types.OutcomeItem) error {
	// Attempt to fetch cached outcomes
	cachedOutcomesJSON, err := database.RedisDB.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return nil
	}
	if cachedOutcomesJSON == nil {
		// If cache miss, initialize empty array
		return nil
		// cachedOutcomesJSON = []byte("[]")
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
