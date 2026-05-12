// Package lockutil provides non-blocking TryLock functions for sync.Mutex and
// sync.RWMutex, as well as context-aware blocking Lock variants with
// configurable polling.
//
// Go 1.18 introduced Mutex.TryLock and RWMutex.TryRLock in the standard
// library. This package exposes those operations through a uniform API and
// adds context-cancellable lock acquisition for callers that need to bound
// wait time.
package lockutil

import (
	"context"
	"sync"
	"time"
)

// TryLockMutex attempts to lock mu without blocking.
// Returns true if the lock was acquired, false if it is already held.
func TryLockMutex(mu *sync.Mutex) bool {
	return mu.TryLock()
}

// TryLockRW attempts to acquire a read lock on rw without blocking.
// Returns true if the read lock was acquired, false if a write lock is held.
func TryLockRW(rw *sync.RWMutex) bool {
	return rw.TryRLock()
}

// LockMutex acquires a lock on mu, blocking until successful or ctx is done.
// opts can specify a polling interval via WithPollInterval.
// Returns true if the lock was acquired, false if ctx expired.
func LockMutex(ctx context.Context, mu *sync.Mutex, opts ...Option) bool {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return pollLock(ctx, func() bool { return mu.TryLock() }, cfg.interval)
}

// LockRW acquires a read lock on rw, blocking until successful or ctx is done.
// opts can specify a polling interval via WithPollInterval.
// Returns true if the read lock was acquired, false if ctx expired.
func LockRW(ctx context.Context, rw *sync.RWMutex, opts ...Option) bool {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return pollLock(ctx, func() bool { return rw.TryRLock() }, cfg.interval)
}

// pollLock repeatedly calls try until it returns true or ctx is done.
func pollLock(ctx context.Context, try func() bool, interval time.Duration) bool {
	for {
		if try() {
			return true
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(interval):
		}
	}
}
