package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	cache "github.com/nordew/go-stash"
)

func main() {
	// Create a new in-memory cache instance.
	c := cache.NewCache()

	// Set an example key with a TTL of 5 seconds.
	c.SetWithTTL("example", "Hello, Cache!", 5*time.Second)

	// Create a context that can be cancelled for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a stop channel to signal the cache worker to stop.
	stopCh := make(chan struct{})

	// Start the cache cleanup worker in a separate goroutine.
	go cache.StartCacheWorker(ctx, cache.CacheWorkerConfig{
		Cache:    c,
		Interval: 2 * time.Second, // Cleanup every 2 seconds.
		StopCh:   stopCh,
	})

	// Set up a channel to listen for OS interrupt signals (e.g., Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	log.Println("Application started. Press Ctrl+C to exit.")

	// Use a ticker to periodically check the value of the key.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Retrieve the "example" key from the cache.
			if val, ok := c.Get("example"); ok {
				log.Printf("Key 'example' found with value: %v", val)
			} else {
				log.Println("Key 'example' not found or expired.")
			}
		case sig := <-sigCh:
			// Handle an OS interrupt signal for graceful shutdown.
			log.Printf("Received signal: %v. Shutting down.", sig)
			close(stopCh) // Signal the cache worker to stop.
			cancel()      // Cancel the context.
			// Allow some time for the worker to finish cleanup.
			time.Sleep(1 * time.Second)
			return
		}
	}
}
