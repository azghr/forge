package workerpool

// Option configures a Pool.
type Option func(*options)

type options struct {
	taskBuf int
}

// WithTaskBuffer sets the size of the internal task channel buffer. The
// default is the pool size (n passed to New). A larger buffer allows
// more tasks to be queued without blocking Submit.
func WithTaskBuffer(size int) Option {
	return func(o *options) {
		o.taskBuf = size
	}
}
