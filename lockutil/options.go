package lockutil

import "time"

// Option configures lock acquisition behaviour.
type Option func(*config)

// config holds settings for LockMutex and LockRW.
type config struct {
	interval time.Duration
}

// defaultConfig returns a config with sensible defaults.
func defaultConfig() config {
	return config{
		interval: 10 * time.Microsecond,
	}
}

// WithPollInterval sets the interval between TryLock attempts in the
// context-aware Lock functions. Smaller values mean faster lock detection at
// the cost of more CPU usage. Values below time.Microsecond are clamped.
func WithPollInterval(d time.Duration) Option {
	if d < time.Microsecond {
		d = time.Microsecond
	}
	return func(c *config) {
		c.interval = d
	}
}
