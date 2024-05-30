package main

import (
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
	"sportsbook-backend/pkg/client"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	cfg := config.LoadConfig() // Load configuration

	database.InitPostgresDB(cfg)
	database.InitRedis(cfg)
	Preload()

	gamesClient := client.NewGamesClient(cfg.ThirdPartyAPIBaseURL, cfg.APIKey)

	// Start the scheduler to fetch games data periodically
	prematchData := &grpc.PrematchData{}                                        // assuming grpc.GamesData is a thread-safe struct
	scheduler.StartPrematchCronJob(gamesClient, prematchData, "36 * * * *")     // Runs every 3 hours
	scheduler.StartMatchStatusCronJob(gamesClient, prematchData, "*/2 * * * *") // Runs every 2 mins

	oddsChannel := make(chan *pb.LiveOddsData)
	wg := &sync.WaitGroup{}

	tournamnets, _ := repositories.TournamentsFindAll()
	for _, tournament := range tournamnets {
		url := fmt.Sprintf("%s/api/v2/stream/odds?sportsbooks=betsson&sportsbooks=bet365&sportsbooks=1XBet&sportsbooks=Pinnacle&league=%s&key=%s", cfg.ThirdPartyAPIBaseURL, tournament.Name, cfg.APIKey)
		wg.Add(1)
		go grpc.ListenToStream(url, oddsChannel, wg)
	}
	// Start the gRPC server
	grpc.StartGRPCServer(cfg.GRPCPort, oddsChannel)
	// Start the HTTP server
	handler := SetupHttpHandler(cfg.APICorsAllowedOrigins)
	port := 9000
	fmt.Printf("Using port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
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
