package main

import (
	"comet"
	"context"
	"crypto/rand"
	"log"
	"net"

	md "comet/metadata_store"
	pb "comet/selection_service/pb"
	"comet/selection_service/policy"
	"comet/selection_service/sessiondb"

	"google.golang.org/grpc"
)

const (
	queryIDLength int = 100
)

// Server is the MSL server
type Server struct {
	pb.SelectionServiceService

	db sessiondb.SessionDB
}

// InitSelectionService initializes a selection service server
func InitSelectionService() *Server {
	// create local db
	localdb := sessiondb.CreateLocalDB()

	return &Server{
		db: localdb,
	}
}

func generateUniqueID(bytelength int) string {
	byteArray := make([]byte, bytelength)
	_, err := rand.Read(byteArray)
	if err != nil {
		panic(err)
	}

	return string(byteArray)
}

// Query takes in a context hash and a set of features, it returns a query ID
func (s *Server) Query(ctx context.Context, query *pb.QueryRequest) (*pb.QueryReply, error) {

	contextuuid := query.GetContextUuid()
	fType := query.GetFType()
	imageVector := comet.ImageVectorType(query.GetImageVector())

	// generate a randomized unique queryID
	queryID := generateUniqueID(queryIDLength)

	var predictionLabel string

	if fType == sessiondb.SingleSelectionType {
		selectionPolicy, err := s.db.GetSingleSelectionPolicy(contextuuid)
		if err != nil {
			log.Printf("[SelectionService: Query]: found no single selection policy for %s: %v\n", contextuuid, err)
			selectionPolicy = policy.GenerateExp3()

			s.db.SetSingleSelectionPolicy(contextuuid, selectionPolicy)
		}

		selection := selectionPolicy.Select(ctx, contextuuid, imageVector)

		// update db with selection
		s.db.SetType(queryID, sessiondb.SingleSelectionType)
		s.db.SetSingleSelection(queryID, selection)

		predictionLabel = selection.PredictionLabel
	} else {
		selectionPolicy, err := s.db.GetEnsembleSelectionPolicy(contextuuid)
		if err != nil {
			log.Printf("[SelectionService: Query]: found no ensemble selection policy for %s: %v\n", contextuuid, err)
			selectionPolicy = policy.GenerateExp4()

			s.db.SetEnsembleSelectionPolicy(contextuuid, selectionPolicy)
		}

		selection := selectionPolicy.Select(ctx, contextuuid, imageVector)

		// update db with selection
		s.db.SetType(queryID, sessiondb.EnsembleSelectionType)
		s.db.SetEnsembleSelection(queryID, selection)

		predictionLabel = selection.PredictionLabel
	}

	return &pb.QueryReply{
		Label:   predictionLabel,
		QueryID: queryID,
	}, nil
}

// Feedback takes in a context hash, query Id and a score
func (s *Server) Feedback(context context.Context, request *pb.FeedbackRequest) (*pb.FeedbackReply, error) {
	contextuuid := request.GetContextUuid()
	queryID := request.GetQueryID()

	fType, err := s.db.GetType(queryID)
	if err != nil {
		panic(err)
	}

	// error in finding selection object or selectionpolicy
	var sErr, pErr error

	if fType == sessiondb.SingleSelectionType {
		var selectionPolicy policy.SingleSelectionPolicy
		var selection *policy.SingleSelection

		if selectionPolicy, pErr = s.db.GetSingleSelectionPolicy(contextuuid); pErr != nil {
			panic(pErr)
		}
		if selection, sErr = s.db.GetSingleSelection(queryID); sErr != nil {
			panic(sErr)
		}

		selectionPolicy.Feedback(selection, request.GetEvaluation())
	} else {

		var selectionPolicy policy.EnsembleSelectionPolicy
		var selection *policy.EnsembleSelection

		if selectionPolicy, pErr = s.db.GetEnsembleSelectionPolicy(contextuuid); pErr != nil {
			panic(pErr)
		}
		if selection, sErr = s.db.GetEnsembleSelection(queryID); sErr != nil {
			panic(sErr)
		}

		selectionPolicy.Feedback(selection, request.GetEvaluation())
	}

	return &pb.FeedbackReply{}, nil
}

func main() {
	addr := md.GetMetadataStoreInstance().GetSelectionServiceAddr()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	server := InitSelectionService()
	selectionService := pb.NewSelectionServiceService(&server)

	pb.RegisterSelectionServiceService(grpcServer, selectionService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
