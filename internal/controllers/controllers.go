package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/internal/types/requests"
	"sportsbook-backend/internal/utils"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-redis/redis"
)

func GetOutcomesByEventId(w http.ResponseWriter, r *http.Request) {
	eventRefId := chi.URLParam(r, "eventRefId")
	// event, _ := repositories.GetSportEventFromRedis(eventRefId)
	event, _ := repositories.SportEventFindByRefId(eventRefId)
	eventId := strconv.Itoa(int(event.Id))
	outcomes, err := GetOutcomes(eventId)
	var activeOutcomes []*types.OutcomeItem = []*types.OutcomeItem{}
	var marketGroups []*types.MarketGroupItem = []*types.MarketGroupItem{}
	addedMarketGroups := make(map[int]bool)
	var collectionInfos []*types.CollectionInfoItem = []*types.CollectionInfoItem{}
	addedCollectionInfos := make(map[int]bool)
	for _, outcome := range outcomes {
		// if outcome.Active {
		groupId := int(outcome.GroupId)
		collectionInfoId := outcome.CollectionInfoId
		if _, exists := addedMarketGroups[groupId]; !exists {
			marketGroup, err := repositories.GetMarketGroupFromRedis(groupId)
			if err != nil {
				log.Printf("Failed to get market group from Redis for groupId %d: %v", groupId, err)
				continue
			}
			marketGroups = append(marketGroups, marketGroup)
			addedMarketGroups[groupId] = true
		}

		if _, exists := addedCollectionInfos[int(collectionInfoId)]; !exists {
			collectionInfo, err := repositories.GetCollectionInfoFromRedis(collectionInfoId)
			if err != nil {
				log.Printf("Failed to get market group from Redis for groupId %d: %v", groupId, err)
				continue
			}
			collectionInfos = append(collectionInfos, collectionInfo)
			addedCollectionInfos[int(collectionInfoId)] = true
		}
		activeOutcomes = append(activeOutcomes, outcome)
		// }
	}
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": fmt.Sprintf("Error fetching outcomes: %v", err)})
		return
	}

	render.Render(w, r, &types.OutcomeListResponse{
		SportEvent:      event,
		MarketGroups:    marketGroups,
		CollectionInfos: collectionInfos,
		Outcome:         activeOutcomes,
	})
}

func GetOutcomes(eventId string) ([]*types.OutcomeItem, error) {
	// var err error

	ctx := context.Background()

	// Attempt to fetch outcomes from cache
	outcomes, _ := getCachedOutcomes(ctx, eventId)

	// if outcomes == nil {
	// 	// If outcomes are not found in cache, fetch from Redis
	// 	outcomes, err = fetchOutcomesFromRedis(ctx, eventId)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Cache fetched outcomes
	// 	if err = cacheOutcomes(ctx, eventId, outcomes); err != nil {
	// 		return nil, err
	// 	}
	// }
	return outcomes, nil
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

	// Define the expiration time as 90 days
	expiration := 90 * 24 * time.Hour

	// Store the JSON array in Redis under a single key
	if err := database.RedisDB.Set(ctx, cacheKey, outcomesJSON, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func GetSportEventsWithOdds(w http.ResponseWriter, r *http.Request) {
	betType := chi.URLParam(r, "betType")
	_sportId := chi.URLParam(r, "sportId")
	sportId, _ := strconv.Atoi(_sportId)
	_countryId := chi.URLParam(r, "countryId")
	countryId, _ := strconv.Atoi(_countryId)
	_leagueId := chi.URLParam(r, "leagueId")
	leagueId, _ := strconv.Atoi(_leagueId)
	// currentUser, _ := services.AuthCurrentUser(r)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values for pagination
	page := 0
	limit := 10

	// Parse pagination parameters
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := page * limit

	// sportEvents, err := repositories.SportEventsFindByFilters(int32(currentUser.SystemId), 1, betType, int32(sportId), int32(countryId), int32(leagueId), offset, limit)
	sportEvents, _ := repositories.SportEventsFindByFilters(62, 2, betType, int32(sportId), int32(countryId), int32(leagueId), offset, limit)
	var results []*types.SportEventOddsItem
	for _, event := range sportEvents {
		outcomes, _ := GetOutcomes(strconv.Itoa(int(event.Id)))
		// var activeOutcomes []*types.OutcomeItem = []*types.OutcomeItem{}
		// for _, outcome := range outcomes {
		// 	if outcome.Active {
		// 		activeOutcomes = append(activeOutcomes, outcome)
		// 	}
		// }
		var entity types.SportEventOddsItem
		entity.SportEvent = event
		entity.Outcome = outcomes
		results = append(results, &entity)
	}

	render.Render(w, r, &types.SportEventListResponse{
		SportEventList: results,
	})
}

func GetSportEventsWithLiveOdds(w http.ResponseWriter, r *http.Request) {

	var req requests.EventsList
	if err := utils.JsonParseRequestBody(w, r, false, &req); err != nil {
		return
	}

	var results []*types.SportEventsOddsItem
	for _, eventRefId := range req.EventIdList {
		event, _ := repositories.GetSportEventFromRedis(eventRefId)
		outcomes, _ := GetOutcomes(strconv.Itoa(int(event.Id)))
		var activeOutcomes []*types.OutcomeItem = []*types.OutcomeItem{}
		for _, outcome := range outcomes {
			if outcome.Active {
				activeOutcomes = append(activeOutcomes, outcome)
			}
		}
		var entity types.SportEventsOddsItem
		entity.SportEvent = event
		entity.Outcome = activeOutcomes
		results = append(results, &entity)
	}
	render.Render(w, r, &types.SportEventsListResponse{
		SportEventsList: results,
	})
}
