package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"

	"github.com/go-redis/redis"
)

func MarketGroupsFindAll() ([]*types.MarketGroupItem, error) {
	var marketGroups []*types.MarketGroupItem
	if err := database.DB.Table("market_groups").Find(&marketGroups).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return marketGroups, fmt.Errorf("MarketGroupsFindAll: %v", err)
	}
	return marketGroups, nil
}

func saveMarketGroupToRedis(ctx context.Context, marketGroup *types.MarketGroupItem) error {
	marketGroupJSON, err := json.Marshal(marketGroup)
	if err != nil {
		return fmt.Errorf("saveMarketGroupToRedis: error marshaling marketGroup: %v", err)
	}

	key := fmt.Sprintf("marketGroup:%d", marketGroup.Id)
	err = database.RedisDB.Set(ctx, key, marketGroupJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("saveMarketGroupToRedis: error saving marketGroup to Redis: %v", err)
	}

	return nil
}

func MarketGroupsPreload() {
	marketGroups, _ := MarketGroupsFindAll()
	ctx := context.Background()
	for _, marketGroup := range marketGroups {
		if err := saveMarketGroupToRedis(ctx, marketGroup); err != nil {
			fmt.Printf("saveMarketGroupToRedis: error saving marketGroup to Redis: %v\n", err)
			continue
		}
	}
}

func GetMarketGroupFromRedis(Id int) (*types.MarketGroupItem, error) {
	// Construct the key for the marketGroup item
	ctx := context.Background()
	key := fmt.Sprintf("marketGroup:%d", Id)

	// Retrieve the marketGroup item from Redis
	marketGroupJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// marketGroup item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetMarketGroupFromRedis: error fetching marketGroup item from Redis: %v", err)
	}

	// Deserialize the marketGroup item from JSON
	var marketGroup types.MarketGroupItem
	if err := json.Unmarshal(marketGroupJSON, &marketGroup); err != nil {
		return nil, fmt.Errorf("GetMarketGroupFromRedis: error unmarshaling marketGroup item: %v", err)
	}

	return &marketGroup, nil
}
