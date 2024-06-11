package grpc

import (
	"bufio"
	"encoding/json"
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

func ListenToStream(url string, oddsChannel chan *pb.LiveOddsData, wg *sync.WaitGroup, rabbitMQ *queue.RabbitMQ) {
	defer wg.Done()

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// scanner := bufio.NewScanner(resp.Body)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	if !strings.HasPrefix(line, "data: {") {
	// 		continue
	// 	}

	// 	jsonStr := strings.TrimPrefix(line, "data: ")
	// 	var oddsData types.OddsStream
	// 	if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
	// 		log.Printf("Error unmarshaling JSON: %v", err)
	// 		continue
	// 	}

	// 	if err == nil {
	// 		rabbitMQ.Publish([]byte(jsonStr))
	// 	}

	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatalf("Error reading from stream: %v", err)
	// }

	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadString('\n')
		if !strings.HasPrefix(line, "data: {") {
			continue
		}
		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	log.Fatalf("Error reading from stream: %v", err)
		// }

		jsonStr := strings.TrimPrefix(line, "data: ")
		var oddsData types.OddsStream
		if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		rabbitMQ.Publish([]byte(jsonStr))
	}
}
