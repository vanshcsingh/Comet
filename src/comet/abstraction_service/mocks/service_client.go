package mocks

import (
	malpb "comet/abstraction_service/pb"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// MockAbstractionServiceClientStruct implements the AbstractionServiceClient interface
type MockAbstractionServiceClientStruct struct {
	modelIDToLabelMap map[int32]string
}

// CreateAbstractionServiceMock returns an instance of the MockAbstractionServiceClientStruct
func CreateAbstractionServiceMock(modelIDToLabelMap map[int32]string) malpb.AbstractionServiceClient {
	return &MockAbstractionServiceClientStruct{
		modelIDToLabelMap: modelIDToLabelMap,
	}
}

// Predict returns the label set for a particular modelID
func (m *MockAbstractionServiceClientStruct) Predict(ctx context.Context, in *malpb.PredictRequest, opts ...grpc.CallOption) (*malpb.PredictReply, error) {
	modelID := in.GetModelId()

	if label, exists := m.modelIDToLabelMap[modelID]; exists {
		return &malpb.PredictReply{
			Label: label,
		}, nil
	}

	return nil, fmt.Errorf("Mock does not have label for modelID %d", modelID)
}
