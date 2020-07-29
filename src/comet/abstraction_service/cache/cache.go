package cache

import (
	"comet"
	"comet/abstraction_service/batch"

	"github.com/hashicorp/golang-lru"
	"sync"
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
	cache *lru.Cache

	predictProducer batch.PredictProducer
	resultConsumer batch.ResultConsumer
	
	// channel per rpc call hash
	rpcChanMap map[string] chan *comet.PredictResult
}

// CreateLocalCache implements MALCache
func CreateLocalCache (size int, predictProducer batch.PredictProducer, resultConsumer batch.ResultConsumer) MALCache {
	cache, error := lru.New(size)	
	if error != nil {
		lc := LocalCache{
			cacheWriteLock: &sync.Mutex{},
			cache: cache,
			predictProducer: predictProducer,
			resultConsumer: resultConsumer,
		}	

		go lc.pollResults()

		return &lc
	}
	return nil
}

// Start starts a thread polling the message queue of results
func (lc *LocalCache) pollResults() {
	for {
		result := lc.resultConsumer.Consume()
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
	result := <- lc.rpcChanMap[predictParams.Hash]

	// Wake up and delete channel
	delete(lc.rpcChanMap, predictParams.Hash)

	// update cache
	lc.cache.Add(predictParams.Hash, result.Label)

	return result.Label
}
