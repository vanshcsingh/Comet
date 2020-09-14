package batch

import (
	"comet"
	"comet/abstraction_service/batch/mq"
	md "comet/metadata_store"
	"context"

	modelpb "comet/models/container_models"

	"fmt"
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
	resultProducer  mq.ResultProducer

	modelLockMap     map[comet.ModelIDType]*sync.Mutex
	modelLockMapLock *sync.Mutex
	modelQueue       map[comet.ModelIDType][]*comet.PredictParams

	// sendBatch channel alerts us of which model type is ready to send batched RPCs to models
	sendBatch     chan comet.ModelIDType
	sendBatchLock *sync.Mutex
	sendBatchMap  map[comet.ModelIDType]bool

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
		resultProducer:  producer,

		modelLockMap:     make(map[comet.ModelIDType]*sync.Mutex),
		modelLockMapLock: &sync.Mutex{},
		modelQueue:       make(map[comet.ModelIDType][]*comet.PredictParams),

		// Channel alerts run thread when to send batch of predict calls
		sendBatchLock: &sync.Mutex{},
		sendBatchMap:  make(map[comet.ModelIDType]bool),
		sendBatch:     make(chan comet.ModelIDType, 100),

		mdStore: mdStore,
	}
	go lb.Run(batchThreshold, duration)
	return lb
}

// Run runs the service
func (lb *LocalBatcher) Run(batchThreshold int, period time.Duration) {
	// Poll data from predictConsumer
	go lb.consumerThread(batchThreshold)

	log.Println("[Batcher]: Batcher is starting up")

	ticker := time.NewTicker(period)
	for {
		select {
		case mID := <-lb.sendBatch:
			go lb.batchPredictCalls(mID, lb.extractPredictParams(mID))

		case <-ticker.C:
			for mID := range lb.modelLockMap {
				go lb.batchPredictCalls(mID, lb.extractPredictParams(mID))
			}
		}
	}
}

func (lb *LocalBatcher) consumerThread(batchThreshold int) {
	log.Println("[Batcher]: consumerThread initialized")

	// forever loop
	for {
		// poll from consumer
		predictParams := lb.predictConsumer.Consume()
		mID := predictParams.ModelID

		log.Println("[Batcher: consumerThread] inbound params for mID: ", mID)

		lb.lockOnModelID(mID)
		lb.modelQueue[mID] = append(lb.modelQueue[mID], predictParams)
		lb.unlockOnModelID(mID)

		log.Println("[Batcher: consumerThread]: model queue has length", len(lb.modelQueue[mID]))

		if len(lb.modelQueue[mID]) >= batchThreshold {
			lb.sendBatchLock.Lock()
			if !lb.sendBatchMap[mID] {
				log.Println("[Batcher: consumerThread] sending batch signal for mID", mID)
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

// It performs batched predict calls and publishes labels
func (lb *LocalBatcher) batchPredictCalls(mID comet.ModelIDType, params []*comet.PredictParams) {

	log.Printf("[Batcher: batchPredictCalls] for mid: %d\n", mID)

	client, err := lb.mdStore.GetClient(mID)
	if err != nil {
		panic("Could not get model client")
	}

	log.Printf("[Batcher: batchPredictionCalls] Received a client: ", client)

	modelImageVectors := make([]*modelpb.ImageVector, 0)

	for _, p := range params {
		modelImageVectors = append(modelImageVectors, &modelpb.ImageVector{
			Pixels: p.ImageVector,
		})
	}

	modelPredictRequest := &modelpb.PredictRequest{Images: modelImageVectors}

	log.Printf("[Batcher: batchPredictionCalls] about to make a predict call")

	reply, err := client.Predict(
		context.Background(),
		modelPredictRequest,
	)

	log.Printf("[Batcher: batchPredictionCalls] Got a reply back from the client")

	if err != nil {
		panic(fmt.Sprintf("Received error from making RPC to container: %v", err))
	}

	// results are ordered in the same way that they were provided
	for idx, p := range params {
		lb.resultProducer.Publish(&comet.PredictResult{
			Label: reply.Labels[idx],
			Hash:  p.Hash,
		})
	}

	lb.sendBatchLock.Lock()
	lb.sendBatchMap[mID] = false
	lb.sendBatchLock.Unlock()

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
