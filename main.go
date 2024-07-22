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

	scheduler.StartOddsAIScheduleCronJob(oddsAIClient, "11 * * * *") // Runs every 1 mins

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

	// httpsCertFile := "cert.pem" // Replace with the path to your cert file
	// httpsKeyFile := "key.pem"   // Replace with the path to your key file

	// err = http.ListenAndServeTLS(fmt.Sprintf(":%d", port), httpsCertFile, httpsKeyFile, handler)
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
	// fmt.Println(APICorsAllowedOrigins)
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

	// Add this middleware for API key check
	// r.Use(CheckAPIKeyMiddleware)

	routes.SetupRouter(r)

	return r
}

// func CheckAPIKeyMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		apiKey := r.Header.Get("X-API-Key")
// 		fmt.Printf("Origin: %s\n", r.Header.Get("Origin"))
// 		if apiKey != "12345678" { // Replace with your actual API key
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func consumeRabbitMQOdds(rabbitMQ *queue.RabbitMQ, oddsChannel chan<- *pb.LiveData) {
	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	var convertedOddsData *pb.LiveData

	for d := range msgs {
		// var oddsData types.OddsStream
		var feedUpdateData datafeed.FeedUpdateData
		err := protojson.Unmarshal(d.Body, &feedUpdateData)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue
		}
		switch data := feedUpdateData.Data.(type) {
		case *datafeed.FeedUpdateData_Match:
			fmt.Printf("MatchUpdate Data: %+v\n", data.Match)
		case *datafeed.FeedUpdateData_Game:
			// fmt.Printf("GameUpdate Data: %+v\n", data.Game)
			// convertedOddsData = getGameUpdateData(data.Game)
			convertedOddsData = &pb.LiveData{
				Data: getGameUpdateData(data.Game),
			}

		case *datafeed.FeedUpdateData_MatchResult_:
			fmt.Printf("MatchResult Data: %+v\n", data.MatchResult)
		case *datafeed.FeedUpdateData_Settlement_:
			fmt.Printf("Settlement Data: %+v\n", data.Settlement)
		default:
			fmt.Println("Unknown data type", data)
		}

		oddsChannel <- convertedOddsData
	}
}

func getGameUpdateData(game *datafeed.FeedUpdateData_GameUpdate) *pb.LiveData_Odds {
	var convertedOdds []*pb.LiveData_LiveOddsData_OddsData
	for _, market := range game.Markets {
		marketRefId := strconv.Itoa(int(market.GroupId)) + "-" + strconv.Itoa(int(market.MarketTemplate))
		marketConstant, _ := repositories.GetMarketConstantFromRedis(marketRefId)
		param := 0
		high := 0
		low := 0
		hasParam := false
		hasOneParam := false
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
			outcomeConstant, _ := repositories.GetOutcomeConstantFromRedis(strconv.Itoa(int(odd.OutcomeId)))
			if marketConstant == nil || outcomeConstant == nil {
				continue
			}
			convertedOdd := &pb.LiveData_LiveOddsData_OddsData{
				ReferenceId: marketConstant.Description + ";" + outcomeConstant.Name,
				BetPrice:    utils.ToAmericanOdds(int(odd.Value)),
				GameId:      strconv.Itoa(int(game.MatchId)),
				Active:      odd.Active,
			}
			// fmt.Println("convertedOdd.ReferenceId", convertedOdd.ReferenceId)
			if hasParam {
				if hasOneParam {
					convertedOdd.ReferenceId = marketConstant.Description + "," + fmt.Sprintf("%g", float64(param)/100.0) + ";" + outcomeConstant.Name
				} else {
					if high != 0 && low != 0 && high != 255 && low != 255 {
						convertedOdd.ReferenceId = marketConstant.Description + "," + fmt.Sprintf("%d or %d", low, high) + ";" + outcomeConstant.Name
					} else if high == 0 {
						convertedOdd.ReferenceId = marketConstant.Description + "," + fmt.Sprintf("exact %d", low) + ";" + outcomeConstant.Name
					} else if high == 255 {
						convertedOdd.ReferenceId = marketConstant.Description + "," + fmt.Sprintf("%d and more", low) + ";" + outcomeConstant.Name
					} else if low == 0 && high != 0 {
						convertedOdd.ReferenceId = marketConstant.Description + "," + fmt.Sprintf("%d and less", high) + ";" + outcomeConstant.Name
					}
				}
			}
			convertedOdds = append(convertedOdds, convertedOdd)
		}
	}

	var convertedOddsData *pb.LiveData_Odds = &pb.LiveData_Odds{
		Odds: &pb.LiveData_LiveOddsData{
			MatchId:  strconv.Itoa(int(game.MatchId)),
			Status:   pb.EventStatus(game.Status),
			GameInfo: game.GameInfo,
			Odds:     convertedOdds,
		},
	}
	sportEvent, _ := repositories.GetSportEventFromRedis(strconv.Itoa(int(game.MatchId)))
	var decodedData []byte
	if len(game.GameInfo) > 0 {
		if game.GameInfo[0] == '{' {
			decodedData = game.GameInfo
		} else {
			decodedData, _ = base64.StdEncoding.DecodeString(string(game.GameInfo))
		}
		var gameInfo types.GameInfo
		err := json.Unmarshal(decodedData, &gameInfo)
		if err != nil {
			fmt.Println("Failed to unmarshal JSON: %v", err)
		}
		if sportEvent != nil {
			sportEvent.HomeScore = int32(gameInfo.HomeScore)
			sportEvent.AwayScore = int32(gameInfo.AwayScore)
			if gameInfo.ScoreInfo != nil {
				sportEvent.RoundInfo = *gameInfo.ScoreInfo
			}
		}
	}

	if sportEvent != nil {
		if pb.EventStatus(game.Status) == pb.EventStatus_live {
			sportEvent.Status = "Live"
		} else if pb.EventStatus(game.Status) == pb.EventStatus_not_started {
			sportEvent.Status = "unplayed"
		} else {
			sportEvent.Status = "Completed"
		}

		repositories.UpdateSportEventStatus(sportEvent)
	}

	return convertedOddsData
}

