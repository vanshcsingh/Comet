package mq

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

// Publish to local channel
func (p *LocalPredictProducer) Publish(pp *comet.PredictParams) {
	// don't block calling thread on publish
	go func() {
		p.Pipe <- pp
	}()
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

// Publish to local channel
func (p *LocalResultProducer) Publish(res *comet.PredictResult) {
	// don't block calling thread on publish
	go func() {
		p.Pipe <- res
	}()
}