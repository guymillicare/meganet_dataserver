package main

import (
	"fmt"
	"net/http"
	"sportsbook-backend/internal/config"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/grpc"
	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/routes"
	"sportsbook-backend/internal/scheduler"
	"sportsbook-backend/pkg/client"
	"sync"
)

func main() {
	cfg := config.LoadConfig() // Load configuration

	database.InitPostgresDB(cfg)
	database.InitRedis()
	Preload()

	gamesClient := client.NewGamesClient(cfg.ThirdPartyAPIBaseURL, cfg.APIKey)

	// Start the scheduler to fetch games data periodically
	prematchData := &grpc.PrematchData{}                                        // assuming grpc.GamesData is a thread-safe struct
	scheduler.StartPrematchCronJob(gamesClient, prematchData, "55 * * * *")     // Runs every 3 hours
	scheduler.StartMatchStatusCronJob(gamesClient, prematchData, "*/2 * * * *") // Runs every 2 mins

	oddsChannel := make(chan *pb.LiveOddsData)
	// go grpc.ListenToStream(oddsChannel)

	urls := []string{
		"https://api.opticodds.com/api/v2/stream/odds?sportsbooks=bwin&league=England - Premier League&key=88f9bd7f-463c-44ca-b938-fd5bf2704e52",
		"https://api.opticodds.com/api/v2/stream/odds?sportsbooks=bwin&league=Spain - La Liga&key=88f9bd7f-463c-44ca-b938-fd5bf2704e52",
		// Add more URLs as needed
	}
	wg := &sync.WaitGroup{}

	// Start listening to all URLs concurrently
	for _, url := range urls {
		wg.Add(1)
		go grpc.ListenToStream(url, oddsChannel, wg)
	}
	router := routes.SetupRouter()
	fmt.Printf("Using port %d\n", 9000)
	http.ListenAndServe(":9000", router)
	// Start the gRPC server
	grpc.StartGRPCServer(cfg.GRPCPort, prematchData, oddsChannel)

}

func Preload() {
	repositories.SportsPreload()
	repositories.CountriesPreload()
	repositories.TournamentsPreload()
}
