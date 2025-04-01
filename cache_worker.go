package cache

import (
	"context"
	"log"
	"time"
)

// CacheWorkerConfig holds the configuration for starting the cache worker.
type CacheWorkerConfig struct {
	Cache    Cache           // Cache instance to clean.
	Interval time.Duration   // Interval between cache cleanup cycles.
	StopCh   <-chan struct{} // Channel used to signal the worker to stop.
}

// StartCacheWorker starts a background worker that periodically cleans expired items from the cache.
// The worker will exit when the provided context is done or when a signal is received on StopCh.
func StartCacheWorker(ctx context.Context, cfg CacheWorkerConfig) {
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	log.Println("Cache worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Cache worker: context done, stopping worker")
			return
		case <-cfg.StopCh:
			log.Println("Cache worker: stop channel signaled, stopping worker")
			return
		case <-ticker.C:
			cleanupCache(cfg.Cache)
		}
	}
}

// cleanupCache removes expired items from the cache.
// This function only works with the inMemoryCache implementation.
func cleanupCache(cache Cache) {
	memCache, ok := cache.(*inMemoryCache)
	if !ok {
		log.Println("Cache worker: cache type is not *inMemoryCache, skipping cleanup")
		return
	}

	now := time.Now()

	memCache.mu.Lock()
	defer memCache.mu.Unlock()

	for key, item := range memCache.items {
		if !item.expiration.IsZero() && now.After(item.expiration) {
			delete(memCache.items, key)
			log.Printf("Cache worker: deleted expired key: %s", key)
		}
	}
}
