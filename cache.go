package cache

import (
	"time"
)

// Cache defines the interface for the cache.
type Cache interface {
	// Set assigns a value to the specified key without expiration.
	Set(key string, value any)
	// SetWithTTL assigns a value to the specified key with a given time-to-live (TTL).
	// If ttl <= 0, the item does not expire.
	SetWithTTL(key string, value any, ttl time.Duration)
	// Get retrieves the value for the specified key.
	// Returns (nil, false) if the key does not exist or if the item is expired.
	Get(key string) (any, bool)
	// Delete removes the item associated with the specified key.
	Delete(key string)
	// Clear removes all items from the cache.
	Clear()
}
