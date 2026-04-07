// Package pokecache provides an in-memory cache with automatic expiration.
package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is a thread-safe in-memory key-value store with time-based expiration.
type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
	done     chan struct{}
}

// NewCache creates a cache that automatically reaps entries older than interval.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		done:     make(chan struct{}),
	}
	go c.reapLoop()
	return c
}

// Add stores val under key with the current timestamp.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves the value for key. The bool indicates whether the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// Stop terminates the background reap goroutine.
func (c *Cache) Stop() {
	close(c.done)
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case now := <-ticker.C:
			c.mu.Lock()
			for k, v := range c.entries {
				if now.Sub(v.createdAt) > c.interval {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
		}
	}
}
