# Retry Strategies for Production Systems

Transient failures are inevitable in distributed systems. A good retry strategy distinguishes a resilient service from a fragile one.

## When to Retry

Retry only for **transient** failures — network timeouts, temporary unavailability, lock contention. Do **not** retry for:
- Invalid input (4xx errors)
- Authentication failures
- Resource not found

## Exponential Backoff with Jitter

Fixed-interval retries cause thundering herd problems. Use exponential backoff with full jitter:

```go
err := retry.RetryContext(ctx, retry.RetryConfig{
    MaxTries:   3,
    InitDelay:  50 * time.Millisecond,
    Multiplier: 2.0,
    MaxDelay:   2 * time.Second,
}, func() error {
    return callExternalService()
})
```

This gives:
- **Exponential backoff**: 50ms, 100ms, 200ms
- **Full jitter**: randomizes each delay to avoid synchronization
- **Max delay cap**: prevents unbounded waits
- **Context respect**: cancellation interrupts waiting

## Choosing Configuration

| Workload | MaxTries | InitDelay | Multiplier | MaxDelay |
|---|---|---|---|---|
| In-process call | 2 | 10ms | 2.0 | 100ms |
| Local network | 3 | 50ms | 2.0 | 1s |
| External API | 3-5 | 100ms | 2.0 | 5s |
| Database | 3 | 10ms | 1.5 | 500ms |
| File system | 3 | 5ms | 2.0 | 50ms |

## Context Integration

Always pass a `context.Context` with appropriate timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := retry.RetryContext(ctx, config, fn)
```

The retry checks context before each attempt and during the backoff wait. If the context is cancelled, retry stops immediately and returns `ctx.Err()`.

## Error Classification

Not all errors should trigger a retry:

```go
err := retry.RetryContext(ctx, config, func() error {
    result, err := api.Call()
    if err != nil {
        if isTransient(err) {
            return err // triggers retry
        }
        return nil // non-nil, stops retry — wait, this won't work
    }
    return nil
})
```

Note: `RetryContext` retries on **any** non-nil error. For selective retry, wrap the function:

```go
err := retry.RetryContext(ctx, config, func() error {
    err := api.Call()
    if err != nil && !isTransient(err) {
        // Return a sentinel to stop retrying
        return &fatalError{err: err}
    }
    return err
})
```

## Key Rules

1. Always use jitter — never retry at fixed intervals
2. Cap maximum delay — unbounded backoff is worse than no retry
3. Respect context — retry must be cancellable
4. Limit retry count — infinite retries cause cascading failures
5. Log retries — `slog.Debug` for each attempt, `slog.Warn` on final failure
