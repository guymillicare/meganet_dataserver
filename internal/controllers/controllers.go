package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/types"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-redis/redis"
)

func GetOutcomesByEventId(w http.ResponseWriter, r *http.Request) {
	eventId := chi.URLParam(r, "eventId")

	ctx := context.Background()

	// Attempt to fetch outcomes from cache
	outcomes, _ := getCachedOutcomes(ctx, eventId)
	// if err != nil {
	// 	render.Status(r, http.StatusInternalServerError)
	// 	render.JSON(w, r, map[string]string{"error": fmt.Sprintf("Error fetching cached outcomes: %v", err)})
	// 	return
	// }

	if outcomes == nil {
		// If outcomes are not found in cache, fetch from Redis
		outcomes, err := fetchOutcomesFromRedis(ctx, eventId)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": fmt.Sprintf("Error fetching outcomes from Redis: %v", err)})
			return
		}

		// Cache fetched outcomes
		if err := cacheOutcomes(ctx, eventId, outcomes); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": fmt.Sprintf("Error caching outcomes: %v", err)})
			return
		}
	}

	render.Render(w, r, &types.OutcomeListResponse{
		OutcomeList: outcomes,
	})
}

func getCachedOutcomes(ctx context.Context, eventId string) ([]*types.OutcomeItem, error) {
	cacheKey := fmt.Sprintf("event:%s-outcomes", eventId)

	// Attempt to fetch cached outcomes
	cachedOutcomesJSON, err := database.RedisDB.Get(ctx, cacheKey).Bytes()
	if err == redis.Nil {
		// Cache miss
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Deserialize cached outcomes
	var cachedOutcomes []*types.OutcomeItem
	if err := json.Unmarshal(cachedOutcomesJSON, &cachedOutcomes); err != nil {
		return nil, err
	}

	return cachedOutcomes, nil
}

func fetchOutcomesFromRedis(ctx context.Context, eventId string) ([]*types.OutcomeItem, error) {
	pattern := fmt.Sprintf("event:%s-outcome:*", eventId)

	// Fetch keys matching the pattern
	keys, err := database.RedisDB.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var outcomes []*types.OutcomeItem

	// Limit concurrency to avoid overwhelming Redis
	concurrencyLimit := 10
	semaphore := make(chan struct{}, concurrencyLimit)

	// Iterate through each key and retrieve the corresponding outcome
	for _, key := range keys {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(key string) {
			defer func() {
				wg.Done()
				<-semaphore // Release semaphore
			}()

			// Get JSON data from Redis
			outcomeJSON, err := database.RedisDB.Get(ctx, key).Bytes()
			if err != nil {
				fmt.Printf("error fetching outcome data from Redis: %v", err)
				// Handle the error
				return
			}

			// Unmarshal JSON data into OutcomeItem struct
			var outcome types.OutcomeItem
			if err := json.Unmarshal(outcomeJSON, &outcome); err != nil {
				fmt.Printf("error unmarshaling outcome data: %v", err)
				// Handle the error
				return
			}

			// Update outcomes slice
			mu.Lock()
			defer mu.Unlock()
			outcomes = append(outcomes, &outcome)
		}(key)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return outcomes, nil
}

func cacheOutcomes(ctx context.Context, eventId string, outcomes []*types.OutcomeItem) error {
	cacheKey := fmt.Sprintf("event:%s-outcomes", eventId)

	// Serialize outcomes into JSON
	outcomesJSON, err := json.Marshal(outcomes)
	if err != nil {
		return err
	}

	// Store the JSON array in Redis under a single key
	if err := database.RedisDB.Set(ctx, cacheKey, outcomesJSON, 0).Err(); err != nil {
		return err
	}

	return nil
}
