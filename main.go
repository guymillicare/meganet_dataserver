package main

import (
	"bufio"
	"bytes"
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
	"sportsbook-backend/pkg/client"
	"sportsbook-backend/pkg/queue"
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

	scheduler.StartOddsAIScheduleCronJob(oddsAIClient, "50 * * * *") // Runs every 1 mins

	oddsChannel := make(chan *pb.LiveData, 100) // Buffer size of 100
	scoreChannel := make(chan *pb.LiveScoreData, 100)
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
	go subscribeToFeed(rabbitMQ_odds, oddsChannel)
	go consumeRabbitMQOdds(rabbitMQ_odds, oddsChannel)
	go consumeRabbitMQScore(rabbitMQ_score, scoreChannel)

	// Start the gRPC server
	go grpc.StartGRPCServer(cfg.GRPCPort, oddsChannel, scoreChannel)
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
			fmt.Printf("GameUpdate Data: %+v\n", data.Game)
		case *datafeed.FeedUpdateData_MatchResult_:
			fmt.Printf("MatchResult Data: %+v\n", data.MatchResult)
		case *datafeed.FeedUpdateData_Settlement_:
			fmt.Printf("Settlement Data: %+v\n", data.Settlement)
		default:
			fmt.Println("Unknown data type", data)
		}
	}
}

func consumeRabbitMQScore(rabbitMQ *queue.RabbitMQ, scoreChannel chan<- *pb.LiveScoreData) {
	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for d := range msgs {
		var scoreData types.ScoreStream
		err := json.Unmarshal(d.Body, &scoreData)
		if err != nil {
			// log.Printf("Error decoding JSON: %v", err)
			// continue
		}

		// convertedScore := &pb.ScoreData{
		// 	GameId: scoreData.Data.GameID,
		// }
		convertedScore := &pb.ScoreData{
			GameId: "2024",
		}
		// score := &pb.Score{
		// 	Clock:             scoreData.Data.Score.Clock,
		// 	ScoreAwayPeriod_1: scoreData.Data.Score.ScoreAwayPeriod1,
		// 	// ScoreAwayPeriod_1Tiebreak: scoreData.Data.Score.ScoreAwayPeriod1Tiebreak,
		// 	ScoreAwayPeriod_2: scoreData.Data.Score.ScoreAwayPeriod2,
		// 	ScoreAwayTotal:    scoreData.Data.Score.ScoreAwayTotal,
		// 	ScoreHomePeriod_1: scoreData.Data.Score.ScoreHomePeriod1,
		// 	// ScoreHomePeriod_1Tiebreak: scoreData.Data.Score.ScoreHomePeriod1Tiebreak,
		// 	ScoreHomePeriod_2: scoreData.Data.Score.ScoreHomePeriod2,
		// 	ScoreHomeTotal:    scoreData.Data.Score.ScoreHomeTotal,
		// }
		score := &pb.Score{
			Clock: "2024-1-1",
		}
		convertedScore.Score = score

		// convertedScoreData := &pb.LiveScoreData{
		// 	Data:    convertedScore,
		// 	EntryId: scoreData.EntryId,
		// }
		convertedScoreData := &pb.LiveScoreData{
			Data:    convertedScore,
			EntryId: "scoreData.EntryId",
		}
		fmt.Printf("Consumer: %v\n", convertedScoreData)
		// Send live data to gRPC clients
		scoreChannel <- convertedScoreData
		// 	}
		// }
	}
}

func subscribeToFeed(rabbitMQ *queue.RabbitMQ, oddsChannel chan<- *pb.LiveData) {
	serverAddr := "demofeed.betapi.win:443"
	serviceMethod := "datafeed.FeedService/SubscribeToFeed"

	// Call the method
	cmdArgs := []string{serverAddr, serviceMethod}

	// Print the command to be executed
	fmt.Printf("Running command: grpcurl %s\n", strings.Join(cmdArgs, " "))

	// Create the command
	callCmd := exec.Command("grpcurl", cmdArgs...)

	// Get the stdout pipe to read the command's output
	stdout, err := callCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start the command
	if err := callCmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// Read the command's output in real-time
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

		// Count opening and closing braces to detect the end of the JSON object
		openBracesCount += strings.Count(line, "{")
		openBracesCount -= strings.Count(line, "}")

		if openBracesCount == 0 && buffer.Len() > 0 {
			// var oddsData types.FeedUpdateData
			jsonData := buffer.String()
			// fmt.Println("Complete JSON:", jsonData)

			rabbitMQ.Publish([]byte(jsonData))
			// err := json.Unmarshal([]byte(jsonData), &oddsData)
			// if err != nil {
			// 	log.Printf("Failed to unmarshal JSON: %v", err)
			// } else {
			// 	fmt.Printf("Unmarshalled Data: %+v\n", oddsData.Game.ID)

			// }
			buffer.Reset()
		}
	}
}
