package accessTokenCache

import (
	"errors"
	"sync"
	"time"
)

type TokenCache struct {
	tokenToTokenRateInfo map[string]*TokenRateInfo
	mu                   sync.RWMutex
}

func NewTokenCache() *TokenCache {
	return &TokenCache{
		tokenToTokenRateInfo: make(map[string]*TokenRateInfo),
	}
}

func (cache *TokenCache) SetIfNotExpired(token string, tokenRateInfo *TokenRateInfo) error {
	if isExpired(tokenRateInfo.createdAt) {
		return errors.New(TokenExpiredMsg)
	}

	cache.tokenToTokenRateInfo[token] = tokenRateInfo
	return nil
}

func (cache *TokenCache) HandleHitForToken(token string, getTokenCreatedAtIfNotInCache func(token string) (time.Time, error)) error {
	cache.mu.RLock()
	if tokenRateInfo, exists := cache.tokenToTokenRateInfo[token]; exists {
		defer cache.mu.RUnlock()
		return tokenRateInfo.handleHit()
	}
	cache.mu.RUnlock()

	createdAt, err := getTokenCreatedAtIfNotInCache(token)
	if err != nil {
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()
	if tokenRateInfo, exists := cache.tokenToTokenRateInfo[token]; exists {
		return tokenRateInfo.handleHit()
	}

	tokenRateInfo := newTokenRateInfo(createdAt)
	if err = cache.SetIfNotExpired(token, tokenRateInfo); err != nil {
		return err
	}

	return tokenRateInfo.handleHit()
}

var TOKEN_CACHE = NewTokenCache()
