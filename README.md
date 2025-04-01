# GoStash

GoStash is a lightweight, in-memory cache package written in Go. It provides a simple API for caching data with support for time-to-live (TTL) expirations, safe concurrent access, and an optional background worker that cleans up expired items automatically.

## Features

- **Simple API**: Easily set, retrieve, and delete cached items.
- **TTL Support**: Optionally set a TTL for each cache entry.
- **Concurrency Safe**: Built-in thread safety using sync.RWMutex.
- **Cache Worker**: Background worker for automatic cleanup of expired items.
- **Modular Design**: Clean and well-organized code, making it easy to integrate into any project.

## Installation

Install the package using go get:

```bash
go get github.com/yourusername/gocache
```

Note: Replace `github.com/yourusername/gocache` with the actual module path of the repository.

## Usage

Below is an example demonstrating how to use GoCache, including setting cache entries, retrieving them, and running the background cleanup worker.

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/yourusername/gocache/cache"
)

func main() {
	// Create a new in-memory cache instance.
	c := cache.NewCache()

	// Set an example key with a TTL of 5 seconds.
	c.SetWithTTL("example", "Hello, Cache!", 5*time.Second)

	// Create a context that can be canceled for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a stop channel to signal the cache worker to stop.
	stopCh := make(chan struct{})

	// Start the cache cleanup worker in a separate goroutine.
	go cache.StartCacheWorker(ctx, cache.CacheWorkerConfig{
		Cache:    c,
		Interval: 2 * time.Second, // Cleanup interval: every 2 seconds.
		StopCh:   stopCh,
	})

	// Set up a channel to listen for OS interrupt signals (e.g., Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	log.Println("Application started. Press Ctrl+C to exit.")

	// Periodically check for the existence of the "example" key.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if val, ok := c.Get("example"); ok {
				log.Printf("Key 'example' found with value: %v", val)
			} else {
				log.Println("Key 'example' not found or expired.")
			}
		case sig := <-sigCh:
			// Gracefully shutdown on interrupt signal.
			log.Printf("Received signal: %v. Shutting down.", sig)
			close(stopCh) // Signal the cache worker to stop.
			cancel()      // Cancel the context.
			// Allow some time for the worker to finish.
			time.Sleep(1 * time.Second)
			return
		}
	}
}
```

## API Reference

### Cache Interface

The primary interface for interacting with the cache:

```go
type Cache interface {
    // Set assigns a value to the specified key without expiration.
    Set(key string, value any)
    // SetWithTTL assigns a value to the specified key with a TTL.
    // If ttl <= 0, the item will not expire.
    SetWithTTL(key string, value any, ttl time.Duration)
    // Get retrieves the value associated with the specified key.
    // Returns (nil, false) if the key does not exist or if the item is expired.
    Get(key string) (any, bool)
    // Delete removes the specified key from the cache.
    Delete(key string)
    // Clear removes all items from the cache.
    Clear()
}
```

### Cache Worker

The cache worker automatically cleans up expired items. Configure it using `CacheWorkerConfig` and start it with `StartCacheWorker`.

```go
type CacheWorkerConfig struct {
    Cache    Cache         // Cache instance to clean.
    Interval time.Duration // Interval between cleanup cycles.
    StopCh   <-chan struct{} // Channel to signal the worker to stop.
}
```

Start the worker with:

```go
func StartCacheWorker(ctx context.Context, cfg CacheWorkerConfig)
```

## Contributing

Contributions are welcome! If you have ideas, bug fixes, or enhancements, please fork the repository and open a pull request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License.

## Support

If you encounter any issues or have any questions, feel free to open an issue on the repository or contact the maintainer.

Happy caching with GoCache!
