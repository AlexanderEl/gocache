package gocache

import (
	"fmt"
	"sync"
	"time"
)

type GoCache interface {
	// Get the byte slice if key exists or error if it does not
	Get(key string) ([]byte, error)

	// Set a new value at given key with a max time to live
	Set(key, value string, ttl time.Duration)

	// Check if the current key exists in the cache
	Has(key string) bool

	// Remove the entry for the given key
	Delete(key string)

	// Get the number of stored elements in the cache (non-locking)
	Size() int

	// Statistic on cache operations since initialization (non-locking)
	GetStats() CacheStats
}

type CacheStats struct {
	TotalOperations uint
	NumGets         uint
	NumSets         uint
	NumHasChecks    uint
	NumDeletes      uint
}

type Cache struct {
	lock  sync.Mutex        // A lock to allow for parallelized access
	data  map[string]string // The cache data
	stats CacheStats        // Stats on the cache operations
}

func New() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.stats.NumGets++

	if value, ok := c.data[key]; ok {
		return []byte(value), nil
	} else {
		return nil, fmt.Errorf("no element in cache at given key")
	}

}

// If Set is called twice for the same key, the earliest TTL will execute
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = value
	c.stats.NumSets++

	if ttl > 0 {
		go func() {
			<-time.NewTimer(ttl).C
			c.Delete(key)
		}()
	}
}

func (c *Cache) Has(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.stats.NumHasChecks++

	_, ok := c.data[key]
	return ok
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.stats.NumDeletes++

	delete(c.data, key)
}

// Size is a non-locking method, value may not be 100% accurate if cache is actively used
func (c *Cache) Size() int {
	return len(c.data)
}

// GetStats is a non-locking method, value may not be 100% accurate if cache is actively used
func (c *Cache) GetStats() CacheStats {
	c.stats.TotalOperations = c.stats.NumGets +
		c.stats.NumSets +
		c.stats.NumHasChecks +
		c.stats.NumDeletes
	return c.stats
}
