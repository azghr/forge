package cache

import "time"

type config struct {
	cleanupInterval time.Duration
}

// Option configures a Cache created with NewCache.
type Option func(*config)

// WithCleanupInterval starts a background goroutine that removes expired
// entries every interval. The goroutine is stopped via Cache.Stop().
//
// Use this when you have many short-lived entries and want to reclaim
// memory between Get calls. Without this option, expired entries are
// removed lazily on Get, which is sufficient for most use cases.
func WithCleanupInterval(d time.Duration) Option {
	return func(c *config) {
		c.cleanupInterval = d
	}
}
