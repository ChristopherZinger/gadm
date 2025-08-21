package handlers

import (
	"sync"
	"time"
)

type GlobalRateLimiter struct {
	mu          sync.RWMutex
	lastRequest time.Time
	interval    time.Duration
}

func NewGlobalRateLimiter(interval time.Duration) *GlobalRateLimiter {
	return &GlobalRateLimiter{
		interval: interval,
	}
}

func (rl *GlobalRateLimiter) IsAllowed() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if rl.lastRequest.IsZero() || now.Sub(rl.lastRequest) >= rl.interval {
		rl.lastRequest = now
		return true
	}

	return false
}

var tokenCreationRateLimiter = NewGlobalRateLimiter(2 * time.Second)
