package scheduler

import (
	"log"
	"sportsbook-backend/internal/grpc"
	"sportsbook-backend/pkg/client"

	"github.com/robfig/cron/v3"
)

// StartCronJob sets up and starts a cron job to fetch games data periodically.
func StartCronJob(client *client.GamesClient, prematchData *grpc.PrematchData, spec string) {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(spec, func() {
		log.Println("Fetching games data...")

		// Call the client method to fetch new games data.
		newPrematchData, err := client.FetchGames()
		if err != nil {
			log.Printf("Error fetching games: %v", err)
			return
		}

		// Safely update the shared games data structure.
		prematchData.UpdatePrematchData(newPrematchData)

		log.Println("Games data successfully updated.")
	})
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
	}

	// Start the cron scheduler.
	c.Start()

	// Optionally, you can handle graceful shutdown:
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, os.Interrupt)
	// <-sig
	// c.Stop()
}
