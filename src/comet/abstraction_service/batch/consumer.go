package batch

import (
	"comet"
)

// PredictConsumer allows us to develop around this interface
// so that we can start on a local channel-based implementation and
// then scale to a Kafka implementation
// Consume the parameters for a predict call
type PredictConsumer interface {
	Consume() *comet.PredictParams
}

// LocalPredictConsumer is a local consumer for our single-node channel based 
// messagequeue system
type LocalPredictConsumer struct {
	// read only channel
	pipe <-chan *comet.PredictParams
}

// Consume consumes from shared local channel
func (c *LocalPredictConsumer) Consume() *comet.PredictParams {
	return <- c.pipe
}

// ResultConsumer allows us to develop around this interface
// so that we can start on a local channel-based implementation and
// then scale to a Kafka implementation
// Consume the result of a predict call
type ResultConsumer interface {
	Consume() *comet.PredictResult
}

// LocalResultConsumer is a local consumer for our single-node channel based 
// messagequeue system
type LocalResultConsumer struct {
	// read only channel
	pipe <-chan *comet.PredictResult
}

// Consume consumes from shared local channel
func (c *LocalResultConsumer) Consume() *comet.PredictResult {
	return <- c.pipe
}