package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
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
	for _, sport := range database.SPORTS {
		url := fmt.Sprintf(
			"%s/api/v2/games?key=%s&sport=%s",
			gc.BaseURL,
			gc.APIKey,
			sport.Slug,
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
					"%s/api/v2/game-odds?key=%s&sportsbook=bet365",
					gc.BaseURL,
					gc.APIKey,
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
			if odds, ok := oddsMap[prematch.Id]; ok && len(odds) > 0 {
				prematch.Odds = odds
				filteredData = append(filteredData, prematch)
				repositories.CreateCompetitor(prematch)
				sportEvent, _ := repositories.CreateSportEvent(prematch)
				repositories.CreateOutcome(prematch, sportEvent)
			}
		}
		gameListResponse.Data = filteredData

		responses.Data = append(responses.Data, gameListResponse.Data...)
	}
	return &responses, nil
}

func (gc *GamesClient) FetchStatus() (*proto.ListPrematchResponse, error) {
	var responses proto.ListPrematchResponse
	for _, sport := range database.SPORTS {
		url := fmt.Sprintf(
			"%s/api/v2/games?key=%s&sport=%s",
			gc.BaseURL,
			gc.APIKey,
			sport.Slug,
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
		for _, prematch := range gameListResponse.Data {
			repositories.CreateSportEvent(prematch)
		}
	}
	return &responses, nil
}
