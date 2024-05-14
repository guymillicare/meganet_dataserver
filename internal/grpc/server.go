package grpc

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/repositories"
	"sportsbook-backend/internal/types"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSportsbookServiceServer
	lock        sync.Mutex // Protects the clients map
	clients     map[string]chan *pb.LiveOddsData
	oddsChannel chan *pb.LiveOddsData
}

// ListPrematch returns a prematch data response.
// func (s *server) ListPrematch(ctx context.Context, req *pb.ListPrematchRequest) (*pb.ListPrematchResponse, error) {
// 	if s.prematchData == nil || s.prematchData.GetPrematchData() == nil {
// 		return nil, nil
// 	}
// 	return s.prematchData.GetPrematchData(), nil
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
func StartGRPCServer(port string, oddsChannel chan *pb.LiveOddsData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sportsbookServer := &server{
		// prematchData: prematchData,
		oddsChannel: oddsChannel,
	}
	pb.RegisterSportsbookServiceServer(grpcServer, sportsbookServer)
	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	log.Printf("gRPC server started")
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

		sportEvent, _ := repositories.GetSportEventFromRedis(oddsData.Data[0].GameId)
		marketConstant, _ := repositories.GetMarketConstantFromRedis(oddsData.Data[0].BetType)
		outcome := &types.OutcomeItem{
			ReferenceId: oddsData.Data[0].BetType + ":" + oddsData.Data[0].BetName,
			EventId:     sportEvent.Id,
			MarketId:    marketConstant.Id,
			Name:        oddsData.Data[0].BetName,
			Odds:        oddsData.Data[0].BetPrice,
			Active:      oddsData.Type == "updated",
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		repositories.SaveOutcomeToRedis(outcome)
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
