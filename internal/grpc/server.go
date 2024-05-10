package grpc

import (
	"context"
	"log"
	"net"

	pb "sportsbook-backend/internal/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSportsbookServiceServer
	prematchData *PrematchData
	clients      map[string]chan *pb.LiveOddsData
}

func (s *server) ListPrematch(ctx context.Context, req *pb.ListPrematchRequest) (*pb.ListPrematchResponse, error) {
	if s.prematchData == nil || s.prematchData.GetPrematchData() == nil {
		return nil, nil
	}
	return s.prematchData.GetPrematchData(), nil
}

func (s *server) SendLiveOdds(req *pb.LiveOddsRequest, stream pb.SportsbookService_SendLiveOddsServer) error {
	ch := make(chan *pb.LiveOddsData, 10)
	clientID := "unique-client-id" // Generate or manage client ID appropriately
	s.clients[clientID] = ch

	defer func() {
		delete(s.clients, clientID)
		close(ch)
	}()

	for data := range ch {
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func StartGRPCServer(port string, prematchData *PrematchData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSportsbookServiceServer(grpcServer, &server{
		prematchData: prematchData,
		clients:      make(map[string]chan *pb.LiveOddsData)})

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