// func consumeRabbitMQScore(rabbitMQ *queue.RabbitMQ, scoreChannel chan<- *pb.LiveScoreData) {
// 	msgs, err := rabbitMQ.Consume()
// 	if err != nil {
// 		log.Fatalf("Failed to consume messages: %v", err)
// 	}

// 	for d := range msgs {
// 		var scoreData types.ScoreStream
// 		err := json.Unmarshal(d.Body, &scoreData)
// 		if err != nil {
// 			// log.Printf("Error decoding JSON: %v", err)
// 			// continue
// 		}

// 		// convertedScore := &pb.ScoreData{
// 		// 	GameId: scoreData.Data.GameID,
// 		// }
// 		convertedScore := &pb.ScoreData{
// 			GameId: "2024",
// 		}
// 		// score := &pb.Score{
// 		// 	Clock:             scoreData.Data.Score.Clock,
// 		// 	ScoreAwayPeriod_1: scoreData.Data.Score.ScoreAwayPeriod1,
// 		// 	// ScoreAwayPeriod_1Tiebreak: scoreData.Data.Score.ScoreAwayPeriod1Tiebreak,
// 		// 	ScoreAwayPeriod_2: scoreData.Data.Score.ScoreAwayPeriod2,
// 		// 	ScoreAwayTotal:    scoreData.Data.Score.ScoreAwayTotal,
// 		// 	ScoreHomePeriod_1: scoreData.Data.Score.ScoreHomePeriod1,
// 		// 	// ScoreHomePeriod_1Tiebreak: scoreData.Data.Score.ScoreHomePeriod1Tiebreak,
// 		// 	ScoreHomePeriod_2: scoreData.Data.Score.ScoreHomePeriod2,
// 		// 	ScoreHomeTotal:    scoreData.Data.Score.ScoreHomeTotal,
// 		// }
// 		score := &pb.Score{
// 			Clock: "2024-1-1",
// 		}
// 		convertedScore.Score = score

// 		// convertedScoreData := &pb.LiveScoreData{
// 		// 	Data:    convertedScore,
// 		// 	EntryId: scoreData.EntryId,
// 		// }
// 		convertedScoreData := &pb.LiveScoreData{
// 			Data:    convertedScore,
// 			EntryId: "scoreData.EntryId",
// 		}
// 		fmt.Printf("Consumer: %v\n", convertedScoreData)
// 		// Send live data to gRPC clients
// 		scoreChannel <- convertedScoreData
// 		// 	}
// 		// }
// 	}
// }

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
