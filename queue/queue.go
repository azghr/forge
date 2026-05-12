// Package queue provides a generic FIFO queue backed by a ring-buffer.
// It is concurrency-safe and supports both non-blocking and context-cancellable
// dequeue operations. Zero external dependencies.
package queue

import (
	"context"
	"sync"
)

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

// Queue is a generic FIFO queue backed by a ring-buffer.
// A zero-value Queue is NOT usable; use New to create one.
type Queue[T any] struct {
	mu    sync.Mutex
	cond  *sync.Cond
	buf   []T
	head  int
	tail  int
	count int
}

// New returns an empty Queue. The initial capacity can be configured with
// WithCapacity.
func New[T any](opts ...Option) *Queue[T] {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	q := &Queue[T]{buf: make([]T, cfg.capacity)}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue adds v to the back of the queue.
func (q *Queue[T]) Enqueue(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == len(q.buf) {
		q.grow()
	}
	q.buf[q.tail] = v
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
	q.cond.Signal()
}

// Dequeue removes and returns the front element. If the queue is empty
// it returns the zero value and false.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		var zero T
		return zero, false
	}
	v := q.buf[q.head]
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	return v, true
}

// DequeueContext blocks until an element is available and then dequeues it,
// or returns zero/false if ctx is cancelled before an element arrives.
func (q *Queue[T]) DequeueContext(ctx context.Context) (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == 0 {
		// Register context cancellation to wake us up.
		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-ctx.Done():
				q.cond.Broadcast()
			case <-done:
			}
		}()

		for q.count == 0 {
			if ctx.Err() != nil {
				var zero T
				return zero, false
			}
			q.cond.Wait()
		}
	}

	v := q.buf[q.head]
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	return v, true
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count
}

// grow doubles the buffer size and repositions elements contiguously.
func (q *Queue[T]) grow() {
	newBuf := make([]T, len(q.buf)*2)
	if q.head < q.tail {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}
	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}
