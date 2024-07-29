package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sportsbook-backend/internal/config"
	"sportsbook-backend/internal/controllers"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/datafeed"
	"sportsbook-backend/internal/grpc"
	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/routes"
	"sportsbook-backend/internal/scheduler"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/internal/utils"
	"sportsbook-backend/pkg/client"
	"sportsbook-backend/pkg/queue"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	cfg := config.LoadConfig() // Load configuration

	database.InitPostgresDB(cfg)
	database.InitRedis(cfg)
	Preload()

	// Initialize RabbitMQ
	rabbitMQ_odds, err := queue.NewRabbitMQ("live_odds_queue")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ_odds.Close()
	rabbitMQ_score, err := queue.NewRabbitMQ("live_score_queue")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ_score.Close()

	// opticOddsClient := client.NewGamesClient(cfg.OpticOddsAPIBaseURL, cfg.APIKey)
	oddsAIClient := client.NewGamesClient(cfg.OddsAIAPIBaseURL, "")

	// Start the scheduler to fetch games data periodically
	// prematchData := &grpc.PrematchData{}                                          // assuming grpc.GamesData is a thread-safe struct
	// scheduler.StartPrematchCronJob(opticOddsClient, prematchData, "32 */2 * * *") // Runs every 3 hours
	// scheduler.StartMatchStatusCronJob(opticOddsClient, prematchData, "* * * * *") // Runs every 1 mins

	scheduler.StartOddsAIScheduleCronJob(oddsAIClient, "36 * * * *") // Runs every 1 mins

	oddsChannel := make(chan *pb.LiveData, 100) // Buffer size of 100
	// wg := &sync.WaitGroup{}

	// tournamnets, _ := repositories.TournamentsFindAll()
	// const batchSize = 10
	// for i := 0; i < len(tournamnets); i += batchSize {
	// 	end := i + batchSize
	// 	if end > len(tournamnets) {
	// 		end = len(tournamnets)
	// 	}
	// 	batch := tournamnets[i:end]

	// 	urlLiveOdds := fmt.Sprintf("%s/api/v2/stream/odds?sportsbooks=1XBet&key=%s", cfg.OpticOddsAPIBaseURL, cfg.APIKey)
	// 	// urlLiveOdds := fmt.Sprintf("%s/api/v2/stream/odds?sportsbooks=bodog&sportsbooks=fanduel&sportsbooks=bet365&sportsbooks=1XBet&sportsbooks=Pinnacle&key=%s", cfg.ThirdPartyAPIBaseURL, cfg.APIKey)
	// 	urlLiveGameScore := fmt.Sprintf("%s/api/v2/stream/results?key=%s", cfg.OpticOddsAPIBaseURL, cfg.APIKey)

	// 	for _, tournament := range batch {
	// 		league := tournament.Name
	// 		if tournament.Name != tournament.CountryName {
	// 			league = tournament.CountryName + " - " + tournament.Name
	// 		}

	// 		urlLiveOdds += fmt.Sprintf("&league=%s", league)
	// 		urlLiveGameScore += fmt.Sprintf("&league=%s", league)
	// 	}

	// 	wg.Add(1)
	// 	go func(url string) {
	// 		defer wg.Done()
	// 		grpc.ListenToOddsStream(url, oddsChannel, wg, rabbitMQ_odds)
	// 	}(urlLiveOdds)

	// 	// wg.Add(1)
	// 	// go func(url string) {
	// 	// 	defer wg.Done()
	// 	// 	grpc.ListenToScoreStream(url, scoreChannel, wg, rabbitMQ_score)
	// 	// }(urlLiveGameScore)
	// }

	// Start the RabbitMQ consumer in a separate goroutine
	go subscribeToFeed(rabbitMQ_odds)
	go consumeRabbitMQOdds(rabbitMQ_odds, oddsChannel)
	// go consumeRabbitMQScore(rabbitMQ_score, scoreChannel)

	// Start the gRPC server
	go grpc.StartGRPCServer(cfg.GRPCPort, oddsChannel)
	// Start the HTTP server
	handler := SetupHttpHandler(cfg.APICorsAllowedOrigins)
	port := 9000
	fmt.Printf("Using port %d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}

func Preload() {
	repositories.SportsPreload()
	repositories.CountriesPreload()
	repositories.TournamentsPreload()
	repositories.SportEventsPreload()
	repositories.MarketConstantsPreload()
	repositories.OutcomeConstantsPreload()
	repositories.CompetitorsPreload()
	repositories.MarketGroupsPreload()
	repositories.CollectionInfosPreload()
}

func SetupHttpHandler(APICorsAllowedOrigins string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(APICorsAllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	routes.SetupRouter(r)

	return r
}

func consumeRabbitMQOdds(rabbitMQ *queue.RabbitMQ, oddsChannel chan<- *pb.LiveData) {
	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	var wg sync.WaitGroup

	// Use a worker pool to process messages concurrently
	numWorkers := 50 // Adjust the number of workers as needed
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for d := range msgs {
				var feedUpdateData datafeed.FeedUpdateData
				err := protojson.Unmarshal(d.Body, &feedUpdateData)
				if err != nil {
					log.Printf("Error decoding JSON: %v", err)
					continue
				}
				processFeedUpdateData(feedUpdateData, oddsChannel)
			}
		}()
	}

	wg.Wait()
}

