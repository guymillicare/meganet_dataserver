package grpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/pkg/queue"

	"google.golang.org/grpc"
)

type clientStream struct {
	stream   pb.SportsbookService_SendLiveOddsServer
	stopChan chan struct{}
}

type server struct {
	pb.UnimplementedSportsbookServiceServer
	lock        sync.Mutex
	clients     map[string]*clientStream
	oddsChannel chan *pb.LiveOddsData
}

// Start a separate goroutine to handle broadcasting data to clients
func (s *server) broadcastOddsData() {
	for oddsData := range s.oddsChannel {
		s.lock.Lock()
		for id, client := range s.clients {
			select {
			case <-client.stopChan:
				// Client is done, remove it
				delete(s.clients, id)
			default:
				if err := client.stream.Send(oddsData); err != nil {
					// Error sending to client, remove it
					close(client.stopChan)
					delete(s.clients, id)
				}
			}
		}
		s.lock.Unlock()
	}
}

func (s *server) SendLiveOdds(req *pb.LiveOddsRequest, stream pb.SportsbookService_SendLiveOddsServer) error {
	clientID := fmt.Sprintf("%p", stream) // Unique client ID
	stopChan := make(chan struct{})

	s.lock.Lock()
	s.clients[clientID] = &clientStream{stream: stream, stopChan: stopChan}
	s.lock.Unlock()

	// Wait for client to disconnect
	<-stream.Context().Done()

	s.lock.Lock()
	if s.clients[clientID] != nil {
		fmt.Print("client", clientID)
		fmt.Print("length", len(s.clients))
		close(s.clients[clientID].stopChan)
		delete(s.clients, clientID)
	}
	s.lock.Unlock()

	return nil
}

func StartGRPCServer(port string, oddsChannel chan *pb.LiveOddsData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sportsbookServer := &server{
		clients:     make(map[string]*clientStream),
		oddsChannel: oddsChannel,
	}
	pb.RegisterSportsbookServiceServer(grpcServer, sportsbookServer)
	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start broadcasting odds data to all connected clients
	go sportsbookServer.broadcastOddsData()

	log.Printf("gRPC server started")
}

func ListenToStream(url string, oddsChannel chan *pb.LiveOddsData, wg *sync.WaitGroup, rabbitMQ *queue.RabbitMQ) {
	defer wg.Done()

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadString('\n')
		if !strings.HasPrefix(line, "data: {") {
			continue
		}
		jsonStr := strings.TrimPrefix(line, "data: ")
		var oddsData types.OddsStream
		if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// fmt.Printf("Stream: %s\n", jsonStr)
		// for _, odds := range oddsData.Data {
		// 	sportEvent, _ := repositories.GetSportEventFromRedis(odds.GameId)
		// 	marketConstant, _ := repositories.GetMarketConstantFromRedis(odds.BetType)
		// 	if marketConstant != nil && sportEvent != nil {
		// 		outcome := &types.OutcomeItem{
		// 			ReferenceId: odds.BetType + ":" + odds.BetName,
		// 			EventId:     sportEvent.Id,
		// 			MarketId:    marketConstant.Id,
		// 			Name:        odds.BetName,
		// 			Odds:        odds.BetPrice,
		// 			Active:      oddsData.Type == "odds",
		// 			CreatedAt:   time.Now().UTC(),
		// 			UpdatedAt:   time.Now().UTC(),
		// 		}
		// 		repositories.SaveOutcomeToRedis(outcome)
		// 		fmt.Printf("Outcome: %v\n", outcome)

		// 		convertedOdds := &pb.Data{
		// 			BetName:         odds.BetName,
		// 			BetPoints:       odds.BetPoints,
		// 			BetPrice:        odds.BetPrice,
		// 			BetType:         odds.BetType,
		// 			GameId:          odds.GameId,
		// 			Id:              odds.Id,
		// 			IsLive:          odds.IsLive,
		// 			IsMain:          odds.IsMain,
		// 			League:          odds.League,
		// 			PlayerId:        odds.PlayerId,
		// 			Selection:       odds.Selection,
		// 			SelectionLine:   odds.SelectionLine,
		// 			SelectionPoints: odds.SelectionPoints,
		// 			Sport:           odds.Sport,
		// 			Sportsbook:      odds.Sportsbook,
		// 			Timestamp:       odds.Timestamp,
		// 		}

		// 		convertedOddsData := &pb.LiveOddsData{
		// 			EntryId: oddsData.EntryId,
		// 			Type:    oddsData.Type,
		// 			Data:    convertedOdds,
		// 		}
		// 		// Send live data to gRPC clients
		// 		oddsChannel <- convertedOddsData
		// 		// fmt.Printf("Stream: %s\n", jsonStr)
		// 		fmt.Printf("Consumer: %v\n", convertedOddsData)
		// 	}
		// }

		rabbitMQ.Publish([]byte(jsonStr))
	}
}
