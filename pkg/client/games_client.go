package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/internal/types/requests"
	"strconv"
	"strings"
	"sync"
	"time"
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
	// Calculate start_date_before as two months from now
	startDateBefore := time.Now().AddDate(0, 2, 0).Format("2006-01-02")

	url := fmt.Sprintf(
		"%s/api/v2/games?key=%s&start_date_before=%s&include_team_info=true&include_statsperform_ids=true",
		gc.BaseURL,
		gc.APIKey,
		startDateBefore,
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

			// oddsURL := fmt.Sprintf(
			// 	"%s/api/v2/game-odds?key=%s&sportsbook=%s&sportsbook=%s&sportsbook=%s&sportsbook=%s&sportsbook=%s",
			// 	gc.BaseURL,
			// 	gc.APIKey,
			// 	"bet365",
			// 	"bodog",
			// 	"Pinnacle",
			// 	"1XBet",
			// 	"fanduel",
			// )
			oddsURL := fmt.Sprintf(
				"%s/api/v2/game-odds?key=%s&sportsbook=%s",
				gc.BaseURL,
				gc.APIKey,
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
	var allCompetitors []*types.CompetitorItem
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

			// Collect competitors
			competitors := repositories.PrepareCompetitors(gc.BaseURL, gc.APIKey, prematch)
			allCompetitors = append(allCompetitors, competitors...)

			sportEvent, _ := repositories.CreateOrUpdateSportEvent(prematch)
			repositories.CreateOutcome(prematch, sportEvent)
			// Define the expiration time as 90 days
			expiration := 90 * 24 * time.Hour
			err = database.RedisDB.Set(ctx, key, "fetched", expiration).Err()
			if err != nil {
				fmt.Println("Error saving OutcomeItem to Redis:", err)
			}
		}
	}
	// Batch insert all competitors
	err = repositories.CreateCompetitorsBatch(allCompetitors)
	if err != nil {
		fmt.Println("Error batch inserting competitors:", err)
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

			// fmt.Println(url)
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
				// fmt.Println(item.GameId, item.Status, item.StartDate, item.ScoreHomeTotal, item.ScoreAwayTotal)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	var allSportEvents []*types.SportEventItem
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
		// repositories.UpdateSportEventStatus(sportEvent)

		allSportEvents = append(allSportEvents, sportEvent)
	}

	repositories.UpdateSportEvents(allSportEvents)
}

func (gc *GamesClient) FetchOddsAISchedule() {
	requestPayload := requests.MatchesRequest{
		Status: []string{"not_started"},
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling request payload:", err)
		return
	}

	url := fmt.Sprintf("%s/matches", gc.BaseURL)
	oddsUrl := fmt.Sprintf("%s/match-snapshots", gc.BaseURL)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	var gameScheduleResponse types.OddsAIGameScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&gameScheduleResponse); err != nil {
		fmt.Println("Error", err)
		return
	}
	// Ensure saveCompetitors fully runs before proceeding
	saveCompetitorsCompleted := make(chan struct{})
	go func() {
		saveCompetitors(gameScheduleResponse)
		close(saveCompetitorsCompleted)
	}()

	// Wait for saveCompetitors to complete
	<-saveCompetitorsCompleted

	saveSportEventsCompleted := make(chan struct{})
	go func() {
		saveSportEvents(gameScheduleResponse, oddsUrl)
		close(saveSportEventsCompleted)
	}()

	// Wait for saveSportEvents to complete
	<-saveSportEventsCompleted
}

func saveCompetitors(gameScheduleResponse types.OddsAIGameScheduleResponse) {
	sports := gameScheduleResponse.Data
	var allCompetitors []*types.CompetitorItem
	semaphore := make(chan struct{}, 100) // Limit to 100 concurrent goroutines
	var wg sync.WaitGroup

	for _, sport := range sports {
		sportName := strings.ToLower(sport.Name)
		slug := strings.ReplaceAll(sportName, " ", "_")
		sportItem, _ := repositories.GetSportFromRedis(slug)
		infos := sport.SportCountries
		for _, info := range infos {
			country := info.Country
			countryItem, _ := repositories.GetCountryFromRedis(country.Name)
			tournaments := info.Tournaments
			for _, tournament := range tournaments {
				if countryItem == nil {
					continue
				}
				tournamentItem, _ := repositories.GetTournamentFromRedis(sportItem.ReferenceId, countryItem.ReferenceId, tournament.Name)
				matches := tournament.Matches
				if tournamentItem == nil {
					continue
				}
				for _, match := range matches {
					wg.Add(1)
					semaphore <- struct{}{}
					go func(match types.Match, sportItem *types.SportItem, countryItem *types.CountryItem) {
						defer wg.Done()
						defer func() { <-semaphore }()
						home := &types.CompetitorItem{
							ReferenceId: strconv.Itoa(match.HomeTeam.ID),
							CountryId:   countryItem.Id,
							Name:        match.HomeTeam.Name,
							SportId:     sportItem.ReferenceId,
							DataFeed:    "huge_data",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						}
						if match.HomeTeam.HasLogo {
							home.Logo = fmt.Sprintf("https://cdn.betapi.win/cdn/logos/m/%s.png", strconv.Itoa(match.HomeTeam.ID))
						}
						allCompetitors = append(allCompetitors, home)

						away := &types.CompetitorItem{
							ReferenceId: strconv.Itoa(match.AwayTeam.ID),
							CountryId:   countryItem.Id,
							Name:        match.AwayTeam.Name,
							SportId:     sportItem.ReferenceId,
							DataFeed:    "huge_data",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						}
						if match.AwayTeam.HasLogo {
							away.Logo = fmt.Sprintf("https://cdn.betapi.win/cdn/logos/m/%s.png", strconv.Itoa(match.AwayTeam.ID))
						}
						allCompetitors = append(allCompetitors, away)
					}(match, sportItem, countryItem)
				}
			}
		}
	}
	wg.Wait()
	err := repositories.CreateCompetitorsBatch(allCompetitors)
	if err != nil {
		fmt.Println("Error batch inserting competitors:", err)
	}
}

