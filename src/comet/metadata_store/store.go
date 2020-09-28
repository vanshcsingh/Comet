package metadata_store

import (
	"comet"

	modelpb "comet/models/container_models"

	"encoding/json"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
)

// MetadataStore provides us the model address for a given modelID
type MetadataStore interface {
	GetMIDs() []comet.ModelIDType
	GetEntry(comet.ModelIDType) ModelEntry
	GetClient(comet.ModelIDType) (modelpb.ServiceClient, error)

	LabelToVector(string) ([]float64, error)
	VectorToLabel([]float64) (string, error)

	GetNumModels() int
	GetNumLabels() int
	GetSelectionGamma() float64
	GetAbstractionServiceAddr() string
	GetSelectionServiceAddr() string
}

// LocalMdStoreSettings struct contains a list of ModelEntry objects
type LocalMdStoreSettings struct {
	Entries   []ModelEntry    `json:"model_entries"`
	Addresses ServiceAddress  `json:"service_addresses"`
	MSLParams SelectionParams `json:"selection_params"`
}

// ModelEntry defines our local metadata json file entries
type ModelEntry struct {
	ModelID         comet.ModelIDType `json:"model_id"`
	Addr            string            `json:"addr"`
	Desc            string            `json:"desc"`
	DockerContainer string            `json:"docker_container"`
}

// ServiceAddress struct stores the addresses for MAL and MSL
type ServiceAddress struct {
	AbstractionService string `json:"abstraction_service"`
	SelectionService   string `json:"selection_service"`
}

// SelectionParams stores the labels and exploration factor for MSL
type SelectionParams struct {
	Labels []string `json:"labels`
	Gamma  float64  `json:"gamma"`
}

// LocalFileBasedMetadataStore loads metadata information from a json file
type LocalFileBasedMetadataStore struct {
	clientConnMap map[comet.ModelIDType]grpc.ClientConnInterface
	entryMap      map[comet.ModelIDType]ModelEntry
	midList       []comet.ModelIDType

	abstractionServiceAddr string
	selectionServiceAddr   string
	selectionParams        SelectionParams
}

// instance for singleton
var instance MetadataStore

func GetMetadataStoreInstance() MetadataStore {
	if instance == nil {
		instance = createLocalFileBasedMetadataStore("./metadata_store/tmp.json")
	}
	return instance
}

// createLocalFileBasedMetadataStore creates a MetadataStore object from a json file
func createLocalFileBasedMetadataStore(metadataFile string) MetadataStore {

	// Load file contents. Panic if error
	dat, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		panic(err)
	}

	// unmarshal file contents into LocalMdStoreSettings
	var entries LocalMdStoreSettings
	if err = json.Unmarshal(dat, &entries); err != nil {
		panic(err)
	}

	log.Println("[Store:createLocalFileBasedMetadataStore] parsed metadata file into map: ", entries)

	// initialize data structures
	clientConnMap := make(map[comet.ModelIDType]grpc.ClientConnInterface)
	entryMap := make(map[comet.ModelIDType]ModelEntry)
	midList := make([]comet.ModelIDType, len(entries.Entries))

	// create addressMap
	for i, entry := range entries.Entries {
		entryMap[entry.ModelID] = entry
		midList[i] = entry.ModelID
	}

	return &LocalFileBasedMetadataStore{
		clientConnMap:          clientConnMap,
		entryMap:               entryMap,
		midList:                midList,
		abstractionServiceAddr: entries.Addresses.AbstractionService,
		selectionServiceAddr:   entries.Addresses.SelectionService,
		selectionParams:        entries.MSLParams,
	}
}
