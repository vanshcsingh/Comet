package batch

import (
	"comet"
	"log"
	"sync"
	"time"
)

// Service interface
type Service interface {
	// Run runs the batch Service
	Run(batchThreshold int, period time.Duration)
}

// LocalBatcher local batcher
type LocalBatcher struct {
	predictConsumer PredictConsumer
	resultProducer ResultProducer

	modelLockMap map[int32] *sync.Mutex
	modelQueue map[int32] []*comet.PredictParams

	sendBatchLock *sync.Mutex	
	sendBatchMap map[int32] bool
	sendBatch chan int32
}

// CreateLocalBatcher creates a local implementation of Service
func CreateLocalBatcher(consumer PredictConsumer, producer ResultProducer, period time.Duration) Service {
	return &LocalBatcher{
		predictConsumer: consumer,
		resultProducer: producer,

		modelLockMap: make(map[int32] *sync.Mutex),
		modelQueue: make(map[int32] []*comet.PredictParams),

		// Channel alerts run thread when to send batch of predict calls
		sendBatchLock: &sync.Mutex{},
		sendBatchMap: make(map[int32] bool),
		sendBatch: make(chan int32, 100),
	}
}

// Run runs the service
func (lb *LocalBatcher) Run(batchThreshold int, period time.Duration) {
	// Poll data from predictConsumer
	go lb.consumerThread(batchThreshold)

	ticker := time.NewTicker(period)
	for {
		select {
		case mID := <- lb.sendBatch:
			go lb.batchPredictCalls(mID, lb.extractPredictParams(mID))
			
		case <- ticker.C:
			for mID := range lb.modelLockMap {
				go lb.batchPredictCalls(mID, lb.extractPredictParams(mID))
			}
		}
	}
}

func (lb *LocalBatcher) consumerThread(batchThreshold int) {
	// forever loop
	for {
		// poll from consumer
		predictParams := lb.predictConsumer.Consume()
		mID := predictParams.ModelID 

		lb.lockOnModelID(mID)
		lb.modelQueue[mID] = append(lb.modelQueue[mID], predictParams)
		lb.unlockOnModelID(mID)

		if len(lb.modelQueue[mID]) > batchThreshold {
			lb.sendBatchLock.Lock()
			if ! lb.sendBatchMap[mID] {
				lb.sendBatchMap[mID] = true
				lb.sendBatch <- mID
			}
			lb.sendBatchLock.Unlock()
		}
	}
}

func (lb *LocalBatcher) extractPredictParams(mID int32) []*comet.PredictParams {
	lb.lockOnModelID(mID)
	copyParams := lb.modelQueue[mID]
	
	// clear modelQueue for mID
	lb.modelQueue[mID] = nil
	lb.unlockOnModelID(mID)

	return copyParams
}

// This function is currently stubbed to test batchPredictCalls
// It simply publishes a result to a label
func (lb *LocalBatcher) batchPredictCalls(mID int32, params []*comet.PredictParams) {
	log.Println("batching", len(params), "predict calls on model#", mID)
	log.Println("Predict params are: ", params)
	log.Println("---------------------------------------------")
	log.Println()


	labels := []string{"hot dog", "not hot dog", "cat"}
	label := labels[int(mID) % len(labels)]

	for _, p := range params {
		lb.resultProducer.Publish(&comet.PredictResult{
			Label: label,
			Hash: p.Hash,
		})
	}
}

// lockOnModelID locks a model's queue on its modelID
func (lb *LocalBatcher) lockOnModelID(mID int32) {
	_, exists := lb.modelLockMap[mID]
	if !exists {
		lb.modelLockMap[mID] = &sync.Mutex{}
	}
	lb.modelLockMap[mID].Lock()
}

// unlockOnModelID locks a model's queue on its modelID
func (lb *LocalBatcher) unlockOnModelID(mID int32) {
	lb.modelLockMap[mID].Unlock()
}