package cache

import (
	"sync"
	"time"
)

type MemoryCache struct {
	mu    sync.Mutex
	cache map[string]Bucket
}

func NewMemoryCache() Cache {
	return &MemoryCache{
		cache: make(map[string]Bucket),
	}
}

func (c *MemoryCache) TryConsume(key string) (bool, float64, float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bucket, ok := c.cache[key]
	if !ok {
		bucket = Bucket{
			Tokens:       100,
			LastSeenTime: time.Now(),
		}
	}

	tokensBefore := bucket.Tokens
	accepted := bucket.RefillAndConsume()
	c.cache[key] = bucket

	return accepted, tokensBefore, bucket.Tokens
}
