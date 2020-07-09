package abstraction_service

import (
	"context"
	"io"
	"log"
	"net"

	"comet/abstraction_service/pb"
	"google.golang.org/grpc"
)

const (
	port = ":424242"
)

type MAL_Server struct {
	pb.UnimplementedAbstractionServiceServer
}

func (s *MAL_Server) PredictSetup(ctx context.Context, in *pb.PredictSetupRequest) *pb.PredictSetupReply {
	log.Printf("Received: %v", in)

	return &pb.PredictSetupReply{}
}

func (s *MAL_Server) Predict(ctx context.Context, stream pb.AbstractionService_PredictServer) {
	for {
		pixel, err := stream.Recv()
		if err == io.EOF {
			log.Printf("End of stream")
		}

		log.Printf("Recieved: %v", pixel)
	}
}


func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAbstractionServiceServer(s, &MAL_Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
