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

func CountriesFindAll() ([]*types.CountryItem, error) {
	var countries []*types.CountryItem
	if err := database.DB.Table("countries").Find(&countries).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return countries, fmt.Errorf("CountriesFindAll: %v", err)
	}
	return countries, nil
}

func saveCountryToRedis(ctx context.Context, country *types.CountryItem) error {
	countryJSON, err := json.Marshal(country)
	if err != nil {
		return fmt.Errorf("saveCountryToRedis: error marshaling country: %v", err)
	}

	key := fmt.Sprintf("country:%s", country.Name)

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	err = database.RedisDB.Set(ctx, key, countryJSON, expiration).Err()
	if err != nil {
		return fmt.Errorf("saveCountryToRedis: error saving country to Redis: %v", err)
	}

	return nil
}

func CountriesPreload() {
	countries, _ := CountriesFindAll()
	ctx := context.Background()
	for _, country := range countries {
		if err := saveCountryToRedis(ctx, country); err != nil {
			fmt.Printf("saveCountryToRedis: error saving country to Redis: %v\n", err)
			continue
		}
	}
}

func GetCountryFromRedis(name string) (*types.CountryItem, error) {
	// Construct the key for the country item
	ctx := context.Background()
	key := fmt.Sprintf("country:%s", name)

	// Retrieve the country item from Redis
	countryJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Country item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetCountryFromRedis: error fetching country item from Redis: %v", err)
	}

	// Deserialize the sport item from JSON
	var country types.CountryItem
	if err := json.Unmarshal(countryJSON, &country); err != nil {
		return nil, fmt.Errorf("GetCountryFromRedis: error unmarshaling country item: %v", err)
	}

	return &country, nil
}
