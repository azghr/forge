// Package rate provides a token-bucket rate limiter.
//
// It supports configurable rate and burst, non-blocking Allow checks, and
// context-aware blocking Wait. The implementation uses a goroutine-free
// design: tokens are refilled lazily on each call based on elapsed time.
package rate

import (
	"context"
	"sync"
	"time"
)

// Limiter controls how frequently events are allowed.
//
// A zero Limiter is not usable; create one with New.
type Limiter struct {
	mu       sync.Mutex
	rate     float64
	burst    int
	tokens   float64
	lastTime time.Time
}

// New creates a limiter with the given rate (events per second) and burst
// (maximum accumulated tokens). At minimum one event per second is allowed
// even with a rate of 0.
func New(rate int, burst int) *Limiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	return &Limiter{
		rate:     float64(rate),
		burst:    burst,
		tokens:   float64(burst),
		lastTime: time.Now(),
	}
}

// Allow reports whether an event is allowed now. If not, the caller should
// wait or drop the event. Allow is safe for concurrent use.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.refill(now)

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// Wait blocks until an event is allowed or ctx is cancelled. It returns
// ctx.Err() if the context is cancelled before a token is available.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		l.refill(now)

		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}

		// Time until the next token becomes available.
		wait := time.Duration((1 - l.tokens) / l.rate * float64(time.Second))
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// Limit returns the configured rate (events per second).
func (l *Limiter) Limit() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return int(l.rate)
}

// Burst returns the configured maximum burst.
func (l *Limiter) Burst() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.burst
}

// refill adds tokens based on elapsed time since the last refill.
// Must be called with l.mu held.
func (l *Limiter) refill(now time.Time) {
	elapsed := now.Sub(l.lastTime)
	tokens := float64(elapsed) / float64(time.Second) * l.rate
	l.tokens += tokens
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}
	l.lastTime = now
}
