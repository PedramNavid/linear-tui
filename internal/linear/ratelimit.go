package linear

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	refillRate int
	lastRefill time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		tokens:     1000,
		maxTokens:  1000,
		refillRate: 1000, // tokens per hour
		lastRefill: time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refill()

	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	tokensToAdd := int(elapsed.Hours() * float64(r.refillRate))

	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
		r.lastRefill = now
	}
}

func (r *RateLimiter) TokensRemaining() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refill()
	return r.tokens
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
