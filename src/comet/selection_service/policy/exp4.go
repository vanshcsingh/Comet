package policy

import (
	"comet"
	malpb "comet/abstraction_service/pb"
	md "comet/metadata_store"
	"context"
	"math"
	"math/rand"
	"sort"
	"sync"
)

// Exp4 algorithm as described by Auer et al in http://rob.schapire.net/papers/AuerCeFrSc01.pdf
type Exp4 struct {
	EnsembleSelectionPolicy

	// Exploration factor in (0, 1]
	Gamma float64

	// Number of actions; for exp4 that corresponds to number of labels
	K int

	// Number of Advisors
	N int

	// Array of N ModelIds that act as our Advisors
	Advisors []comet.ModelIDType

	// Weights for our K actions
	Weights       []float64
	Probabilities []float64

	// Used for sampling from our distribution
	CumulativeProbabilities []float64

	AbstractionService malpb.AbstractionServiceClient
}

// Select for Exp4 is written as per Auer
func (e *Exp4) Select(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) *EnsembleSelection {

	adviceVectors := e.getAdviceVectors(ctx, contextuuid, imageVector)
	modelID := e.sampleModelID(adviceVectors)

	AdviceForModelID := make([]float64, e.N)
	for idx, vector := range adviceVectors {
		AdviceForModelID[idx] = vector[modelID]
	}

	label, err := md.VectorToLabel(adviceVectors[modelID])
	if err != nil {
		panic(err)
	}

	return &EnsembleSelection{
		SingleSelection: SingleSelection{
			Probability:     e.Probabilities[modelID],
			PredictionLabel: label,
			ModelID:         modelID,
		},
		Advice: AdviceForModelID,
	}
}

// Feedback updates the weights of the Exp4 policy
func (e *Exp4) Feedback(ensembleSelection *EnsembleSelection, prediction string, actual string) {
	var reward float64 = 0
	if prediction == actual {
		reward = 1
	}

	estimatedReward := reward / ensembleSelection.Probability

	for idx := range e.Advisors {
		e.Weights[idx] *= math.Exp(estimatedReward * e.Gamma / float64(e.K))
	}

	e.Weights[int(ensembleSelection.ModelID)] *= math.Exp(estimatedReward * e.Gamma / float64(e.K))
}

func (e *Exp4) updateProbabilities(adviceVectors [][]float64) {
	explorationFactor := e.Gamma / float64(e.K)

	var weightSum float64 = 0
	for _, weight := range e.Weights {
		weightSum += weight
	}

	// update cumulative probabilities sum
	var currProbabilitySum float64 = 0
	for idx := range e.Weights {

		var recommendationWeight float64 = 0
		for advIdx := range e.Advisors {
			recommendationWeight += adviceVectors[advIdx][idx]
		}

		exploitationFactor := (1 - e.Gamma) * recommendationWeight / weightSum
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

func (e *Exp4) sampleModelID(adviceVectors [][]float64) comet.ModelIDType {

	// update probabilities
	e.updateProbabilities(adviceVectors)

	// choose float between 0 and 1
	randFloat := rand.Float64()

	// search for greatest index with probability >= randFloat
	modelID := sort.Search(len(e.CumulativeProbabilities), func(i int) bool {
		return e.CumulativeProbabilities[i] >= randFloat
	})

	return comet.ModelIDType(modelID)
}

func (e *Exp4) getAdviceVectors(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) [][]float64 {
	adviceList := make([][]float64, len(e.Advisors))

	var wg sync.WaitGroup
	wg.Add(len(e.Advisors))

	for idx, modelID := range e.Advisors {
		go func(i int, m comet.ModelIDType) {
			reply, err := e.AbstractionService.Predict(ctx, &malpb.PredictRequest{
				ImageVector: comet.ImageVectorType(imageVector),
				ModelId:     int32(m),
				ContextUuid: contextuuid,
			})

			if err != nil {
				panic(err)
			}

			adviceList[i], _ = md.LabelToVector(reply.Label)
			wg.Done()
		}(idx, modelID)
	}

	wg.Wait()

	return adviceList
}
