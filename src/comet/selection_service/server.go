package selection_service

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "comet/selection_service/pb"
)

const (
	port = ":424242"
)

// Server is the MSL server
type Server struct {
	pb.SelectionServiceServer
}

// Query takes in a context hash and a set of features, it returns a query ID
func (s *Server) Query(context context.Context, query *pb.QueryRequest) (*pb.QueryReply, error) {
	return nil, nil
}

// Feedback takes in a context hash, query Id and a score
func (s *Server) Feedback(context context.Context, request *pb.FeedbackRequest) (*pb.FeedbackReply, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSelectionServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
