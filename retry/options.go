package retry

import "time"

// RetryOption configures a RetryConfig using the functional options pattern.
type RetryOption func(*RetryConfig)

// WithMaxTries sets the maximum number of attempts.
func WithMaxTries(n int) RetryOption {
	return func(c *RetryConfig) {
		c.MaxTries = n
	}
}

// WithInitDelay sets the initial backoff duration.
func WithInitDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) {
		c.InitDelay = d
	}
}

// WithMultiplier sets the backoff multiplier applied after each retry.
func WithMultiplier(m float64) RetryOption {
	return func(c *RetryConfig) {
		c.Multiplier = m
	}
}

// WithMaxDelay sets an optional cap on the per-attempt backoff delay.
func WithMaxDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) {
		c.MaxDelay = d
	}
}

// NewConfig returns a RetryConfig with sensible defaults overridden by the
// supplied options. Defaults: MaxTries=3, InitDelay=100ms, Multiplier=2.0.
func NewConfig(opts ...RetryOption) RetryConfig {
	c := RetryConfig{
		MaxTries:   3,
		InitDelay:  100 * time.Millisecond,
		Multiplier: 2.0,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
