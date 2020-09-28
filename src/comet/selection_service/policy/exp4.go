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
	"time"
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

	// Weights for our N advisors
	Weights []float64

	// Probabilities for our K actions
	Probabilities []float64

	// Used for sampling from our distribution
	CumulativeProbabilities []float64

	AbstractionService malpb.AbstractionServiceClient
}

// CreateExp4 returns an exp4 single selection policy
func CreateExp4(gamma float64, numLabels int, numModels int, malClient malpb.AbstractionServiceClient) EnsembleSelectionPolicy {

	// initialize all weights to 1
	weights := make([]float64, numModels)
	for i := 0; i < numModels; i++ {
		weights[i] = 1
	}

	return &Exp4{
		Gamma:                   gamma,
		K:                       numLabels,
		N:                       numModels,
		Weights:                 weights,
		Probabilities:           make([]float64, numLabels),
		CumulativeProbabilities: make([]float64, numLabels+1),
		AbstractionService:      malClient,
	}
}

// Select for Exp4 is written as per Auer
func (e *Exp4) Select(ctx context.Context, contextuuid string, imageVector comet.ImageVectorType) *EnsembleSelection {

	adviceVectors := e.getAdviceVectors(ctx, contextuuid, imageVector)
	modelID := e.sampleModelID(adviceVectors)

	AdviceForModelID := make([]float64, e.N)
	for idx, vector := range adviceVectors {
		AdviceForModelID[idx] = vector[modelID]
	}

	label, err := md.GetMetadataStoreInstance().VectorToLabel(adviceVectors[modelID])
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

// Feedback updates the weights of the Exp4 policy. Returns if model predicted correctly
func (e *Exp4) Feedback(ensembleSelection *EnsembleSelection, rewardInt int32) bool {
	reward := float64(rewardInt)

	for idx, advice := range ensembleSelection.Advice {
		estimatedReward := reward * advice / ensembleSelection.Probability
		e.Weights[idx] *= math.Exp(estimatedReward * e.Gamma / float64(e.K))
	}

	return reward > 0
}

func (e *Exp4) updateProbabilities(adviceVectors [][]float64) {
	explorationFactor := e.Gamma / float64(e.K)

	var weightSum float64 = 0
	for _, weight := range e.Weights {
		weightSum += weight
	}

	// update cumulative probabilities sum
	var currProbabilitySum float64 = 0

	// iterate through all K actions
	for idx := 0; idx < e.K; idx++ {

		var recommendationWeight float64 = 0
		for advIdx := range e.Advisors {
			recommendationWeight += adviceVectors[advIdx][idx] * e.Weights[advIdx]
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
			ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
			reply, err := e.AbstractionService.Predict(ctxWithTimeout, &malpb.PredictRequest{
				ImageVector: comet.ImageVectorType(imageVector),
				ModelId:     int32(m),
				ContextUuid: contextuuid,
			})

			if err != nil {
				panic(err)
			}

			adviceList[i], _ = md.GetMetadataStoreInstance().LabelToVector(reply.Label)
			wg.Done()
		}(idx, modelID)
	}

	wg.Wait()

	return adviceList
}
