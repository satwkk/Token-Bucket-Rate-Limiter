package cache

import (
	"math"
	"time"
)

type Bucket struct {
	Tokens       float64
	LastSeenTime time.Time
}

func (b *Bucket) RefillAndConsume() bool {
	now := time.Now()
	elapsed := now.Sub(b.LastSeenTime).Seconds()

	refillAmt := elapsed * 0.1
	b.Tokens += refillAmt

	b.Tokens = math.Min(b.Tokens, 100)

	b.LastSeenTime = time.Now()

	if b.Tokens > 1 {
		b.Tokens -= 1
		return true
	}

	return false
}
