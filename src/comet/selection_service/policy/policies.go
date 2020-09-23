package policy

import (
	"comet"
	"context"
)

// SingleSelection records the state of our Policy when we made a single selection
type SingleSelection struct {
	ModelID         comet.ModelIDType
	PredictionLabel string
	ActionID        int
	Probability     float64
}

// EnsembleSelection encapsulates the fields of SingleSelection and records the advice
// given by its experts
type EnsembleSelection struct {
	SingleSelection

	// Advice[idx] corresponds to the advice given by the idx Advisor on the selection
	Advice        []float64
	Probabilities []float64
}

// SingleSelectionPolicy is for policies that return the result of a single model
type SingleSelectionPolicy interface {
	Select(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) *SingleSelection
	Feedback(singleSelection *SingleSelection, actual string) bool
}

// EnsembleSelectionPolicy is for policies that ensemble the results of multiple models
type EnsembleSelectionPolicy interface {
	Select(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) *EnsembleSelection
	Feedback(ensembleSelection *EnsembleSelection, actual string) bool
}
