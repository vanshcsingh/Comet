package cache

import (
	"comet"
	"comet/abstraction_service/batch/mq"
	"fmt"
	"log"

	"sync"

	lru "github.com/hashicorp/golang-lru"
)

// MALCache interface abstracts the caching layer
type MALCache interface {
	Fetch(*comet.PredictParams) (string, bool)
	Request(*comet.PredictParams) string
}

// LocalCache implements MALCache
// TODO: use existing cache of fixed size
type LocalCache struct {
	cacheWriteLock *sync.Mutex
	cache          *lru.Cache

	predictProducer mq.PredictProducer
	resultConsumer  mq.ResultConsumer

	// channel per rpc call hash
	rpcChanMap map[string]chan *comet.PredictResult
}

// CreateAndStartLocalCache implements MALCache
func CreateAndStartLocalCache(size int, predictProducer mq.PredictProducer, resultConsumer mq.ResultConsumer) (MALCache, error) {
	cache, error := lru.New(size)
	if error == nil {
		lc := LocalCache{
			cacheWriteLock:  &sync.Mutex{},
			cache:           cache,
			predictProducer: predictProducer,
			resultConsumer:  resultConsumer,
			rpcChanMap:      make(map[string]chan *comet.PredictResult),
		}

		go lc.pollResults()

		return &lc, nil
	}
	return nil, fmt.Errorf("Could not create cache")
}

// Start starts a thread polling the message queue of results
func (lc *LocalCache) pollResults() {
	log.Println("[MALCache] pollResults thread has started")
	for {
		result := lc.resultConsumer.Consume()

		log.Printf("[MALCache] pollResults consumed result: %v\n", result.Label)

		lc.rpcChanMap[result.Hash] <- result
	}
}

// Fetch from cache if exists
func (lc *LocalCache) Fetch(predictParams *comet.PredictParams) (string, bool) {
	label, exists := lc.cache.Get(predictParams.Hash)
	strLabel, isString := label.(string)
	return strLabel, exists && isString
}

// Request if does not exist in cache, Predict from batching service
// Blocking call
func (lc *LocalCache) Request(predictParams *comet.PredictParams) string {
	// Fetch from local cache
	cachedLabel, exists := lc.Fetch(predictParams)
	if exists {
		return cachedLabel
	}

	// TODO: test performance agains having all Request calls wait on a shared condition variable
	// create channel
	lc.rpcChanMap[predictParams.Hash] = make(chan *comet.PredictResult)
	lc.predictProducer.Publish(predictParams)

	log.Printf("[MALCache] published request for mID: %v\n", predictParams.ModelID)

	result := <-lc.rpcChanMap[predictParams.Hash]

	// modify cache
	lc.cacheWriteLock.Lock()
	delete(lc.rpcChanMap, predictParams.Hash)      // Wake up and delete channel
	lc.cache.Add(predictParams.Hash, result.Label) // update cache
	lc.cacheWriteLock.Unlock()

	return result.Label
}
