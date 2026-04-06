package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

// Create a new cache instance with an interval for reaping old entries.
func NewCache(interval time.Duration) *Cache {
	new := &Cache{
		cache: make(map[string]cacheEntry),
	}

	// Start the reaping loop in a separate goroutine.
	ch := make(chan time.Duration)
	go new.reapLoop(ch)
	ch <- interval

	return new
}

// Add a new entry to the cache with the current timestamp.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

// Get an entry from the cache. Returns the value and a boolean indicating if the key exists.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

// reapLoop will periodically delete cache entries older than a certain threshold (e.g., 5 seconds).
func (c *Cache) reapLoop(ch chan time.Duration) {
	interval := <-ch
	// Create a ticker that ticks at the specified interval.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// When ticker ticks, check for old entries and delete them
	for range ticker.C { // will block until the ticker ticks
		c.mu.Lock()
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) >= interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
