package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
	"sync"

	"github.com/go-redis/redis"
)

var outcomeConstantCache sync.Map

func OutcomeConstantsFindAll() ([]*types.OutcomeConstantItem, error) {
	var outcomeConstants []*types.OutcomeConstantItem
	if err := database.DB.Table("outcome_constants").Where("outcome_constants.data_feed='huge_data'").Find(&outcomeConstants).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcomeConstants, fmt.Errorf("OutcomeConstantsFindAll: %v", err)
	}
	return outcomeConstants, nil
}

func OutcomeConstantFind(reference_id string) (*types.OutcomeConstantItem, error) {
	// Check if the item is in the cache
	if value, ok := outcomeConstantCache.Load(reference_id); ok {
		return value.(*types.OutcomeConstantItem), nil
	}

	// If not in the cache, query the database
	var outcomeConstant *types.OutcomeConstantItem
	if err := database.DB.Table("outcome_constants").Where("reference_id =?", reference_id).First(&outcomeConstant).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return outcomeConstant, fmt.Errorf("OutcomeConstantFind: %v", err)
	}

	// Store the result in the cache
	outcomeConstantCache.Store(reference_id, outcomeConstant)
	return outcomeConstant, nil
}

func saveOutcomeConstantToRedis(ctx context.Context, outcomeConstant *types.OutcomeConstantItem) error {
	outcomeConstantJSON, err := json.Marshal(outcomeConstant)
	if err != nil {
		return fmt.Errorf("saveOutcomeConstantToRedis: error marshaling outcomeConstant: %v", err)
	}

	key := fmt.Sprintf("outcomeConstant:%s", outcomeConstant.ReferenceId)
	err = database.RedisDB.Set(ctx, key, outcomeConstantJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("saveOutcomeConstantToRedis: error saving outcomeConstant to Redis: %v", err)
	}

	return nil
}

func OutcomeConstantsPreload() {
	outcomeConstants, _ := OutcomeConstantsFindAll()
	ctx := context.Background()
	for _, outcomeConstant := range outcomeConstants {
		if err := saveOutcomeConstantToRedis(ctx, outcomeConstant); err != nil {
			fmt.Printf("saveOutcomeConstantToRedis: error saving outcomeConstant to Redis: %v\n", err)
			continue
		}
	}
}

func GetOutcomeConstantFromRedis(refId string) (*types.OutcomeConstantItem, error) {
	// Construct the key for the outcomeConstant item
	ctx := context.Background()
	key := fmt.Sprintf("outcomeConstant:%s", refId)

	// Retrieve the outcomeConstant item from Redis
	outcomeConstantJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// outcomeConstant item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetOutcomeConstantFromRedis: error fetching outcomeConstant item from Redis: %v", err)
	}

	// Deserialize the outcomeConstant item from JSON
	var outcomeConstant types.OutcomeConstantItem
	if err := json.Unmarshal(outcomeConstantJSON, &outcomeConstant); err != nil {
		return nil, fmt.Errorf("GetOutcomeConstantFromRedis: error unmarshaling outcomeConstant item: %v", err)
	}

	return &outcomeConstant, nil
}