func processFeedUpdateData(feedUpdateData datafeed.FeedUpdateData, oddsChannel chan<- *pb.LiveData) {
	var convertedData *pb.LiveData

	switch data := feedUpdateData.Data.(type) {
	case *datafeed.FeedUpdateData_Match:
		convertedData = &pb.LiveData{
			Data: getMatchUpdateData(data.Match),
		}
	case *datafeed.FeedUpdateData_Game:
		convertedData = &pb.LiveData{
			Data: getGameUpdateData(data.Game),
		}
	case *datafeed.FeedUpdateData_MatchResult_:
		fmt.Printf("MatchResult Data: %+v\n", data.MatchResult)
	case *datafeed.FeedUpdateData_Settlement_:
		fmt.Printf("Settlement Data: %+v\n", data.Settlement)
	default:
		fmt.Println("Unknown data type", data)
	}

	if convertedData != nil {
		oddsChannel <- convertedData
	}
}

func getMatchUpdateData(match *datafeed.FeedUpdateData_MatchUpdate) *pb.LiveData_Match {
	eventRefId := strconv.Itoa(int(match.Id))
	sportEvent, _ := repositories.GetSportEventFromRedis(eventRefId)
	if sportEvent == nil {
		return nil
	}
	sportEvent.StartAt = time.Unix(match.MatchDate, 0).Format(time.RFC3339)
	if match.Status == datafeed.EventStatus_live {
		sportEvent.Status = "Live"
	} else if match.Status == datafeed.EventStatus_not_started {
		sportEvent.Status = "unplayed"
	} else {
		sportEvent.Status = "Completed"
	}

	var decodedData []byte
	if len(match.MatchInfo) > 0 {
		if match.MatchInfo[0] == '{' {
			decodedData = match.MatchInfo
		} else {
			decodedData, _ = base64.StdEncoding.DecodeString(string(match.MatchInfo))
		}
		var gameInfo types.GameInfo
		err := json.Unmarshal(decodedData, &gameInfo)
		if err != nil {
			fmt.Println("Failed to unmarshal JSON: %v", err)
		}
		if sportEvent != nil {
			sportEvent.HomeScore = int32(*gameInfo.HomeScore)
			sportEvent.AwayScore = int32(*gameInfo.AwayScore)
			if gameInfo.ScoreInfo != nil {
				sportEvent.RoundInfo = *gameInfo.ScoreInfo
			}
			if gameInfo.Tmr != nil {
				sportEvent.Tmr = *gameInfo.Tmr
			}
			if gameInfo.TmrRunning != nil {
				sportEvent.TmrRunning = *gameInfo.TmrRunning
				sportEvent.Tmr = true
			}
			if gameInfo.TmrUpdate != nil {
				sportEvent.TmrUpdate = int64(*gameInfo.TmrUpdate)
				sportEvent.Tmr = true
			}
			if gameInfo.TmrSecond != nil {
				sportEvent.TmrUpdate = time.Now().Unix() - int64(*gameInfo.TmrSecond)
				sportEvent.Tmr = true
			}
		}
	}

	repositories.UpdateSportEventStatus(sportEvent)

	status := pb.EventStatus_live
	if match.Status == datafeed.EventStatus_live {
		status = pb.EventStatus_live
	} else if match.Status == datafeed.EventStatus_not_started {
		status = pb.EventStatus_not_started
	} else {
		status = pb.EventStatus_not_active
	}

	var convertedMatchData *pb.LiveData_Match = &pb.LiveData_Match{
		Match: &pb.LiveData_MatchUpdate{
			MatchDate: match.MatchDate,
			MatchInfo: match.MatchInfo,
			Id:        match.Id,
			Status:    status,
		},
	}

	return convertedMatchData
}

