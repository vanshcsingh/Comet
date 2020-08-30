package cache

import (
	"comet"
	"testing"

	"comet/abstraction_service/batch/mq"

	"github.com/stretchr/testify/assert"
)

func setupBatcher(t *testing.T) (MALCache, mq.PredictConsumer, mq.ResultProducer) {

	predictPipe := make(chan *comet.PredictParams)
	predictConsumer := &mq.LocalPredictConsumer{Pipe: predictPipe}
	predictProducer := &mq.LocalPredictProducer{Pipe: predictPipe}

	resultPipe := make(chan *comet.PredictResult)
	resultConsumer := &mq.LocalResultConsumer{Pipe: resultPipe}
	resultProducer := &mq.LocalResultProducer{Pipe: resultPipe}

	cacheSize := 2

	// create and start caching service
	cache, cacheCreationError := CreateAndStartLocalCache(cacheSize, predictProducer, resultConsumer)
	assert.Nil(t, cacheCreationError)

	return cache, predictConsumer, resultProducer
}

// Struct to store local test's fields. Do not export
type testParams struct {
	Hash  string
	Param *comet.PredictParams
	Label string
}

// TestCache tests the Fetch and Request methods of our local cache
func TestCache(t *testing.T) {
	cache, predictConsumer, resultProducer := setupBatcher(t)

	hash0 := "hash0"
	tp0 := testParams{
		Hash: hash0,
		Param: &comet.PredictParams{
			ImageVector: comet.ImageVectorType{},
			ModelID:     0,
			ContextUUID: "",
			Hash:        hash0,
		},
		Label: "hot dog",
	}

	hash1 := "hash1"
	tp1 := testParams{
		Hash: hash1,
		Param: &comet.PredictParams{
			ImageVector: comet.ImageVectorType{},
			ModelID:     0,
			ContextUUID: "",
			Hash:        hash1,
		},
		Label: "not hot dog",
	}

	for _, tp := range []testParams{tp0, tp1} {

		_, res := cache.Fetch(tp.Param)
		assert.False(t, res)

		// Test Request is successful
		requestCallDone := make(chan bool)

		go func(done chan bool) {
			assert.Equal(t, cache.Request(tp.Param), tp.Label)
			requestCallDone <- true
		}(requestCallDone)

		assert.Equal(t, predictConsumer.Consume(), tp.Param)
		resultProducer.Publish(&comet.PredictResult{
			Label: tp.Label,
			Hash:  tp.Hash,
		})

		// Wait for Request call to finish so that cache is populated
		<-requestCallDone
		prediction, _ := cache.Fetch(tp.Param)
		assert.Equal(t, prediction, tp.Label)
	}
}
