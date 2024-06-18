package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sportsbook-backend/internal/config"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/grpc"
	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/routes"
	"sportsbook-backend/internal/scheduler"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/pkg/client"
	"sportsbook-backend/pkg/queue"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	cfg := config.LoadConfig() // Load configuration

	database.InitPostgresDB(cfg)
	database.InitRedis(cfg)
	Preload()

	// Initialize RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQ("live_odds_queue")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	gamesClient := client.NewGamesClient(cfg.ThirdPartyAPIBaseURL, cfg.APIKey)

	// Start the scheduler to fetch games data periodically
	prematchData := &grpc.PrematchData{}                                        // assuming grpc.GamesData is a thread-safe struct
	scheduler.StartPrematchCronJob(gamesClient, prematchData, "57 4 * * *")     // Runs every 3 hours
	scheduler.StartMatchStatusCronJob(gamesClient, prematchData, "*/2 * * * *") // Runs every 2 mins

	oddsChannel := make(chan *pb.LiveOddsData)
	wg := &sync.WaitGroup{}

	tournamnets, _ := repositories.TournamentsFindAll()
	for _, tournament := range tournamnets {
		urlLiveOdds := fmt.Sprintf("%s/api/v2/stream/odds?sportsbooks=bodog&sportsbooks=fanduel&sportsbooks=bet365&sportsbooks=1XBet&sportsbooks=Pinnacle&league=%s&key=%s", cfg.ThirdPartyAPIBaseURL, tournament.Name, cfg.APIKey)
		urlLiveGameScore := fmt.Sprintf("%s/api/v2/stream/results?league=%s&key=%s", cfg.ThirdPartyAPIBaseURL, tournament.Name, cfg.APIKey)
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			grpc.ListenToStream(url, oddsChannel, wg, rabbitMQ)
		}(urlLiveOdds)

		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			grpc.ListenToStream(url, oddsChannel, wg, rabbitMQ)
		}(urlLiveGameScore)
	}

	// Start the RabbitMQ consumer in a separate goroutine
	go consumeRabbitMQ(rabbitMQ, oddsChannel)

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
}

func SetupHttpHandler(APICorsAllowedOrigins string) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(middleware.UserContext) // add ctx["auth-user"] = user
	// r.Use(middleware.UserTracking)
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

func consumeRabbitMQ(rabbitMQ *queue.RabbitMQ, oddsChannel chan<- *pb.LiveOddsData) {
	msgs, err := rabbitMQ.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for d := range msgs {
		var oddsData types.OddsStream
		err := json.Unmarshal(d.Body, &oddsData)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue
		}

		// Process the oddsData
		for _, odds := range oddsData.Data {
			sportEvent, _ := repositories.GetSportEventFromRedis(odds.GameId)
			marketConstant, _ := repositories.GetMarketConstantFromRedis(odds.BetType)
			if marketConstant != nil && sportEvent != nil {
				outcome := &types.OutcomeItem{
					ReferenceId: odds.BetType + ":" + odds.BetName,
					EventId:     sportEvent.Id,
					MarketId:    marketConstant.Id,
					Name:        odds.BetName,
					Odds:        odds.BetPrice,
					Active:      oddsData.Type == "odds",
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
				}
				repositories.SaveOutcomeToRedis(outcome)

				convertedOdds := &pb.Data{
					BetName:         odds.BetName,
					BetPoints:       odds.BetPoints,
					BetPrice:        odds.BetPrice,
					BetType:         odds.BetType,
					GameId:          odds.GameId,
					Id:              odds.Id,
					IsLive:          odds.IsLive,
					IsMain:          odds.IsMain,
					League:          odds.League,
					PlayerId:        odds.PlayerId,
					Selection:       odds.Selection,
					SelectionLine:   odds.SelectionLine,
					SelectionPoints: odds.SelectionPoints,
					Sport:           odds.Sport,
					Sportsbook:      odds.Sportsbook,
					Timestamp:       odds.Timestamp,
				}

				convertedOddsData := &pb.LiveOddsData{
					EntryId: oddsData.EntryId,
					Type:    oddsData.Type,
					Data:    convertedOdds,
				}
				// Send live data to gRPC clients
				oddsChannel <- convertedOddsData
				// fmt.Printf("Consumer: %v\n", convertedOddsData)
			}
		}
	}
}
