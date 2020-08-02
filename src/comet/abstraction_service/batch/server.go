package batch

import (
	"comet"
	md "comet/metadata_store"
	"comet/abstraction_service/batch/mq"

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
	predictConsumer mq.PredictConsumer
	resultProducer mq.ResultProducer

	modelLockMap map[comet.ModelIDType] *sync.Mutex
	modelLockMapLock *sync.Mutex
	modelQueue map[comet.ModelIDType] []*comet.PredictParams

	sendBatchLock *sync.Mutex	
	sendBatchMap map[comet.ModelIDType] bool
	sendBatch chan comet.ModelIDType

	mdStore md.MetadataStore
}

// CreateAndStartLocalBatcher creates a local implementation of Service
func CreateAndStartLocalBatcher(
	consumer mq.PredictConsumer, 
	producer mq.ResultProducer, 
	batchThreshold int, 
	duration time.Duration,
	mdStore md.MetadataStore,
) Service {
	lb := &LocalBatcher{
		predictConsumer: consumer,
		resultProducer: producer,

		modelLockMap: make(map[comet.ModelIDType] *sync.Mutex),
		modelLockMapLock: &sync.Mutex{},
		modelQueue: make(map[comet.ModelIDType] []*comet.PredictParams),

		// Channel alerts run thread when to send batch of predict calls
		sendBatchLock: &sync.Mutex{},
		sendBatchMap: make(map[comet.ModelIDType] bool),
		sendBatch: make(chan comet.ModelIDType, 100),

		mdStore: mdStore,
	}
	go lb.Run(batchThreshold, duration)
	return lb
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

func (lb *LocalBatcher) extractPredictParams(mID comet.ModelIDType) []*comet.PredictParams {
	lb.lockOnModelID(mID)
	copyParams := lb.modelQueue[mID]
	
	// clear modelQueue for mID
	lb.modelQueue[mID] = nil
	lb.unlockOnModelID(mID)

	return copyParams
}

// This function is currently stubbed to test batchPredictCalls
// It simply publishes a result to a label
func (lb *LocalBatcher) batchPredictCalls(mID comet.ModelIDType, params []*comet.PredictParams) {
	log.Println("batching", len(params), "predict calls on model#", mID, "to address: ", lb.mdStore.GetEntry(mID).Addr)
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
func (lb *LocalBatcher) lockOnModelID(mID comet.ModelIDType) {
	_, exists := lb.modelLockMap[mID]
	if !exists {
		// Lock inside of if statement so that codepath of 99.99% of calls
		// are not impeded by synchronization congestion
		lb.modelLockMapLock.Lock()
		_, stillExists := lb.modelLockMap[mID]
		if !stillExists {
			lb.modelLockMap[mID] = &sync.Mutex{}
		}
		lb.modelLockMapLock.Unlock()
	}
	lb.modelLockMap[mID].Lock()
}

// unlockOnModelID locks a model's queue on its modelID
func (lb *LocalBatcher) unlockOnModelID(mID comet.ModelIDType) {
	lb.modelLockMap[mID].Unlock()
}
