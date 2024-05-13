package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"

	"github.com/go-redis/redis"
)

func MarketConstantsFindAll() ([]*types.MarketConstantItem, error) {
	var marketConstants []*types.MarketConstantItem
	if err := database.DB.Table("market_constants").Find(&marketConstants).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return marketConstants, fmt.Errorf("MarketConstantsFindAll: %v", err)
	}
	return marketConstants, nil
}

func saveMarketConstantToRedis(ctx context.Context, marketConstant *types.MarketConstantItem) error {
	marketConstantJSON, err := json.Marshal(marketConstant)
	if err != nil {
		return fmt.Errorf("saveMarketConstantToRedis: error marshaling marketConstant: %v", err)
	}

	key := fmt.Sprintf("marketConstant:%s", marketConstant.Description)

	err = database.RedisDB.Set(ctx, key, marketConstantJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("saveMarketConstantToRedis: error saving marketConstant to Redis: %v", err)
	}

	return nil
}

func MarketConstantsPreload() {
	marketConstants, _ := MarketConstantsFindAll()
	ctx := context.Background()
	for _, marketConstant := range marketConstants {
		if err := saveMarketConstantToRedis(ctx, marketConstant); err != nil {
			fmt.Printf("saveMarketConstantToRedis: error saving marketConstant to Redis: %v\n", err)
			continue
		}
	}
}

func GetMarketConstantFromRedis(description string) (*types.MarketConstantItem, error) {
	// Construct the key for the marketConstant item
	ctx := context.Background()
	key := fmt.Sprintf("marketConstant:%s", description)

	// Retrieve the marketConstant item from Redis
	marketConstantJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// marketConstant item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetMarketConstantFromRedis: error fetching marketConstant item from Redis: %v", err)
	}

	// Deserialize the marketConstant item from JSON
	var marketConstant types.MarketConstantItem
	if err := json.Unmarshal(marketConstantJSON, &marketConstant); err != nil {
		return nil, fmt.Errorf("GetMarketConstantFromRedis: error unmarshaling marketConstant item: %v", err)
	}

	return &marketConstant, nil
}