func getGameUpdateData(game *datafeed.FeedUpdateData_GameUpdate) *pb.LiveData_Odds {
	var convertedOdds []*pb.LiveData_LiveOddsData_OddsData
	sportEvent, err := repositories.GetSportEventFromRedis(strconv.Itoa(int(game.MatchId)))
	if err != nil || sportEvent == nil {
		return nil
	}

	outcomes, err := controllers.GetOutcomes(strconv.Itoa(int(sportEvent.Id)))
	if err != nil {
		return nil
	}

	for _, outcome := range outcomes {
		if outcome.GroupId == game.GameType {
			outcome.Active = false
			if err := repositories.SaveOutcomeToRedis(outcome); err != nil {
				return nil
			}
		}
	}

	for _, market := range game.Markets {
		marketRefId := fmt.Sprintf("%d-%d", game.GameType, market.MarketTemplate)
		marketConstant, err := repositories.GetMarketConstantFromRedis(marketRefId)
		if err != nil || marketConstant == nil {
			continue
		}

		var high, low, param int
		var hasParam, hasOneParam bool

		if market.Param != nil {
			if market.Param.High != nil {
				high = int(*market.Param.High)
				low = int(*market.Param.Low)
			} else {
				param = int(*market.Param.Param)
				hasOneParam = true
			}
			hasParam = true
		}

		for _, odd := range market.Odds {
			outcomeConstant, err := repositories.GetOutcomeConstantFromRedis(strconv.Itoa(int(odd.OutcomeId)))
			if err != nil || outcomeConstant == nil {
				continue
			}

			referenceID := generateReferenceID(marketConstant.Description, outcomeConstant.Name, hasParam, hasOneParam, param, low, high)
			collectionInfo, err := repositories.GetCollectionInfoFromRedis(int32(market.GroupId))
			if err != nil {
				continue
			}

			outcome := &types.OutcomeItem{
				ReferenceId:      referenceID,
				EventId:          sportEvent.Id,
				MarketId:         marketConstant.Id,
				GroupId:          game.GameType,
				CollectionInfoId: collectionInfo.Id,
				Odds:             float64(utils.ToAmericanOdds(int(odd.Value))),
				Name:             outcomeConstant.Name,
				Active:           !odd.Blocked && odd.Active,
				OutcomeId:        outcomeConstant.Id,
				OutcomeOrder:     int(odd.OutcomeId),
			}

			if err := repositories.SaveOutcomeToRedis(outcome); err != nil {
				continue
			}
		}
	}

	updatedOutcomes, err := controllers.GetOutcomes(strconv.Itoa(int(sportEvent.Id)))
	if err != nil {
		return nil
	}

	for _, outcome := range updatedOutcomes {
		if outcome.GroupId == game.GameType {
			convertedOdds = append(convertedOdds, convertToPbOddsData(outcome, game.MatchId))
		}
	}

	convertedOddsData := &pb.LiveData_Odds{
		Odds: &pb.LiveData_LiveOddsData{
			MatchId:  strconv.Itoa(int(game.MatchId)),
			Status:   pb.EventStatus(game.Status),
			GameInfo: game.GameInfo,
			Odds:     convertedOdds,
		},
	}

	updateSportEventStatus(sportEvent, game)

	return convertedOddsData
}

