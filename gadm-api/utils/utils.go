package utils

import (
	"context"
	"fmt"
	"time"
)

func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ArrayToStrings(args ...interface{}) []string {
	var result []string
	for _, item := range args {
		result = append(result, fmt.Sprintf("%s", item))
	}
	return result
}

func Sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func Retry[T any](
	ctx context.Context,
	fn func(ctx context.Context) (T, error),
	retries int,
	initialBackoff time.Duration,
	maxBackoff time.Duration,
) (T, error) {
	var result T
	if retries < 0 {
		retries = 0
	}

	var lastErr error
	backoff := initialBackoff
	for attempt := 0; attempt <= retries; attempt++ {
		if err := ctx.Err(); err != nil {
			return result, err
		}

		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		lastErr = err
		if attempt == retries {
			break
		}

		if err := Sleep(ctx, backoff); err != nil {
			return result, err
		}

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return result, fmt.Errorf("retry_exhausted: attempts=%d: %w", retries+1, lastErr)
}
