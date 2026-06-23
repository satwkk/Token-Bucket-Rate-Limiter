package cache

import (
	"math"
	"time"
)

type Bucket struct {
	CurrentToken float64
	LastSeenTime time.Time
}

func (b *Bucket) RefillAndConsume(cfg *CacheConfig) bool {
	now := time.Now()
	elapsed := now.Sub(b.LastSeenTime).Seconds()

	refillAmt := elapsed * cfg.RefillRate
	b.CurrentToken += refillAmt

	b.CurrentToken = math.Min(b.CurrentToken, cfg.MaxToken)

	b.LastSeenTime = time.Now()

	if b.CurrentToken > 1 {
		b.CurrentToken -= 1.0
		return true
	}

	return false
}
