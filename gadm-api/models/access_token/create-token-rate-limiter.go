package access_token

import (
	"sync"
	"time"
)

type AccessTokenCreationRateLimiter struct {
	mu          sync.RWMutex
	lastRequest time.Time
	interval    time.Duration
}

var tokenCreationRateLimiter *AccessTokenCreationRateLimiter

func NewAccessTokenCreationRateLimiter() *AccessTokenCreationRateLimiter {
	interval := 2 * time.Second
	if tokenCreationRateLimiter == nil {
		tokenCreationRateLimiter = &AccessTokenCreationRateLimiter{
			interval: interval,
		}
		return tokenCreationRateLimiter
	}
	return tokenCreationRateLimiter
}

func (rl *AccessTokenCreationRateLimiter) IsAllowed() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if rl.lastRequest.IsZero() || now.Sub(rl.lastRequest) >= rl.interval {
		rl.lastRequest = now
		return true
	}

	return false
}
