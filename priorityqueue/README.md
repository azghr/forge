# priorityqueue

Generic binary heap (min-heap or max-heap) with concurrency-safe push and pop.

## Problem

Go's standard library provides `container/heap` as an interface, but using it
requires boilerplate wrapper types and manual type assertions. This package
provides a ready-to-use generic priority queue with a clean API, zero external
dependencies, and built-in concurrency safety.

## Quick start

```go
import "github.com/azghr/forge/priorityqueue"

pq := priorityqueue.New[string]()
pq.Push(priorityqueue.Item[string]{Value: "a", Priority: 10})
pq.Push(priorityqueue.Item[string]{Value: "b", Priority: 5})

it, _ := pq.Pop()
fmt.Println(it.Value) // "b" (lower priority = higher precedence in min-heap)
```

## API

### Functions

- **`New[T any](opts ...Option) *Queue[T]`** — create an empty priority queue
  (min-heap by default).

### Types

- **`Item[T any]`** — holds a `Value T` and an `int Priority`.

### Methods

- **`Push(item Item[T])`** — insert an item into the queue.
- **`Pop() (Item[T], bool)`** — remove and return the highest-priority item.
  Returns `false` if the queue is empty.
- **`PopContext(ctx context.Context) (Item[T], bool)`** — block until an item
  is available or `ctx` is cancelled.
- **`Len() int`** — number of items in the queue.

### Options

- **`WithMaxHeap() Option`** — flip to max-heap (highest priority popped first).

## Performance

| Operation | Time   | Notes                     |
|-----------|--------|---------------------------|
| Push      | O(log n) | Single sift-up          |
| Pop       | O(log n) | Single sift-down        |
| Len       | O(1)   | Returns stored length     |

All operations are concurrency-safe via `sync.Mutex`. Benchmark (Apple M1 Max):
- Push+Pop pair: ~790 ns/op
- 100 Push+100 Pop batch: ~910 ns per op
