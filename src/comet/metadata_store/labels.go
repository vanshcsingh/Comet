package metadatastore

// This file contains the mapping of labels to vectors that all models have been trained on

import (
	"fmt"
	"strconv"
)

var (
	LABELS = [...]string{
		"0",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
	}
)

// Converts a label provided by a model container to a vector
func LabelToVector(label string) ([]float64, error) {
	vector := make([]float64, len(LABELS))

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

func VectorToLabel(vector []float64) (string, error) {
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
