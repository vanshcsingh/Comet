package sessiondb

import (
	"comet/selection_service/policy"
	"fmt"
)

const (
	// SingleSelectionType is the enumeration designating the exp3 type
	SingleSelectionType int = iota

	// EnsembleSelectionType is the enumeration designatin the exp4 typ
	EnsembleSelectionType
)

// SessionDB holds the weights for selection policies for contextuuids
type SessionDB interface {
	// Must be stored on a persistent datastore
	GetSingleSelectionPolicy(contextUUID string) (policy.SingleSelectionPolicy, error)
	SetSingleSelectionPolicy(contextUUID string, policy policy.SingleSelectionPolicy)

	// Must be stored on a persistent datastore
	GetEnsembleSelectionPolicy(contextUUID string) (policy.EnsembleSelectionPolicy, error)
	SetEnsembleSelectionPolicy(contextUUID string, policy policy.EnsembleSelectionPolicy)

	GetType(queryID string) int
	SetType(queryID string, fType int)

	// Stored on cache
	// for a given queryID, set/get the selectionPolicy used and the selection object
	GetSingleSelection(queryID string) (*policy.SingleSelection, error)
	SetSingleSelection(queryID string, selection *policy.SingleSelection)

	// for a given queryID, set/get the selectionPolicy used and the selection object
	GetEnsembleSelection(queryID string) (*policy.EnsembleSelection, error)
	SetEnsembleSelection(queryID string, selection *policy.EnsembleSelection)
}

// LocalDB is an in-memory implementation of SessionDB
type LocalDB struct {
	contextUUIDToSingleSelectionPolicy   map[string]policy.SingleSelectionPolicy
	contextUUIDToEnsembleSelectionPolicy map[string]policy.EnsembleSelectionPolicy

	queryIDToType map[string]int

	queryIDToSingleSelection   map[string]*policy.SingleSelection
	queryIDToEnsembleSelection map[string]*policy.EnsembleSelection
}

// CreateLocalDB instantiates a local implementation of SessionDB
func CreateLocalDB() SessionDB {
	return &LocalDB{
		contextUUIDToSingleSelectionPolicy:   make(map[string]policy.SingleSelectionPolicy),
		contextUUIDToEnsembleSelectionPolicy: make(map[string]policy.EnsembleSelectionPolicy),
		queryIDToType:                        make(map[string]int),
		queryIDToSingleSelection:             make(map[string]*policy.SingleSelection),
		queryIDToEnsembleSelection:           make(map[string]*policy.EnsembleSelection),
	}
}

// GetSingleSelectionPolicy ...
func (l *LocalDB) GetSingleSelectionPolicy(contextUUID string) (policy.SingleSelectionPolicy, error) {
	s, exists := l.contextUUIDToSingleSelectionPolicy[contextUUID]

	var err error = nil
	if !exists {
		err = fmt.Errorf("No Single Selection Policy exists for contextUUID: %s", contextUUID)
	}
	return s, err
}

// SetSingleSelectionPolicy ...
func (l *LocalDB) SetSingleSelectionPolicy(contextUUID string, policy policy.SingleSelectionPolicy) {
	l.contextUUIDToSingleSelectionPolicy[contextUUID] = policy
}

// GetEnsembleSelectionPolicy ...
func (l *LocalDB) GetEnsembleSelectionPolicy(contextUUID string) (policy.EnsembleSelectionPolicy, error) {
	s, exists := l.contextUUIDToEnsembleSelectionPolicy[contextUUID]

	var err error = nil
	if !exists {
		err = fmt.Errorf("No Ensemble Selection Policy exists for contextUUID: %s", contextUUID)
	}
	return s, err
}

// SetEnsembleSelectionPolicy ...
func (l *LocalDB) SetEnsembleSelectionPolicy(contextUUID string, policy policy.EnsembleSelectionPolicy) {
	l.contextUUIDToEnsembleSelectionPolicy[contextUUID] = policy
}

// GetType ...
func (l *LocalDB) GetType(queryID string) int {
	t, exists := l.queryIDToType[queryID]
	if !exists {
		panic("No type exists for queryID")
	}
	return t
}

// SetType ...
func (l *LocalDB) SetType(queryID string, fType int) {
	l.queryIDToType[queryID] = fType
}

// GetSingleSelection ...
func (l *LocalDB) GetSingleSelection(queryID string) (*policy.SingleSelection, error) {
	s, exists := l.queryIDToSingleSelection[queryID]

	var err error = nil
	if !exists {
		err = fmt.Errorf("No SingleSelection object exists for queryID: %s. It may have been cycled out of cache", queryID)
	}
	return s, err
}

// SetSingleSelection ...
func (l *LocalDB) SetSingleSelection(queryID string, selection *policy.SingleSelection) {
	l.queryIDToSingleSelection[queryID] = selection
}

// GetEnsembleSelection ...
func (l *LocalDB) GetEnsembleSelection(queryID string) (*policy.EnsembleSelection, error) {
	s, exists := l.queryIDToEnsembleSelection[queryID]

	var err error = nil
	if !exists {
		err = fmt.Errorf("No EnsembleSelection object exists for queryID: %s. It may have been cycled out of cache", queryID)
	}
	return s, err
}

// SetEnsembleSelection ...
func (l *LocalDB) SetEnsembleSelection(queryID string, selection *policy.EnsembleSelection) {
	l.queryIDToEnsembleSelection[queryID] = selection
}
