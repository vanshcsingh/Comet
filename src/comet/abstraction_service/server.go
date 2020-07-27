package abstraction_service

import (
	"context"
	"log"
	"net"

	"comet/abstraction_service/pb"
	"google.golang.org/grpc"
)

const (
	port = ":424242"
)

// Server is the MAL server
type Server struct {
	pb.UnimplementedAbstractionServiceServer
}

// Predict takes in a context and predict request
func (s *Server) Predict(ctx context.Context, pr *pb.PredictRequest) (*pb.PredictReply, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAbstractionServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