func saveSportEvents(gameScheduleResponse types.OddsAIGameScheduleResponse, oddsUrl string) {
	sports := gameScheduleResponse.Data
	var allSportEvents []*types.SportEventItem
	allMatchIds := make([]int, 0)
	semaphore := make(chan struct{}, 100) // Limit to 100 concurrent goroutines
	var wg sync.WaitGroup

	for _, sport := range sports {
		sportName := strings.ToLower(sport.Name)
		slug := strings.ReplaceAll(sportName, " ", "_")
		sportItem, _ := repositories.GetSportFromRedis(slug)
		infos := sport.SportCountries
		for _, info := range infos {
			country := info.Country
			countryItem, _ := repositories.GetCountryFromRedis(country.Name)
			tournaments := info.Tournaments
			for _, tournament := range tournaments {
				if countryItem == nil {
					continue
				}
				tournamentItem, _ := repositories.GetTournamentFromRedis(sportItem.ReferenceId, countryItem.ReferenceId, tournament.Name)
				matches := tournament.Matches
				if tournamentItem == nil {
					continue
				}
				for _, match := range matches {
					wg.Add(1)
					semaphore <- struct{}{}
					go func(match types.Match, sportItem *types.SportItem, countryItem *types.CountryItem) {
						defer wg.Done()
						defer func() { <-semaphore }()
						homeTeam, _ := repositories.GetCompetitorFromRedis(strconv.Itoa(match.HomeTeam.ID))
						awayTeam, _ := repositories.GetCompetitorFromRedis(strconv.Itoa(match.AwayTeam.ID))

						fmt.Println(match.HomeTeam.ID, match.HomeTeam.Name+" vs "+match.AwayTeam.Name, match.AwayTeam.ID)
						if homeTeam != nil && awayTeam != nil {
							sportEvent := &types.SportEventItem{
								ProviderId:   2,
								ReferenceId:  strconv.Itoa(match.ID),
								SportId:      sportItem.Id,
								CountryId:    countryItem.Id,
								TournamentId: tournamentItem.Id,
								Name:         match.HomeTeam.Name + " vs " + match.AwayTeam.Name,
								HomeTeamId:   homeTeam.Id,
								AwayTeamId:   awayTeam.Id,
								Status:       "unplayed",
								Active:       1,
								DataFeed:     "huge_data",
								CreatedAt:    time.Now(),
								UpdatedAt:    time.Now(),
							}
							t := time.Unix(match.MatchDate, 0)
							sportEvent.StartAt = t.Format("2006-01-02 15:04:05.999999-07")
							if match.Status == "live" {
								sportEvent.HomeScore = int32(match.MatchInfo.HomeScore)
								sportEvent.AwayScore = int32(match.MatchInfo.AwayScore)
							} else {
								sportEvent.HomeScore = 0
								sportEvent.AwayScore = 0
							}
							// sportEvent.RoundInfo = match.MatchInfo.ScoreInfo
							allSportEvents = append(allSportEvents, sportEvent)
							allMatchIds = append(allMatchIds, match.ID)
						}
					}(match, sportItem, countryItem)
				}
			}
		}
	}
	wg.Wait()
	repositories.UpdateSportEvents(allSportEvents)
	getOdds(allMatchIds, oddsUrl)
}

