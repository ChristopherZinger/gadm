package accessTokenCache

import (
	"errors"
	"sync"
	"time"
)

const RATE_LIMIT_DURATION = 1 * time.Second
const NUM_HITS_PER_RATE_LIMIT = 10

var TOKEN_LIVE_DURATION = struct {
	years  int
	months int
	days   int
}{
	years:  0,
	months: 3,
	days:   0,
}

type TokenRateInfo struct {
	hitHistory []time.Time
	createdAt  time.Time
	mu         sync.RWMutex
}

func newTokenRateInfo(createdAt time.Time) *TokenRateInfo {
	return &TokenRateInfo{
		hitHistory: []time.Time{},
		createdAt:  createdAt,
		mu:         sync.RWMutex{},
	}
}

func (tri *TokenRateInfo) handleHit() error {
	tri.mu.Lock()
	defer tri.mu.Unlock()

	now := time.Now()
	if isTokenExpired(tri.createdAt) {
		return errors.New(TokenExpiredMsg)
	}

	i := 0
	for _, hitTime := range tri.hitHistory {
		if now.Sub(hitTime) > RATE_LIMIT_DURATION {
			i++
			continue
		}
		break
	}
	tri.hitHistory = tri.hitHistory[i:]

	if len(tri.hitHistory) >= NUM_HITS_PER_RATE_LIMIT {
		return errors.New(RateLimitExceededMsg)
	}

	tri.hitHistory = append(tri.hitHistory, now)
	return nil
}

func isTokenExpired(tokenCreationTime time.Time) bool {
	tokenExpirationTime := tokenCreationTime.AddDate(
		TOKEN_LIVE_DURATION.years,
		TOKEN_LIVE_DURATION.months,
		TOKEN_LIVE_DURATION.days)

	return tokenExpirationTime.Before(time.Now())
}
