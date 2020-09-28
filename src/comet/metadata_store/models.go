package metadata_store

import (
	"comet"
	modelpb "comet/models/container_models"
	"log"

	"google.golang.org/grpc"
)

// GetMIDs returns list of ModelIDs that we store
func (md *LocalFileBasedMetadataStore) GetMIDs() []comet.ModelIDType {
	return md.midList
}

// GetEntry gets a model entry for a ModelID
func (md *LocalFileBasedMetadataStore) GetEntry(mid comet.ModelIDType) ModelEntry {
	return md.entryMap[mid]
}

// GetClient returns a grpc client for the Model
func (md *LocalFileBasedMetadataStore) GetClient(mid comet.ModelIDType) (modelpb.ServiceClient, error) {
	conn, exists := md.clientConnMap[mid]

	log.Println("[Store:GetClient] does a connection exist so far? ", exists)

	var err error
	if !exists {
		entry := md.GetEntry(mid)

		log.Println("[Store:GetClient] trying to create a connection for: ", entry.Addr)
		conn, err = grpc.Dial(entry.Addr, grpc.WithInsecure(), grpc.WithBlock())

		log.Println("[Store:GetClient] connection error: ", err)
	}
	return modelpb.NewServiceClient(conn), err
}
