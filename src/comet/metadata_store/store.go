package metadata_store

import (
	"comet"

	"encoding/json"
	"io/ioutil"
)

// MetadataStore provides us the model address for a given modelID
type MetadataStore interface {
	GetMIDs() []comet.ModelIDType
	GetEntry(comet.ModelIDType) ModelEntry
}

// LocalModelEntries struct contains a list of ModelEntry objects
type LocalModelEntries struct {
	Entries []ModelEntry `json:"model_entries"`
}

// ModelEntry defines our local metadata json file entries
type ModelEntry struct {
	ModelID comet.ModelIDType `json:"model_id"`
	Addr string `json:"addr"`
	Desc string `json:"desc"`
}

// LocalFileBasedMetadataStore loads metadata information from a json file
type LocalFileBasedMetadataStore struct {
	midList []comet.ModelIDType
	entryMap map[comet.ModelIDType] ModelEntry
}

// CreateLocalFileBasedMetadataStore creates a MetadataStore object from a json file
func CreateLocalFileBasedMetadataStore(metadataFile string) MetadataStore {

	// Load file contents. Panic if error 
	dat, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		panic (err)
	}

	// unmarshal file contents into LocalModelEntries
	var entries LocalModelEntries
	if err = json.Unmarshal(dat, &entries); err != nil {
		panic (err)
	}

	// list of mids
	midList := make([]comet.ModelIDType, len(entries.Entries))

	// create addressMap
	var entryMap map[comet.ModelIDType] ModelEntry
	for i, entry := range entries.Entries {
		entryMap[entry.ModelID] = entry
		midList[i] = entry.ModelID
	}


	return &LocalFileBasedMetadataStore {
		entryMap: entryMap,
		midList: midList,
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