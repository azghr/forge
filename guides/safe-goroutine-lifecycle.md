# Safe Goroutine Lifecycle

Every goroutine your application starts must be accounted for. Leaked goroutines cause resource exhaustion, memory growth, and hard-to-debug production issues.

## Principles

1. **Know when every goroutine exits** — use `sync.WaitGroup` or channels
2. **Provide cancellation** — every goroutine should accept a `context.Context`
3. **Never start a goroutine you can't stop** — if there's no shutdown path, there's a leak

## Worker Lifecycle

```go
type Worker struct {
    stopCh chan struct{}
    wg     sync.WaitGroup
}

func NewWorker() *Worker {
    return &Worker{stopCh: make(chan struct{})}
}

func (w *Worker) Start(ctx context.Context) {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        for {
            select {
            case <-ctx.Done():
                return
            case <-w.stopCh:
                return
            default:
                // do work
            }
        }
    }()
}

func (w *Worker) Stop() {
    close(w.stopCh)
    w.wg.Wait()
}
```

## Using Workerpool

Forge's `workerpool` handles lifecycle management for you:

```go
pool := workerpool.New[int](3)
defer pool.Close()

pool.Submit(func() int { return doWork() })

// Wait for all results
for result := range pool.Results {
    fmt.Println(result)
}
```

The pool manages worker goroutines, task distribution, and result collection. Close stops accepting tasks and waits for in-flight work to complete.

## Context Cancellation Everywhere

Every blocking operation should respect context cancellation:

```go
func (q *Queue[T]) DequeueContext(ctx context.Context) (T, bool) {
    // Context cancellation wakes the goroutine
    // even when blocked on condition variable
}
```

This ensures a single signal can cascade cancellation through all goroutines.

## Testing Goroutine Lifecycles

Always test that goroutines actually exit:

```go
func TestWorker_Shutdown(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    w := NewWorker()
    w.Start(ctx)
    cancel()
    // Worker should exit within reasonable time
}
```

## Key Rules

1. Every `go` call must have a corresponding shutdown mechanism
2. Use `sync.WaitGroup` to wait for goroutine completion
3. Never use `time.Sleep` for coordination
4. Prefer `pool.Close()` + `pool.Results` over manual goroutine management
5. Test that goroutines exit cleanly under cancellation
