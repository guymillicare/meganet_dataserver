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

func TournamentsFindAll() ([]*types.TournamentItem, error) {
	var tournaments []*types.TournamentItem
	if err := database.DB.Table("tournaments").Find(&tournaments).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return tournaments, fmt.Errorf("TournamentsFindAll: %v", err)
	}
	return tournaments, nil
}

func saveTournamentToRedis(ctx context.Context, tournament *types.TournamentItem) error {
	tournamentJSON, err := json.Marshal(tournament)
	if err != nil {
		return fmt.Errorf("saveTournamentToRedis: error marshaling tournament: %v", err)
	}
	key := fmt.Sprintf("tournament:%s", tournament.Name)

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	err = database.RedisDB.Set(ctx, key, tournamentJSON, expiration).Err()
	if err != nil {
		return fmt.Errorf("saveTournamentToRedis: error saving tournament to Redis: %v", err)
	}
	return nil
}

func TournamentsPreload() {
	tournaments, _ := TournamentsFindAll()
	ctx := context.Background()
	for _, tournament := range tournaments {
		if err := saveTournamentToRedis(ctx, tournament); err != nil {
			fmt.Printf("saveTournamentToRedis: error saving tournament to Redis: %v\n", err)
			continue
		}
	}
}

func GetTournamentFromRedis(name string) (*types.TournamentItem, error) {
	// Construct the key for the tournament item
	ctx := context.Background()
	key := fmt.Sprintf("tournament:%s", name)

	// Retrieve the tournament item from Redis
	tournamentJSON, err := database.RedisDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// tournament item not found in Redis
			return nil, nil
		}
		return nil, fmt.Errorf("GetTournamentFromRedis: error fetching tournament item from Redis: %v", err)
	}

	// Deserialize the sport item from JSON
	var tournament types.TournamentItem
	if err := json.Unmarshal(tournamentJSON, &tournament); err != nil {
		return nil, fmt.Errorf("GetTournamentFromRedis: error unmarshaling tournament item: %v", err)
	}

	return &tournament, nil
}
