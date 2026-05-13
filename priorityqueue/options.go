package priorityqueue

// Option configures Queue behaviour.
type Option func(*config)

type config struct {
	maxHeap bool
}

func defaultConfig() config {
	return config{maxHeap: false}
}

// WithMaxHeap configures the queue as a max-heap (highest priority popped
// first). The default is min-heap (lowest priority popped first).
func WithMaxHeap() Option {
	return func(c *config) {
		c.maxHeap = true
	}
}
