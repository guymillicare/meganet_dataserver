package grpc

import (
	"sync"

	pb "sportsbook-backend/internal/proto" // assuming the protobuf generated code is in this package
)

// GamesData encapsulates the ListGamesResponse with a mutex for thread-safety.
type PrematchData struct {
	mu         sync.RWMutex
	Prematches *pb.ListPrematchResponse
}

// NewGamesData creates a new instance of GamesData.
func NewPrematchData() *PrematchData {
	return &PrematchData{
		Prematches: &pb.ListPrematchResponse{},
	}
}

// UpdateGames safely updates the games data.
func (gd *PrematchData) UpdatePrematchData(prematches *pb.ListPrematchResponse) {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.Prematches = prematches
}

// GetGames safely retrieves the games data.
func (gd *PrematchData) GetPrematchData() *pb.ListPrematchResponse {
	gd.mu.RLock()
	defer gd.mu.RUnlock()
	// Make a deep copy of Games to ensure the data is not modified while being read
	gamesCopy := *gd.Prematches
	return &gamesCopy
}
