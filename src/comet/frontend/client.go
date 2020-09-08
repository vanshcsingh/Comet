package main

import (
	"comet"
	"context"
	"log"
	"time"

	"comet/abstraction_service/pb"

	"google.golang.org/grpc"
)

const (
	address = "localhost:4001"
)

func main() {

	conn1, err1 := grpc.Dial("localhost:4000", grpc.WithInsecure(), grpc.WithBlock())
	if err1 != nil {
		log.Fatalf("Could not connect to localhost:4000", err1)
	}
	log.Printf("[Client] connection successful to localhost:4000")
	conn1.Close()

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAbstractionServiceClient(conn)
	log.Printf("[Client] Created network connection with %v", address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	image := make(comet.ImageVectorType, 784)

	r, err := c.Predict(ctx, &pb.PredictRequest{
		ImageVector: image,
		ModelId:     0,
		ContextUuid: "0xdeadbeef",
	})
	if err != nil {
		log.Fatalf("could not predict: %v", err)
	}

	log.Printf("Prediction resulted in label: %s", r.GetLabel())
}
