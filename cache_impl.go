package cache

import (
	"sync"
	"time"
)

// cachedItem represents an item stored in the cache.
type cachedItem struct {
	value      any
	expiration time.Time
}

// isExpired checks whether the cached item has expired.
func (ci cachedItem) isExpired() bool {
	if ci.expiration.IsZero() {
		return false
	}

	return time.Now().After(ci.expiration)
}

// inMemoryCache is an in-memory cache implementation.
type inMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cachedItem
}

// NewCache creates and returns a new instance of inMemoryCache that implements the Cache interface.
func NewCache() Cache {
	return &inMemoryCache{
		items: make(map[string]cachedItem),
	}
}

// Get retrieves the value for the specified key if it exists and is not expired.
// If the item is expired, it is removed and (nil, false) is returned.
func (c *inMemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if item.isExpired() {
		c.Delete(key)
		return nil, false
	}

	return item.value, true
}

// Set assigns a value to the specified key without setting an expiration.
func (c *inMemoryCache) Set(key string, value any) {
	c.SetWithTTL(key, value, 0)
}

// SetWithTTL assigns a value to the specified key with a TTL.
// If ttl <= 0, the item does not expire.
func (c *inMemoryCache) SetWithTTL(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}
	c.items[key] = cachedItem{
		value:      value,
		expiration: expiration,
	}
}

// Delete removes the item associated with the specified key from the cache.
func (c *inMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *inMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]cachedItem)
}
