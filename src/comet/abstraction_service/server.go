package abstraction_service

import (
	"comet"
	"context"
	"log"
	"net"
	"time"

	"comet/abstraction_service/batch"
	"comet/abstraction_service/cache"
	"comet/abstraction_service/pb"

	"google.golang.org/grpc"
)

const (
	port = ":424242"
	cacheSize = 3000
	duration = time.Second * 1
)

// Server is the MAL server
type Server struct {
	pb.UnimplementedAbstractionServiceServer
	cache cache.MALCache
	batcher batch.Service
}

// CreateServer creates the receiver server
func CreateServer(cacheSize int, batchingPeriod time.Duration) *Server {
	predictPipe := make(chan *comet.PredictParams)
	predictConsumer := &batch.LocalPredictConsumer{predictPipe}
	predictProducer := &batch.LocalPredictProducer{predictPipe}

	resultPipe := make(chan *comet.PredictResult)
	resultConsumer := &batch.LocalResultConsumer{resultPipe}
	resultProducer := &batch.LocalResultProducer{resultPipe}

	return &Server {
		cache: cache.CreateLocalCache(cacheSize, predictProducer, resultConsumer),
		batcher: batch.CreateLocalBatcher(predictConsumer, resultProducer, batchingPeriod),
	}
}

// Predict is a synchronous call that takes in a context and predict request
func (s *Server) Predict(ctx context.Context, pr *pb.PredictRequest) (*pb.PredictReply, error) {
	predictParams := comet.CreatePredictParamsMAL(pr)
	return &pb.PredictReply{
		Label: s.cache.Request(predictParams),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAbstractionServiceServer(s, CreateServer(cacheSize, duration))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
