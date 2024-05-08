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
}

func (s *server) ListPrematch(ctx context.Context, req *pb.ListPrematchRequest) (*pb.ListPrematchResponse, error) {
	if s.prematchData == nil || s.prematchData.GetPrematchData() == nil {
		return nil, nil
	}
	return s.prematchData.GetPrematchData(), nil
}

func StartGRPCServer(port string, prematchData *PrematchData) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSportsbookServiceServer(grpcServer, &server{prematchData: prematchData})

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
