// Package workerpool provides a fixed-size worker goroutine pool for
// executing tasks concurrently and collecting results.
//
// It uses generics for type-safe task and result handling, supports
// configurable task buffering, and is safe for concurrent use.
package workerpool

import "sync"

// Pool represents a fixed-size worker pool that executes submitted tasks
// concurrently and delivers results through the Results channel.
//
// The zero Pool is not ready for use; create one with New.
//
// A Pool is safe for concurrent use. Submit may be called from multiple
// goroutines. Results must be consumed from a single goroutine after Close.
type Pool[T any] struct {
	tasks   chan func() T
	Results chan T
	wg      sync.WaitGroup
}

// New creates a pool with n worker goroutines. At most n tasks execute
// concurrently; additional submissions queue in a task buffer.
//
// Options may be provided to configure the task buffer size (defaults to n).
// The Results channel is buffered with capacity n.
func New[T any](n int, opts ...Option) *Pool[T] {
	o := &options{taskBuf: n}
	for _, fn := range opts {
		fn(o)
	}
	p := &Pool[T]{
		tasks:   make(chan func() T, o.taskBuf),
		Results: make(chan T, n),
	}
	for i := 0; i < n; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

// worker pulls tasks from the task channel, executes them, and delivers
// results to the Results channel. Results are sent directly when the
// Results buffer has space; otherwise a goroutine is spawned so the worker
// never blocks on result delivery. This preserves result ordering for
// workloads that fit within the Results buffer.
func (p *Pool[T]) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		result := task()
		select {
		case p.Results <- result:
		default:
			p.wg.Add(1)
			go func(r T) {
				defer p.wg.Done()
				p.Results <- r
			}(result)
		}
	}
}

// Submit adds a task to the pool for execution. Submit blocks if the task
// buffer is full.
//
// Submitting after Close panics; callers must ensure all submissions are
// complete before calling Close.
func (p *Pool[T]) Submit(task func() T) {
	p.tasks <- task
}

// Close shuts down the pool. It closes the task channel (preventing new
// submissions) and starts a background goroutine that waits for all workers
// and in-flight result goroutines to finish before closing the Results
// channel.
//
// The caller must consume all results from the Results channel after Close
// returns; the channel is closed automatically once every task has
// delivered its result.
//
// Close returns immediately; it does not block until tasks complete.
func (p *Pool[T]) Close() {
	close(p.tasks)
	go func() {
		p.wg.Wait()
		close(p.Results)
	}()
}
