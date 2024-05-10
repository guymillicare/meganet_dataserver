package main

import (
	"sportsbook-backend/internal/config"
	"sportsbook-backend/internal/database"
	"sportsbook-backend/internal/grpc"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/scheduler"
	"sportsbook-backend/internal/websockets"
	"sportsbook-backend/pkg/client"
)

func main() {
	cfg := config.LoadConfig() // Load configuration

	database.InitPostgresDB(cfg)
	Preload()
	gamesClient := client.NewGamesClient(cfg.ThirdPartyAPIBaseURL, cfg.APIKey)

	// Start the scheduler to fetch games data periodically
	prematchData := &grpc.PrematchData{}                                        // assuming grpc.GamesData is a thread-safe struct
	scheduler.StartPrematchCronJob(gamesClient, prematchData, "0 */3 * * *")    // Runs every 3 hours
	scheduler.StartMatchStatusCronJob(gamesClient, prematchData, "*/2 * * * *") // Runs every 2 mins

	websockets.StartWebSocket()
	// Start the gRPC server
	grpc.StartGRPCServer(cfg.GRPCPort, prematchData)
}

func Preload() {
	repositories.SportsPreload()
	repositories.CountriesPreload()
	repositories.TournamentsPreload()
}
