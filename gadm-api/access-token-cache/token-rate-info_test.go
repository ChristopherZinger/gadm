package accessTokenCache

import (
	"sync"
	"testing"
	"time"
)

type DateInfo struct {
	month int
	day   int
	year  int
}

func getCreatedAtFromDateInfo(dateInfo DateInfo) time.Time {
	return time.Now().AddDate(dateInfo.year, dateInfo.month, dateInfo.day)
}
func TestIsTokenExpiredToken(t *testing.T) {
	t.Logf("Test: isTokenExpired")

	expiredDateInfos := []DateInfo{
		{year: 0, month: -3, day: -1},
		{year: 0, month: -4, day: 0},
		{year: -1, month: 0, day: 0},
	}

	for _, expiredDate := range expiredDateInfos {
		if !isTokenExpired(getCreatedAtFromDateInfo(expiredDate)) {
			t.Errorf("failed to detect expired token: date: Y%d M%d D%d",
				expiredDate.year, expiredDate.month, expiredDate.day)
		}
	}

	datesBeforeExpiration := []DateInfo{
		{year: 0, month: -3, day: 1},
		{year: 0, month: 0, day: 0},
		{year: 0, month: -1, day: -3},
		{year: 0, month: -2, day: -29},
	}

	for _, v := range datesBeforeExpiration {
		if isTokenExpired(getCreatedAtFromDateInfo(v)) {
			t.Errorf("valid token date evaluated as expired: date: Y%d M%d D%d",
				v.year, v.month, v.day)
		}
	}
}

func TestHandleHitExpiredToken(t *testing.T) {
	t.Logf("Test: TokenRateInfo.handleHit - expired token")

	expiredDateInfos := []DateInfo{
		{year: 0, month: -3, day: -1},
		{year: 0, month: -4, day: 0},
		{year: -1, month: 0, day: 0},
	}

	for _, expiredDate := range expiredDateInfos {
		tokenRateInfo := newTokenRateInfo(getCreatedAtFromDateInfo(expiredDate))
		err := tokenRateInfo.handleHit()
		if err == nil {
			t.Errorf(
				"expected expired token error for creation date: Y%d M%d D%d",
				expiredDate.year, expiredDate.month, expiredDate.day)
		} else if err.Error() != TokenExpiredMsg {
			t.Errorf(
				"unexpected error message: '%s'. Expected '%s'",
				err.Error(), TokenExpiredMsg)
		}
	}
}

func TestHandleHitRateLimitExceeded(t *testing.T) {
	t.Logf("Test: TokenRateInfo.handleHit - rate limit exceeded")

	createdAt := getCreatedAtFromDateInfo(DateInfo{year: 0, month: 0, day: 0})
	tokenRateInfo := newTokenRateInfo(createdAt)

	_NUM_HITS_PER_RATE_LIMIT := 10

	for i := 0; i < _NUM_HITS_PER_RATE_LIMIT+1; i++ {
		err := tokenRateInfo.handleHit()
		if i < _NUM_HITS_PER_RATE_LIMIT {
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}
		} else {
			if err == nil {
				t.Errorf("expected rate limit exceeded error but got nil")
			} else if err.Error() != RateLimitExceededMsg {
				t.Errorf(
					"unexpected error message: '%s'. Expected '%s'",
					err.Error(), RateLimitExceededMsg)
			}
		}
	}
}

func TestHandleHitRateLimitNotExceeded(t *testing.T) {
	t.Logf("Test: TokenRateInfo.handleHit - hit history is cleared accordingly between hits")

	createdAt := getCreatedAtFromDateInfo(DateInfo{year: 0, month: 0, day: 0})
	tokenRateInfo := newTokenRateInfo(createdAt)

	waitDurationBetweenHits := RATE_LIMIT_DURATION / (NUM_HITS_PER_RATE_LIMIT)

	for i := 0; i < 20; i++ {
		err := tokenRateInfo.handleHit()
		time.Sleep(waitDurationBetweenHits)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			break
		}
	}
}

func TestHandleHitConcurrency(t *testing.T) {
	t.Logf("Test: TokenRateInfo.handleHit - concurrent access and locking")

	createdAt := getCreatedAtFromDateInfo(DateInfo{year: 0, month: 0, day: 0})
	tokenRateInfo := newTokenRateInfo(createdAt)

	var wg sync.WaitGroup
	numCalls := 20
	errCh := make(chan error, numCalls)

	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := tokenRateInfo.handleHit()
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	var rateLimitExceededCount, numSuccessfulHits, otherErrCount int
	for err := range errCh {
		if err == nil {
			numSuccessfulHits++
		} else if err.Error() == RateLimitExceededMsg {
			rateLimitExceededCount++
		} else {
			otherErrCount++
			t.Errorf("unexpected error: %v", err)
		}
	}

	if numSuccessfulHits > NUM_HITS_PER_RATE_LIMIT {
		t.Errorf("expected at most %d successful hits, got %d", NUM_HITS_PER_RATE_LIMIT, numSuccessfulHits)
	}
	if rateLimitExceededCount != numCalls-NUM_HITS_PER_RATE_LIMIT {
		t.Errorf("expected %d rate limit exceeded errors, got %d", numCalls-NUM_HITS_PER_RATE_LIMIT, rateLimitExceededCount)
	}
	if otherErrCount > 0 {
		t.Errorf("unexpected errors occurred: %d", otherErrCount)
	}
}
