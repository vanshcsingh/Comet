package comet

import (
	malpb "comet/abstraction_service/pb"
	"crypto/sha256"
	"fmt"
)

// ModelIDType wraps the type of ModelIDs
type ModelIDType int32

// ImageVectorType defines an MNIST image type
type ImageVectorType []int32

// PredictParams extracts and contains the important components of
// the PredictRequest abstraction_service pb
type PredictParams struct {
	ImageVector ImageVectorType
	ModelID     ModelIDType
	ContextUUID string
	Hash        string
}

// PredictResult is returned by the batch message queue
type PredictResult struct {
	Label string
	Hash  string
}

// CreatePredictParamsMAL object from a Predict RPC
func CreatePredictParamsMAL(pr *malpb.PredictRequest) *PredictParams {
	return &PredictParams{
		ImageVector: pr.GetImageVector(),
		ModelID:     ModelIDType(pr.GetModelId()),
		ContextUUID: pr.GetContextUuid(),
		Hash:        predictRequestHash(pr),
	}
}

// In the future we can devise hashes that are more injective
func predictRequestHash(pr *malpb.PredictRequest) string {
	hash := sha256.New()
	fmt.Fprint(hash, pr)
	return string(hash.Sum(nil))
}
