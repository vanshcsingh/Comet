package abstraction_service

import (
	"comet"
	"context"
	"log"
	"net"
	"time"

	"comet/abstraction_service/batch"
	"comet/abstraction_service/batch/mq"
	"comet/abstraction_service/cache"
	"comet/abstraction_service/pb"

	md "comet/metadata_store"

	"google.golang.org/grpc"
)

const (
	batchThreshold = 10
	cacheSize = 3000
	duration = time.Second * 1
	port = ":424242"
)

// Server is the MAL server
type Server struct {
	pb.UnimplementedAbstractionServiceServer
	cache cache.MALCache
}

// CreateServer creates the receiver server
func CreateServer() *Server {
	predictPipe := make(chan *comet.PredictParams)
	predictConsumer := &mq.LocalPredictConsumer{Pipe: predictPipe}
	predictProducer := &mq.LocalPredictProducer{Pipe: predictPipe}

	resultPipe := make(chan *comet.PredictResult)
	resultConsumer := &mq.LocalResultConsumer{Pipe: resultPipe}
	resultProducer := &mq.LocalResultProducer{Pipe: resultPipe}

	// create and start caching service
	cache := cache.CreateAndStartLocalCache(cacheSize, predictProducer, resultConsumer)

	mdStore := md.CreateLocalFileBasedMetadataStore("./metadata_store/tmp.json")

	// create and start batcher service
	batch.CreateAndStartLocalBatcher(predictConsumer, resultProducer, batchThreshold, duration, mdStore)

	return &Server {
		cache: cache,
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
	pb.RegisterAbstractionServiceServer(s, CreateServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
