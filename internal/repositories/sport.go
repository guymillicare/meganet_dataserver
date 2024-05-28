package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
	"time"

	"github.com/go-redis/redis"
)

func SportsFindAll() ([]*types.SportItem, error) {
	var sports []*types.SportItem
	if err := database.DB.Table("sports").Find(&sports).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return sports, fmt.Errorf("SportsFindAll: %v", err)
	}
	// Save serialized sports data to Redis

	return sports, nil

}

func saveSportToRedis(ctx context.Context, sport *types.SportItem) error {
	sportJSON, err := json.Marshal(sport)
	if err != nil {
		return fmt.Errorf("saveSportToRedis: error marshaling sport: %v", err)
	}

	key := fmt.Sprintf("sport:%s", sport.Slug)

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	err = database.RedisDB.Set(ctx, key, sportJSON, expiration).Err()
	if err != nil {
		return fmt.Errorf("saveSportToRedis: error saving sport to Redis: %v", err)
	}

	return nil
}

func SportsPreload() {
	sports, _ := SportsFindAll()
	ctx := context.Background()
	for _, sport := range sports {
		if err := saveSportToRedis(ctx, sport); err != nil {
			fmt.Printf("saveSportToRedis: error saving sport to Redis: %v\n", err)
			continue
		}
	}
}

func GetSportFromRedis(slug string) (*types.SportItem, error) {
	// Construct the key for the sport item
	ctx := context.Background()
	key := fmt.Sprintf("sport:%s", slug)

	// Retrieve the sport item from Redis
	sportJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Sport item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetSportFromRedis: error fetching sport item from Redis: %v", err)
	}

	// Deserialize the sport item from JSON
	var sport types.SportItem
	if err := json.Unmarshal(sportJSON, &sport); err != nil {
		return nil, fmt.Errorf("GetSportFromRedis: error unmarshaling sport item: %v", err)
	}

	return &sport, nil
}
