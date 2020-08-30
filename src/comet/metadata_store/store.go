package metadata_store

import (
	"comet"

	modelpb "comet/models/pb"

	"encoding/json"
	"io/ioutil"

	"google.golang.org/grpc"
)

// MetadataStore provides us the model address for a given modelID
type MetadataStore interface {
	GetMIDs() []comet.ModelIDType
	GetEntry(comet.ModelIDType) ModelEntry
	GetClient(comet.ModelIDType) (modelpb.ServiceClient, error)
}

// LocalModelEntries struct contains a list of ModelEntry objects
type LocalModelEntries struct {
	Entries []ModelEntry `json:"model_entries"`
}

// ModelEntry defines our local metadata json file entries
type ModelEntry struct {
	ModelID         comet.ModelIDType `json:"model_id"`
	Addr            string            `json:"addr"`
	Desc            string            `json:"desc"`
	DockerContainer string            `json:"docker_container"`
}

// LocalFileBasedMetadataStore loads metadata information from a json file
type LocalFileBasedMetadataStore struct {
	clientConnMap map[comet.ModelIDType]grpc.ClientConnInterface
	midList       []comet.ModelIDType
	entryMap      map[comet.ModelIDType]ModelEntry
}

// CreateLocalFileBasedMetadataStore creates a MetadataStore object from a json file
func CreateLocalFileBasedMetadataStore(metadataFile string) MetadataStore {

	// Load file contents. Panic if error
	dat, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		panic(err)
	}

	// unmarshal file contents into LocalModelEntries
	var entries LocalModelEntries
	if err = json.Unmarshal(dat, &entries); err != nil {
		panic(err)
	}

	// list of mids
	midList := make([]comet.ModelIDType, len(entries.Entries))

	// create addressMap
	var entryMap map[comet.ModelIDType]ModelEntry
	for i, entry := range entries.Entries {
		entryMap[entry.ModelID] = entry
		midList[i] = entry.ModelID
	}

	return &LocalFileBasedMetadataStore{
		entryMap: entryMap,
		midList:  midList,
	}
}

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
	var err error
	if !exists {
		entry := md.GetEntry(mid)
		conn, err = grpc.Dial(entry.Addr)
	}
	return modelpb.NewServiceClient(conn), err
}
