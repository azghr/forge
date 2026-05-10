// Package retry provides simple retry logic with exponential backoff and jitter.
//
// It is designed for transient errors (network calls, file locks) and follows
// a minimal stdlib-first approach with no external dependencies.
//
// The core function is RetryContext, which executes a user-supplied function
// up to MaxTries times with exponential backoff and full-jitter between
// attempts. Context cancellation is respected at every stage.
package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// RetryConfig configures the retry behaviour.
type RetryConfig struct {
	// MaxTries is the maximum number of attempts (defaults to 1 when <= 0).
	MaxTries int

	// InitDelay is the initial backoff duration before the first retry.
	InitDelay time.Duration

	// Multiplier is applied to the delay after each attempt. A value of 2.0
	// produces classic exponential backoff.
	Multiplier float64

	// MaxDelay is an optional cap on the per-attempt delay. When set to a
	// positive value, backoff never exceeds this duration.
	MaxDelay time.Duration
}

// RetryContext executes fn up to MaxTries times with exponential backoff and
// full-jitter between attempts.
//
// Before each attempt and during the backoff sleep ctx is checked. If ctx is
// done the function returns ctx.Err() immediately.
//
// If fn returns nil the function returns nil. If all attempts fail the last
// error returned by fn is returned.
//
// Performance: O(MaxTries) calls to fn. Backoff sleeps use time.Sleep via a
// time.Timer so the goroutine can be interrupted by context cancellation.
func RetryContext(ctx context.Context, config RetryConfig, fn func() error) error {
	maxTries := config.MaxTries
	if maxTries < 1 {
		maxTries = 1
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var lastErr error
	for i := 0; i < maxTries; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(); err != nil {
			lastErr = err
		} else {
			return nil
		}

		if i == maxTries-1 {
			break
		}

		delay := float64(config.InitDelay) * math.Pow(config.Multiplier, float64(i))
		if config.MaxDelay > 0 {
			if max := float64(config.MaxDelay); delay > max {
				delay = max
			}
		}

		if delay > 0 {
			delay = rng.Float64() * delay
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(delay)):
		}
	}

	return lastErr
}
