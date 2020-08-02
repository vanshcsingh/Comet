package mq

import (
	"comet"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createPredictQueueElements() (PredictConsumer, PredictProducer) {
	predictPipe := make(chan *comet.PredictParams)
	predictConsumer := &LocalPredictConsumer{Pipe: predictPipe}
	predictProducer := &LocalPredictProducer{Pipe: predictPipe}

	return predictConsumer, predictProducer
}

func createResultQueueElements() (ResultConsumer, ResultProducer) {
	resultPipe := make(chan *comet.PredictResult)
	resultConsumer := &LocalResultConsumer{Pipe: resultPipe}
	resultProducer := &LocalResultProducer{Pipe: resultPipe}

	return resultConsumer, resultProducer
}

func TestPredictQueue(t *testing.T) {

	vector1 := comet.ImageVectorType([]int32{1})
	vector2 := comet.ImageVectorType([]int32{2,3})
	vector3 := comet.ImageVectorType([]int32{3,4,5})

	params1 := &comet.PredictParams{
		ImageVector: vector1,
		ModelID: comet.ModelIDType(0),
		ContextUUID: "uuid1",
		Hash: "hash1",
	}

	params2 := &comet.PredictParams{
		ImageVector: vector2,
		ModelID: comet.ModelIDType(1),
		ContextUUID: "uuid2",
		Hash: "hash2",
	}

	params3 := &comet.PredictParams{
		ImageVector: vector3,
		ModelID: comet.ModelIDType(2),
		ContextUUID: "uuid3",
		Hash: "hash3",
	}

	predictConsumer, predictProducer := createPredictQueueElements()

	for _, p := range []*comet.PredictParams{params1, params2, params3} {
		predictProducer.Publish(p)
		assert.Equal(t, p, predictConsumer.Consume())
	}
}

func TestResultQueue(t *testing.T) {

	result1 := &comet.PredictResult{
		Label: "label1",
		Hash: "hash1",
	}

	result2 := &comet.PredictResult{
		Label: "label2",
		Hash: "hash2",
	}

	result3 := &comet.PredictResult{
		Label: "label3",
		Hash: "hash3",
	}

	resultConsumer, resultProducer := createResultQueueElements()

	for _, r := range []*comet.PredictResult{result1, result2, result3} {
		resultProducer.Publish(r)
		assert.Equal(t, r, resultConsumer.Consume())
	}
}