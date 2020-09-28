package main

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
	batchThreshold = 1
	cacheSize      = 3000
	duration       = time.Hour * 1
)

// Server is the MAL server
type Server struct {
	pb.AbstractionServiceService
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
	cache, _ := cache.CreateAndStartLocalCache(cacheSize, predictProducer, resultConsumer)

	mdStore := md.GetMetadataStoreInstance()

	// create and start batcher service
	batch.CreateAndStartLocalBatcher(predictConsumer, resultProducer, batchThreshold, duration, mdStore)

	return &Server{
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

	addr := md.GetMetadataStoreInstance().GetAbstractionServiceAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Abstraction service has started listening on %s\n", addr)

	grpcServer := grpc.NewServer()
	abstractionService := pb.NewAbstractionServiceService(CreateServer())

	pb.RegisterAbstractionServiceService(grpcServer, abstractionService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
