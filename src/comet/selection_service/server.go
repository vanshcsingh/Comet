package selection_service

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":424242"
)

// Server is the MSL server
type Server struct {
}

// Query takes in a context hash and a set of features, it returns a query ID
func (s *Server) Query(ctxh string, xDim int32, yDim int32, type string, features []int32) (string, error) {
	return nil, nil
}

// Feedback takes in a context hash, query Id and a score
func (s *Server) Feedback(ctxh string, queryID string, evaluation int32) (error) {
	return nil
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
