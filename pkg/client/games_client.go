package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/types"
	"sync"
)

type GamesClient struct {
	APIKey  string
	BaseURL string
}

func NewGamesClient(baseURL string, apiKey string) *GamesClient {
	return &GamesClient{BaseURL: baseURL, APIKey: apiKey}
}

func (gc *GamesClient) FetchGames() (*proto.ListPrematchResponse, error) {
	var responses proto.ListPrematchResponse
	// sports, _ := repositories.SportsFindAll()
	// for _, sport := range sports {
	url := fmt.Sprintf(
		"%s/api/v2/games?key=%s",
		gc.BaseURL,
		gc.APIKey,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gameListResponse proto.ListPrematchResponse
	if err := json.NewDecoder(resp.Body).Decode(&gameListResponse); err != nil {
		return nil, err
	}

	oddsMap := make(map[string][]*proto.Odds)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Fetch odds for each set of five games concurrently
	for i := 0; i < len(gameListResponse.Data); i += 5 {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			end := start + 5
			if end > len(gameListResponse.Data) {
				end = len(gameListResponse.Data)
			}

			gameIDs := make([]string, 0, 5)
			for j := start; j < end; j++ {
				gameIDs = append(gameIDs, gameListResponse.Data[j].Id)
			}

			oddsURL := fmt.Sprintf(
				"%s/api/v2/game-odds?key=%s&sportsbook=%s&sportsbook=%s&sportsbook=%s&sportsbook=%s",
				gc.BaseURL,
				gc.APIKey,
				"bet365",
				"betsson",
				"Pinnacle",
				"1XBet",
			)
			for _, gameID := range gameIDs {
				oddsURL += fmt.Sprintf("&game_id=%s", gameID)
			}

			res, err := http.Get(oddsURL)
			if err != nil {
				return
			}
			defer res.Body.Close()

			var oddsResponse struct {
				Data []struct {
					Odds   []*proto.Odds `json:"odds"`
					GameID string        `json:"id"`
				} `json:"data"`
			}
			if err := json.NewDecoder(res.Body).Decode(&oddsResponse); err != nil {
				return
			}

			mu.Lock()
			for _, item := range oddsResponse.Data {
				oddsMap[item.GameID] = item.Odds
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// Filter out games with empty odds
	var filteredData []*proto.Prematch
	for _, prematch := range gameListResponse.Data {
		ctx := context.Background()
		key := fmt.Sprintf("fetched:match:%s", prematch.Id)
		val, err := database.RedisDB.Get(ctx, key).Result()
		if err == nil && val == "fetched" {
			repositories.CreateOrUpdateSportEvent(prematch)
			continue
		}
		if odds, ok := oddsMap[prematch.Id]; ok && len(odds) > 0 {
			prematch.Odds = odds
			filteredData = append(filteredData, prematch)
			repositories.CreateCompetitor(prematch)
			sportEvent, _ := repositories.CreateOrUpdateSportEvent(prematch)
			repositories.CreateOutcome(prematch, sportEvent)
		}
		err = database.RedisDB.Set(ctx, key, "fetched", 0).Err()
		if err != nil {
			fmt.Println("Error saving OutcomeItem to Redis:", err)
		}
	}
	gameListResponse.Data = filteredData

	responses.Data = append(responses.Data, gameListResponse.Data...)
	// }
	return &responses, nil
}

func (gc *GamesClient) FetchStatus() {
	sportEvents, err := repositories.SportEventsFindAll()
	if err != nil {
		return
	}
	scoreMap := make(map[string]string)
	scoreStartDateMap := make(map[string]string)
	scoreHomeMap := make(map[string]int)
	scoreAwayMap := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < len(sportEvents); i += 5 {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			end := start + 5
			if end > len(sportEvents) {
				end = len(sportEvents)
			}

			gameIDs := make([]string, 0, 5)
			for j := start; j < end; j++ {
				gameIDs = append(gameIDs, sportEvents[j].ReferenceId)
			}

			url := fmt.Sprintf(
				"%s/api/v2/scores?key=%s",
				gc.BaseURL,
				gc.APIKey,
			)

			for _, gameID := range gameIDs {
				url += fmt.Sprintf("&game_id=%s", gameID)
			}

			fmt.Println(url)
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			var gameScoreResponse types.GameScoreResponse
			if err := json.NewDecoder(resp.Body).Decode(&gameScoreResponse); err != nil {
				fmt.Println("Error", err)
				return
			}
			mu.Lock()
			for _, item := range gameScoreResponse.Data {
				scoreMap[item.GameId] = item.Status
				scoreStartDateMap[item.GameId] = item.StartDate
				scoreHomeMap[item.GameId] = item.ScoreHomeTotal
				scoreAwayMap[item.GameId] = item.ScoreAwayTotal
				fmt.Println(item.GameId, item.Status, item.StartDate, item.ScoreHomeTotal, item.ScoreAwayTotal)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	for _, sportEvent := range sportEvents {
		if status, ok := scoreMap[sportEvent.ReferenceId]; ok {
			sportEvent.Status = status
		}
		if startDate, ok := scoreStartDateMap[sportEvent.ReferenceId]; ok {
			sportEvent.StartAt = startDate
		}
		if homeScore, ok := scoreHomeMap[sportEvent.ReferenceId]; ok {
			sportEvent.HomeScore = int32(homeScore)
		}
		if awayScore, ok := scoreAwayMap[sportEvent.ReferenceId]; ok {
			sportEvent.AwayScore = int32(awayScore)
		}
		repositories.UpdateSportEventStatus(sportEvent)
	}
}
