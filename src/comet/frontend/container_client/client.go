package main

import (
	// "comet"
	pb "comet/models/container_models"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:4000"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewServiceClient(conn)
	log.Printf("[Client] Created network connection with %v", address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	imageVector := &pb.ImageVector{
		Pixels: make([]int32, 784),
	}

	request := pb.PredictRequest{
		Images: []*pb.ImageVector{imageVector, imageVector},
	}

	r, err := c.Predict(ctx, &request)
	if err != nil {
		log.Fatalf("could not predict: %v", err)
	}

	log.Printf("Prediction resulted in labels: %v", r.GetLabels())
}
