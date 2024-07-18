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
	"time"

	pb "sportsbook-backend/internal/proto"
	"sportsbook-backend/internal/types"
	"sportsbook-backend/pkg/queue"

	"google.golang.org/grpc"
)

type clientOddsStream struct {
	stream   pb.SportsbookService_SendLiveOddsServer
	stopChan chan struct{}
}

type clientScoreStream struct {
	stream   pb.SportsbookService_SendLiveScoreServer
	stopChan chan struct{}
}

type server struct {
	pb.UnimplementedSportsbookServiceServer
	lock         sync.Mutex
	oddsClients  map[string]*clientOddsStream
	scoreClients map[string]*clientScoreStream
	oddsChannel  chan *pb.LiveOddsData
	scoreChannel chan *pb.LiveScoreData
}

// Start a separate goroutine to handle broadcasting data to clients
func (s *server) broadcastOddsData() {
	for oddsData := range s.oddsChannel {
		s.lock.Lock()
		for id, client := range s.oddsClients {
			select {
			case <-client.stopChan:
				// Client is done, remove it
				delete(s.oddsClients, id)
			default:
				if err := client.stream.Send(oddsData); err != nil {
					// Error sending to client, remove it
					close(client.stopChan)
					delete(s.oddsClients, id)
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
	s.oddsClients[clientID] = &clientOddsStream{stream: stream, stopChan: stopChan}
	s.lock.Unlock()

	// Wait for client to disconnect
	<-stream.Context().Done()

	s.lock.Lock()
	if s.oddsClients[clientID] != nil {
		fmt.Print("client", clientID)
		fmt.Print("length", len(s.oddsClients))
		close(s.oddsClients[clientID].stopChan)
		delete(s.oddsClients, clientID)
	}
	s.lock.Unlock()

	return nil
}

// Start a separate goroutine to handle broadcasting data to clients
func (s *server) broadcastScoreData() {
	for scoreData := range s.scoreChannel {
		s.lock.Lock()
		for id, client := range s.scoreClients {
			select {
			case <-client.stopChan:
				// Client is done, remove it
				delete(s.scoreClients, id)
			default:
				if err := client.stream.Send(scoreData); err != nil {
					// Error sending to client, remove it
					close(client.stopChan)
					delete(s.scoreClients, id)
				}
			}
		}
		s.lock.Unlock()
	}
}

func (s *server) SendLiveScore(req *pb.LiveScoreRequest, stream pb.SportsbookService_SendLiveScoreServer) error {
	clientID := fmt.Sprintf("%p", stream) // Unique client ID
	stopChan := make(chan struct{})

	s.lock.Lock()
	s.scoreClients[clientID] = &clientScoreStream{stream: stream, stopChan: stopChan}
	s.lock.Unlock()

	// Wait for client to disconnect
	<-stream.Context().Done()

	s.lock.Lock()
	if s.scoreClients[clientID] != nil {
		fmt.Print("client", clientID)
		fmt.Print("length", len(s.scoreClients))
		close(s.scoreClients[clientID].stopChan)
		delete(s.scoreClients, clientID)
	}
	s.lock.Unlock()

	return nil
}
func StartGRPCServer(port string, oddsChannel chan *pb.LiveOddsData, scoreChannel chan *pb.LiveScoreData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	sportsbookServer := &server{
		oddsClients:  make(map[string]*clientOddsStream),
		scoreClients: make(map[string]*clientScoreStream),
		oddsChannel:  oddsChannel,
		scoreChannel: scoreChannel,
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

func ListenToOddsStream(url string, oddsChannel chan *pb.LiveOddsData, wg *sync.WaitGroup, rabbitMQ *queue.RabbitMQ) {
	defer wg.Done()

	for {
		err := connectAndListen(url, oddsChannel, rabbitMQ)
		if err != nil {
			// log.Printf("Error in ListenToOddsStream: %v. Reconnecting in 5 seconds...", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func connectAndListen(url string, oddsChannel chan *pb.LiveOddsData, rabbitMQ *queue.RabbitMQ) error {
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a timeout for the HTTP client
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	// fmt.Printf("ListenToOddsStream: %s\n", "Started")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading line: %w", err)
		}
		if !strings.HasPrefix(line, "data: {") {
			// fmt.Printf("ListenToOddsStream LINE: %s\n", line)
			continue
		}

		jsonStr := strings.TrimPrefix(line, "data: ")
		var oddsData types.OddsStream
		if err := json.Unmarshal([]byte(jsonStr), &oddsData); err != nil {
			// log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}
		// fmt.Printf("ListenToOddsStream: %s\n", jsonStr)

		// Send to RabbitMQ
		err = rabbitMQ.Publish([]byte(jsonStr))
		if err != nil {
			log.Printf("Error publishing to RabbitMQ: %v", err)
		}
	}
}

func ListenToScoreStream(url string, scoreChannel chan *pb.LiveScoreData, wg *sync.WaitGroup, rabbitMQ *queue.RabbitMQ) {
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
		var scoreData types.ScoreStream
		if err := json.Unmarshal([]byte(jsonStr), &scoreData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// fmt.Printf("Stream: %s\n", jsonStr)
		rabbitMQ.Publish([]byte(jsonStr))
	}
}