func generateReferenceID(description, name string, hasParam, hasOneParam bool, param, low, high int) string {
	if !hasParam {
		return fmt.Sprintf("%s;%s", description, name)
	}

	if hasOneParam {
		return fmt.Sprintf("%s,%g;%s", description, float64(param)/100.0, name)
	}

	switch {
	case high != 0 && low != 0 && high != 255 && low != 255:
		return fmt.Sprintf("%s,%d or %d;%s", description, low, high, name)
	case high == 0:
		return fmt.Sprintf("%s,exact %d;%s", description, low, name)
	case high == 255:
		return fmt.Sprintf("%s,%d and more;%s", description, low, name)
	case low == 0 && high != 0:
		return fmt.Sprintf("%s,%d and less;%s", description, high, name)
	default:
		return fmt.Sprintf("%s;%s", description, name)
	}
}

func convertToPbOddsData(outcome *types.OutcomeItem, matchId int32) *pb.LiveData_LiveOddsData_OddsData {
	return &pb.LiveData_LiveOddsData_OddsData{
		ReferenceId:      outcome.ReferenceId,
		Odds:             outcome.Odds,
		GameId:           strconv.Itoa(int(matchId)),
		Active:           outcome.Active,
		GroupId:          outcome.GroupId,
		CollectionInfoId: outcome.CollectionInfoId,
		MarketId:         outcome.MarketId,
		OutcomeId:        int32(outcome.OutcomeId),
		OutcomeOrder:     int32(outcome.OutcomeOrder),
		Name:             outcome.Name,
	}
}

func updateSportEventStatus(sportEvent *types.SportEventItem, game *datafeed.FeedUpdateData_GameUpdate) {
	if pb.EventStatus(game.Status) == pb.EventStatus_live {
		sportEvent.Status = "Live"
	} else if pb.EventStatus(game.Status) == pb.EventStatus_not_started {
		sportEvent.Status = "unplayed"
	} else {
		sportEvent.Status = "Completed"
	}

	var decodedData []byte
	if len(game.GameInfo) > 0 {
		if game.GameInfo[0] == '{' {
			decodedData = game.GameInfo
		} else {
			decodedData, _ = base64.StdEncoding.DecodeString(string(game.GameInfo))
		}
		var gameInfo types.GameInfo
		if err := json.Unmarshal(decodedData, &gameInfo); err != nil {
			fmt.Println("Failed to unmarshal JSON: %v", err)
		}

		if game.GameType == 1 {
			sportEvent.HomeScore = int32(*gameInfo.HomeScore)
			sportEvent.AwayScore = int32(*gameInfo.AwayScore)
			if gameInfo.ScoreInfo != nil {
				sportEvent.RoundInfo = *gameInfo.ScoreInfo
			}
		}
		if gameInfo.Tmr != nil {
			sportEvent.Tmr = *gameInfo.Tmr
		}
		if gameInfo.TmrRunning != nil {
			sportEvent.TmrRunning = *gameInfo.TmrRunning
			sportEvent.Tmr = true
		}
		if gameInfo.TmrUpdate != nil {
			sportEvent.TmrUpdate = int64(*gameInfo.TmrUpdate)
			sportEvent.Tmr = true
		}
		if gameInfo.TmrSecond != nil {
			sportEvent.TmrUpdate = time.Now().Unix() - int64(*gameInfo.TmrSecond)
			sportEvent.Tmr = true
		}
	}

	repositories.UpdateSportEventStatus(sportEvent)
}

func subscribeToFeed(rabbitMQ *queue.RabbitMQ) {
	serverAddr := "demofeed.betapi.win:443"
	serviceMethod := "datafeed.FeedService/SubscribeToFeed"

	cmdArgs := []string{serverAddr, serviceMethod}
	callCmd := exec.Command("grpcurl", cmdArgs...)
	stdout, err := callCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	if err := callCmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	var buffer bytes.Buffer
	openBracesCount := 0
	dataType := 0

	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(line)
		buffer.WriteString("\n")

		if dataType == 0 && line == "match" {
			dataType = 1
		}
		if dataType == 0 && line == "game" {
			dataType = 2
		}
		if dataType == 0 && line == "match_result" {
			dataType = 3
		}
		if dataType == 0 && line == "settlement" {
			dataType = 4
		}

		openBracesCount += strings.Count(line, "{")
		openBracesCount -= strings.Count(line, "}")

		if openBracesCount == 0 && buffer.Len() > 0 {
			jsonData := buffer.String()
			rabbitMQ.Publish([]byte(jsonData))
			buffer.Reset()
		}
	}
}
