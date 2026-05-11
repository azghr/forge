# workerpool

A fixed-size worker goroutine pool for executing tasks concurrently and collecting results. Limits parallelism to a configurable number of workers.

## Problem

Go makes it easy to spawn goroutines, but unconstrained concurrency can overwhelm system resources (goroutine stack, CPU, file descriptors). A worker pool bounds the number of concurrent executions while still allowing an arbitrary number of tasks.

## Quick start

```go
import "github.com/azghr/forge/workerpool"

pool := workerpool.NewPool(3)

for i := 0; i < 5; i++ {
    i := i
    pool.Submit(func() interface{} {
        time.Sleep(time.Duration(i) * 10 * time.Millisecond)
        return i * 2
    })
}

pool.Close()
for res := range pool.Results {
    fmt.Println(res.(int))
}
```

## API

### Functions

- **`NewPool(n int) *Pool`** — create a pool with n workers. The task channel
  buffer defaults to n (configurable via `WithTaskBuffer`).

### Types

```go
type Pool struct {
    Results chan interface{}  // read results from here after Close
}
```

### Methods

- **`(*Pool) Submit(task func() interface{})`** — enqueue a task. Blocks if the
  task buffer is full. Panics if called after Close.
- **`(*Pool) Close()`** — prevent new submissions and start draining results.
  Returns immediately. The Results channel is automatically closed once all
  tasks have delivered their results.

### Error semantics

Neither `NewPool` nor `Submit` return errors. If a task function panics, the
panic crashes the per-task goroutine; tasks that may panic should recover
internally. Submitting after Close panics (sending on a closed channel).

## Performance

| Operation        | Complexity | Notes                        |
|------------------|------------|------------------------------|
| NewPool          | O(n)       | Starts n worker goroutines   |
| Submit           | O(1)       | Blocks if buffer full        |
| Close            | O(1)       | Returns immediately          |
| Per-task overhead| O(1)       | One additional goroutine     |

The pool creates n permanent worker goroutines plus one temporary goroutine
per submitted task (for non-blocking result delivery). Task channel buffer
defaults to n; use `WithTaskBuffer` for burstier workloads.

Pool is safe for concurrent use: Submit can be called from multiple goroutines.
