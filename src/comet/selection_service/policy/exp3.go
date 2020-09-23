package policy

import (
	"comet"
	malpb "comet/abstraction_service/pb"
	"context"
	"math"
	"math/rand"
	"sort"
)

// Exp3 is a MultiArmBandit policy
type Exp3 struct {
	SingleSelectionPolicy

	// Exploration factor in (0, 1]
	Gamma float64

	// Number of actions; in our case actions are models
	K int

	// Weights for our K actions
	Weights       []float64
	Probabilities []float64

	// Used for sampling from our distribution
	CumulativeProbabilities []float64

	AbstractionService malpb.AbstractionServiceClient
}

// CreateExp3 returns an exp3 single selection policy
func CreateExp3(gamma float64, numActions int, malClient malpb.AbstractionServiceClient) SingleSelectionPolicy {

	// initialize all weights to 1
	weights := make([]float64, numActions)
	for i := 0; i < numActions; i++ {
		weights[i] = 1
	}

	return &Exp3{
		Gamma:                   gamma,
		K:                       numActions,
		Weights:                 weights,
		Probabilities:           make([]float64, numActions),
		CumulativeProbabilities: make([]float64, numActions+1),
		AbstractionService:      malClient,
	}
}

// Select for Exp3 is written as per Panagiotopoulou
func (e *Exp3) Select(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) *SingleSelection {

	modelID := e.sampleModelID()

	// Request label
	reply, err := e.AbstractionService.Predict(ctx, &malpb.PredictRequest{
		ImageVector: comet.ImageVectorType(imageVector),
		ModelId:     int32(modelID),
		ContextUuid: contextuuid,
	})

	if err != nil {
		panic(err)
	}

	return &SingleSelection{
		Probability:     e.Probabilities[modelID],
		ModelID:         comet.ModelIDType(modelID),
		PredictionLabel: reply.Label,
	}
}

// Feedback updates the weights of the Exp3 Policy. Returns if model predicted correctly
func (e *Exp3) Feedback(singleSelection *SingleSelection, actual string) bool {
	var reward float64 = 0
	if singleSelection.PredictionLabel == actual {
		reward = 1
	}

	estimatedReward := reward / singleSelection.Probability
	e.Weights[int(singleSelection.ModelID)] *= math.Exp(estimatedReward * e.Gamma / float64(e.K))

	return reward > 0
}

func (e *Exp3) updateProbabilities() {
	explorationFactor := e.Gamma / float64(e.K)

	var weightSum float64 = 0
	for _, weight := range e.Weights {
		weightSum += weight
	}

	// update cumulative probabilities sum
	var currProbabilitySum float64 = 0
	for idx, weight := range e.Weights {
		exploitationFactor := (1 - e.Gamma) * weight / weightSum
		currProbability := explorationFactor + exploitationFactor

		// update probabilities
		e.Probabilities[idx] = currProbability

		currProbabilitySum += currProbability

		// cumulative probabilities list has K+1 values
		e.CumulativeProbabilities[idx+1] = currProbabilitySum
	}

	// enforce that the last value of cumulative probabilities is 1
	// it may not due to floating point arithmetic errors
	e.CumulativeProbabilities[len(e.CumulativeProbabilities)-1] = 1
}

func (e *Exp3) sampleModelID() comet.ModelIDType {

	// update probabilities
	e.updateProbabilities()

	// choose float between 0 and 1
	randFloat := rand.Float64()

	// search for greatest index with probability >= randFloat
	modelIDInt := sort.Search(len(e.CumulativeProbabilities), func(i int) bool {
		return e.CumulativeProbabilities[i] >= randFloat
	})
	return comet.ModelIDType(modelIDInt)
}
