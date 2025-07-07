package linear

import (
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	refillRate int
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
// Linear API allows 1000 requests per hour
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		tokens:     1000,
		maxTokens:  1000,
		refillRate: 1000, // tokens per hour
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed and consumes a token if so
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Refill tokens based on time elapsed
	r.refill()

	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

// refill adds tokens based on time elapsed
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// Calculate tokens to add (rate per hour converted to tokens per elapsed time)
	tokensToAdd := int(elapsed.Hours() * float64(r.refillRate))

	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
		r.lastRefill = now
	}
}

// TokensRemaining returns the number of tokens remaining
func (r *RateLimiter) TokensRemaining() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.refill()
	return r.tokens
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
