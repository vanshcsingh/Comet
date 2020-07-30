package batch

import (
	"comet"
)

// PredictProducer allows us to develop around this interface
// so that we can start on a local channel-based implementation and
// then scale to a Kafka implementation
// Produce the parameters for a predict call 
type PredictProducer interface {
	Publish(*comet.PredictParams)
}

// LocalPredictProducer is a local producer for our single-node channel based 
// messagequeue system
type LocalPredictProducer struct {
	// write only channel
	Pipe chan<- *comet.PredictParams
}

// Publish publishes to local channel
func (p *LocalPredictProducer) Publish(pp *comet.PredictParams) {
	p.Pipe <- pp
}

// ResultProducer allows us to develop around this interface
// so that we can start on a local channel-based implementation and
// then scale to a Kafka implementation
// Produce the result of prediction
type ResultProducer interface {
	Publish(*comet.PredictResult)
}

// LocalResultProducer is a local producer for our single-node channel based 
// messagequeue system
type LocalResultProducer struct {
	// write only channel
	Pipe chan<- *comet.PredictResult
}

// Publish publishes to local channel
func (p *LocalResultProducer) Publish(res *comet.PredictResult) {
	p.Pipe <- res
}