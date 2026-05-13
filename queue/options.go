package queue

// Option configures Queue behaviour.
type Option func(*config)

type config struct {
	capacity int
}

func defaultConfig() config {
	return config{capacity: 8}
}

// WithCapacity sets the initial ring-buffer capacity (default 8).
// The queue grows automatically when full; this merely avoids early resizes.
func WithCapacity(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.capacity = n
		}
	}
}
