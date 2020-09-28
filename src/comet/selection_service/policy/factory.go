package policy

import (
	malpb "comet/abstraction_service/pb"
	md "comet/metadata_store"
	"log"

	"google.golang.org/grpc"
)

// GenerateExp3  currently each policy object creates it's own file descriptor
// for an abstraction-service connection. this can cause a problem where our file descriptors blow up.
func GenerateExp3() SingleSelectionPolicy {
	mdStore := md.GetMetadataStoreInstance()

	malConn, err := grpc.Dial(mdStore.GetAbstractionServiceAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	malClient := malpb.NewAbstractionServiceClient(malConn)

	return CreateExp3(
		mdStore.GetSelectionGamma(),
		mdStore.GetNumModels(),
		malClient,
	)
}

// GenerateExp4  currently each policy object creates it's own file descriptor
// for an abstraction-service connection. this can cause a problem where our file descriptors blow up.
func GenerateExp4() EnsembleSelectionPolicy {
	mdStore := md.GetMetadataStoreInstance()

	malConn, err := grpc.Dial(mdStore.GetAbstractionServiceAddr(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	malClient := malpb.NewAbstractionServiceClient(malConn)

	return CreateExp4(
		mdStore.GetSelectionGamma(),
		mdStore.GetNumLabels(),
		mdStore.GetNumModels(),
		malClient,
	)
}