func getOdds(allMatchIds []int, url string) {
	const chunkSize = 30

	for i := 0; i < len(allMatchIds); i += chunkSize {
		end := i + chunkSize
		if end > len(allMatchIds) {
			end = len(allMatchIds)
		}
		chunk := allMatchIds[i:end]

		jsonData, err := json.Marshal(chunk)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error sending POST request:", err)
			continue
		}
		defer resp.Body.Close()

		var oddsResponse types.ResponseForOdds
		if err := json.NewDecoder(resp.Body).Decode(&oddsResponse); err != nil {
			fmt.Println("Error decoding response:", err)
			continue
		}

		// Ensure saveCompetitors fully runs before proceeding
		saveMarketconstantsCompleted := make(chan struct{})
		go func() {
			saveMarketconstants(oddsResponse)
			close(saveMarketconstantsCompleted)
		}()
		// Wait for saveCompetitors to complete
		<-saveMarketconstantsCompleted

		// Ensure saveCompetitors fully runs before proceeding
		saveSportMarketGroupsCompleted := make(chan struct{})
		go func() {
			saveSportMarketGroups(oddsResponse)
			close(saveSportMarketGroupsCompleted)
		}()
		// Wait for saveCompetitors to complete
		<-saveSportMarketGroupsCompleted

		// SaveOutcomeToRedis() // Uncomment and implement if
		data := oddsResponse.Data
		for _, match := range data {
			sportEventItem, _ := repositories.GetSportEventFromRedis(strconv.Itoa(int(match.ID)))
			games := match.Games
			for _, game := range games {
				groupId := game.GameType
				markets := game.Markets
				for _, market := range markets {
					// marketRefId := strconv.Itoa(market.TemplateID)
					marketRefId := strconv.Itoa(groupId) + "-" + strconv.Itoa(market.TemplateID)
					marketConstant, _ := repositories.GetMarketConstantFromRedis(marketRefId)
					param := -1000000000000
					hasParam := false
					if market.Param != nil {
						param = market.Param.Param
						hasParam = true
					}
					odds := market.Odds
					for _, odd := range odds {
						outcomeConstant, _ := repositories.GetOutcomeConstantFromRedis(strconv.Itoa(int(odd.OutcomeID)))
						odds := toAmericanOdds(odd.Value)
						outcome := &types.OutcomeItem{
							// ReferenceId: strconv.Itoa(int(odd.OutcomeID)),
							ReferenceId: marketConstant.Description + ":" + outcomeConstant.Name,
							EventId:     sportEventItem.Id,
							MarketId:    marketConstant.Id,
							Odds:        odds,
							Name:        outcomeConstant.Name,
							Active:      true,
						}
						if hasParam {
							outcome.ReferenceId = marketConstant.Description + ":" + outcomeConstant.Name + ":" + fmt.Sprintf("%g", float64(param)/100.0)
						}
						repositories.SaveOutcomeToRedis(outcome)
					}
				}
			}
		}
	}
}

func saveMarketconstants(oddsResponse types.ResponseForOdds) {
	var marketConstants []*types.MarketConstantItem
	data := oddsResponse.Data
	for _, match := range data {
		games := match.Games
		for _, game := range games {
			groupId := game.GameType
			groupItem, _ := repositories.GetMarketGroupFromRedis(groupId)
			markets := game.Markets
			for _, market := range markets {
				marketRefId := strconv.Itoa(market.TemplateID)
				marketConstant, _ := repositories.GetMarketConstantFromRedis("1-" + marketRefId)
				if groupId != 1 {
					marketConstantItem := &types.MarketConstantItem{
						ReferenceId:  strconv.Itoa(groupId) + "-" + marketRefId,
						Description:  groupItem.MarketGroup + " " + marketConstant.Description,
						Order:        int32(market.TemplateID),
						IsTranslated: false,
						DataFeed:     "huge_data",
					}
					marketConstants = append(marketConstants, marketConstantItem)
				} else {
					marketConstants = append(marketConstants, marketConstant)
				}
			}
		}
	}
	repositories.UpdateMarketConstants(marketConstants)
}

func saveSportMarketGroups(oddsResponse types.ResponseForOdds) {
	var sportMarketGroups []*types.SportMarketGroupItem
	data := oddsResponse.Data
	for _, match := range data {
		games := match.Games
		for _, game := range games {
			groupId := game.GameType
			groupItem, _ := repositories.GetMarketGroupFromRedis(groupId)
			markets := game.Markets
			for _, market := range markets {
				marketRefId := strconv.Itoa(groupId) + "-" + strconv.Itoa(market.TemplateID)
				marketConstant, _ := repositories.GetMarketConstantFromRedis(marketRefId)
				sportMarketGroup := &types.SportMarketGroupItem{
					GroupId:    int32(groupId),
					MarketId:   marketConstant.Id,
					GroupName:  groupItem.MarketGroup,
					MarketName: marketConstant.Description,
				}
				sportMarketGroups = append(sportMarketGroups, sportMarketGroup)
			}
		}
	}
	repositories.UpdateSportMarketGroup(sportMarketGroups)
}

func toAmericanOdds(value int) float64 {
	decimalOdds := float64(value) / 1000.0
	if decimalOdds < 1.00 {
		return 0
	}

	var americanOdds float64
	if decimalOdds >= 2.00 {
		americanOdds = (decimalOdds - 1) * 100
	} else {
		americanOdds = -100 / (decimalOdds - 1)
	}

	return math.Round(americanOdds*100) / 100 // Round to 2 decimal places
}
