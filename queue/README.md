# queue

Generic FIFO queue backed by a ring-buffer. Concurrency-safe, with both
non-blocking and context-cancellable dequeue operations.

## Problem

Go provides slices and channels as building blocks, but there is no standard
FIFO queue type that combines O(1) amortised operations, concurrency safety,
and a clean API. Channels are fixed-size and lack a non-blocking dequeue;
slice-based queues require managing head/tail indices manually. This package
provides a ready-to-use generic queue with zero external dependencies.

## Quick start

```go
import "github.com/azghr/forge/queue"

q := queue.New[int]()
q.Enqueue(10)
q.Enqueue(20)
x, ok := q.Dequeue()
fmt.Println(x, ok) // 10 true
```

## API

### Functions

- **`New[T any](opts ...Option) *Queue[T]`** — create an empty queue.

### Methods

- **`Enqueue(v T)`** — add `v` to the back of the queue.
- **`Dequeue() (T, bool)`** — remove and return the front element. Returns
  `false` if the queue is empty.
- **`DequeueContext(ctx context.Context) (T, bool)`** — block until an element
  is available, then dequeue it. Returns `false` if `ctx` is cancelled.
- **`Len() int`** — number of elements currently in the queue.

### Options

- **`WithCapacity(n int) Option`** — initial ring-buffer capacity (default 8).
  The queue grows automatically when full.

## Performance

All operations are O(1) amortised. The ring-buffer avoids slice-shifting
overhead: enqueue and dequeue are simple index increments with wrapping.

Benchmark (Apple M1 Max):
- `Enqueue` + `Dequeue` pair: ~800 ns/op
- 100 enqueue + 100 dequeue batch: ~760 ns per op
