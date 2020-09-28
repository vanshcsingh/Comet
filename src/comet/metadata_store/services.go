package metadata_store

// This file contains the mapping of labels to vectors that all models have been trained on

import (
	"fmt"
	"strconv"
)

// Converts a label provided by a model container to a vector
func (md *LocalFileBasedMetadataStore) LabelToVector(label string) ([]float64, error) {
	vector := make([]float64, len(md.selectionParams.Labels))

	labelIdx, err := strconv.Atoi(label)
	if err != nil {
		return nil, err
	}

	if labelIdx < 0 || labelIdx > 9 {
		return nil, fmt.Errorf("label must be between 0 and 9")
	}

	vector[labelIdx] = 1.0
	return vector, nil

}

func (md *LocalFileBasedMetadataStore) VectorToLabel(vector []float64) (string, error) {
	oneIndex := 0
	count := 0

	for idx, val := range vector {
		if val > 0 {
			oneIndex = idx
			count++
		}
	}

	if count != 1 {
		return "", fmt.Errorf("Illformed vector; must have exactly one nonzero element")
	}

	strIndex := strconv.Itoa(oneIndex)
	return strIndex, nil
}

func (md *LocalFileBasedMetadataStore) GetNumModels() int {
	return len(md.entryMap)
}

func (md *LocalFileBasedMetadataStore) GetNumLabels() int {
	return len(md.selectionParams.Labels)
}

func (md *LocalFileBasedMetadataStore) GetSelectionGamma() float64 {
	return float64(md.selectionParams.Gamma)
}

func (md *LocalFileBasedMetadataStore) GetAbstractionServiceAddr() string {
	return md.abstractionServiceAddr
}

func (md *LocalFileBasedMetadataStore) GetSelectionServiceAddr() string {
	return md.selectionServiceAddr
}
