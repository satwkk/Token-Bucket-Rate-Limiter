package cache

import (
	"sync"
	"time"
)

type MemoryCache struct {
	mu     sync.Mutex
	cache  map[string]Bucket
	config CacheConfig
}

type CacheConfig struct {
	MaxToken   float64
	RefillRate float64
}

func NewMemoryCache(config CacheConfig) Cache {
	return &MemoryCache{
		cache:  make(map[string]Bucket),
		config: config,
	}
}

func (c *MemoryCache) TryConsume(key string) (bool, float64, float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bucket, ok := c.cache[key]
	if !ok {
		bucket = Bucket{
			CurrentToken: c.config.MaxToken,
			LastSeenTime: time.Now(),
		}
	}

	tokensBefore := bucket.CurrentToken
	accepted := bucket.RefillAndConsume(&c.config)
	c.cache[key] = bucket

	return accepted, tokensBefore, bucket.CurrentToken
}
