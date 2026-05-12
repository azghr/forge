// Package priorityqueue implements a generic binary heap (min-heap by default)
// with concurrency-safe push and pop operations.
//
// Zero external dependencies; uses a manual heap implementation for
// predictable performance.
package priorityqueue

import (
	"context"
	"sync"
)

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

// Item holds a value with an associated priority. Lower Priority values
// have higher precedence in a min-heap (the default).
type Item[T any] struct {
	Value    T
	Priority int
}

// Queue is a generic priority queue backed by a binary heap.
// A zero-value Queue is NOT usable; use New to create one.
type Queue[T any] struct {
	mu      sync.Mutex
	cond    *sync.Cond
	items   []Item[T]
	maxHeap bool
}

// New returns an empty priority queue. By default it is a min-heap; use
// WithMaxHeap to reverse the ordering.
func New[T any](opts ...Option) *Queue[T] {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	q := &Queue[T]{maxHeap: cfg.maxHeap}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push adds an item to the queue.
func (q *Queue[T]) Push(item Item[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, item)
	q.siftUp(len(q.items) - 1)
	q.cond.Signal()
}

// Pop removes and returns the highest-priority item (smallest priority in
// a min-heap, largest in a max-heap). Returns false if the queue is empty.
func (q *Queue[T]) Pop() (Item[T], bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return Item[T]{}, false
	}
	v := q.items[0]
	q.items[0] = q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	if len(q.items) > 0 {
		q.siftDown(0)
	}
	return v, true
}

// PopContext blocks until an item is available and pops it, or returns
// the zero value and false if ctx is cancelled.
func (q *Queue[T]) PopContext(ctx context.Context) (Item[T], bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		done := make(chan struct{})
		defer close(done)

		go func() {
			select {
			case <-ctx.Done():
				q.cond.Broadcast()
			case <-done:
			}
		}()

		for len(q.items) == 0 {
			if ctx.Err() != nil {
				return Item[T]{}, false
			}
			q.cond.Wait()
		}
	}

	v := q.items[0]
	q.items[0] = q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	if len(q.items) > 0 {
		q.siftDown(0)
	}
	return v, true
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// --- binary heap helpers ---

func (q *Queue[T]) less(i, j int) bool {
	if q.maxHeap {
		return q.items[i].Priority > q.items[j].Priority
	}
	return q.items[i].Priority < q.items[j].Priority
}

func (q *Queue[T]) siftUp(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if !q.less(i, p) {
			break
		}
		q.items[i], q.items[p] = q.items[p], q.items[i]
		i = p
	}
}

func (q *Queue[T]) siftDown(i int) {
	n := len(q.items)
	for {
		smallest := i
		l := 2*i + 1
		r := 2*i + 2
		if l < n && q.less(l, smallest) {
			smallest = l
		}
		if r < n && q.less(r, smallest) {
			smallest = r
		}
		if smallest == i {
			break
		}
		q.items[i], q.items[smallest] = q.items[smallest], q.items[i]
		i = smallest
	}
}
