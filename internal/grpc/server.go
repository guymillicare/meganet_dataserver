package grpc

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSportsbookServiceServer
	prematchData *PrematchData
	lock         sync.Mutex // Protects the clients map
	clients      map[string]chan *pb.LiveOddsData
	oddsChannel  chan *pb.LiveOddsData
}

// ListPrematch returns a prematch data response.
func (s *server) ListPrematch(ctx context.Context, req *pb.ListPrematchRequest) (*pb.ListPrematchResponse, error) {
	if s.prematchData == nil || s.prematchData.GetPrematchData() == nil {
		return nil, nil
	}
	return s.prematchData.GetPrematchData(), nil
}

// SendLiveOdds streams live odds data to clients.
// func (s *server) SendLiveOdds(req *pb.LiveOddsRequest, stream pb.SportsbookService_SendLiveOddsServer) error {
// 	ch := make(chan *pb.LiveOddsData, 10)
// 	clientID := "unique-client-id" // Generate or manage client ID appropriately
// 	s.lock.Lock()
// 	s.clients[clientID] = ch
// 	s.lock.Unlock()

// 	defer func() {
// 		s.lock.Lock()
// 		delete(s.clients, clientID)
// 		close(ch)
// 		s.lock.Unlock()
// 	}()

// 	for data := range ch {
// 		if err := stream.Send(data); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// BroadcastOddsData sends live odds to all connected clients.
func (s *server) SendLiveOdds(req *pb.LiveOddsRequest, stream pb.SportsbookService_SendLiveOddsServer) error {
	for oddsData := range s.oddsChannel {
		s.lock.Lock()
		if err := stream.Send(oddsData); err != nil {
			return err
		}
		s.lock.Unlock()
	}
	return nil
}

// StartGRPCServer initializes and starts the gRPC server.
func StartGRPCServer(port string, prematchData *PrematchData, oddsChannel chan *pb.LiveOddsData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sportsbookServer := &server{
		prematchData: prematchData,
		oddsChannel:  oddsChannel,
	}
	pb.RegisterSportsbookServiceServer(grpcServer, sportsbookServer)
	// go sportsbookServer.BroadcastOddsData()

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// Example usage: sportsbookServer.BroadcastOddsData(&pb.LiveOddsData{})
}

func ListenToStream(url string, oddsChannel chan *pb.LiveOddsData, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: {") {
			continue
		}

		jsonStr := strings.TrimPrefix(line, "data: ")
		var oddsData types.OddsStream
		if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Convert types.OddsStream to pb.LiveOddsData as needed before sending
		convertedOdds := &pb.Data{
			BetName:         oddsData.Data[0].BetName,
			BetPoints:       oddsData.Data[0].BetPoints,
			BetPrice:        oddsData.Data[0].BetPrice,
			BetType:         oddsData.Data[0].BetType,
			GameId:          oddsData.Data[0].GameId,
			Id:              oddsData.Data[0].Id,
			IsLive:          oddsData.Data[0].IsLive,
			IsMain:          oddsData.Data[0].IsMain,
			League:          oddsData.Data[0].League,
			PlayerId:        oddsData.Data[0].PlayerId,
			Selection:       oddsData.Data[0].Selection,
			SelectionLine:   oddsData.Data[0].SelectionLine,
			SelectionPoints: oddsData.Data[0].SelectionPoints,
			Sport:           oddsData.Data[0].Sport,
			Sportsbook:      oddsData.Data[0].Sportsbook,
			Timestamp:       oddsData.Data[0].Timestamp,
		}
		convertedOddsData := &pb.LiveOddsData{
			EntryId: oddsData.EntryId,
			Type:    oddsData.Type,
			Data:    convertedOdds,
		} // Conversion logic here
		oddsChannel <- convertedOddsData
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stream: %v", err)
	}
}
