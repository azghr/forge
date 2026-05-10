# retry

Simple retry logic with exponential backoff and full-jitter for Go.

## Problem

Transient errors (network hiccups, file locks, rate-limit 503s) are common in
distributed systems. Naively retrying immediately often makes things worse.
This package provides a minimal, stdlib-only retry loop that:

- Waits with **exponential backoff** between attempts.
- Adds **full-jitter** to avoid thundering-herd problems.
- Respects **context cancellation** at every step (before calls and during
  backoff sleeps).
- Returns the **last error** when all attempts are exhausted.

## Quick start

```go
import "github.com/azghr/forge/retry"

config := retry.RetryConfig{
    MaxTries:   3,
    InitDelay:  100 * time.Millisecond,
    Multiplier: 2.0,
}

err := retry.RetryContext(context.Background(), config, func() error {
    resp, err := http.Get("http://example.com")
    if err != nil || resp.StatusCode >= 500 {
        return fmt.Errorf("try again")
    }
    resp.Body.Close()
    return nil
})
```

## API

### Types

- **`RetryConfig`** – fields `MaxTries`, `InitDelay`, `Multiplier`, `MaxDelay`.
- **`RetryOption`** – functional option for `NewConfig`.

### Functions

- **`RetryContext(ctx, config, fn)`** – core retry loop.
- **`NewConfig(opts ...)`** – construct a `RetryConfig` with defaults (3 tries,
  100 ms, 2× multiplier) overridden by options.
- **`WithMaxTries(n)`**, **`WithInitDelay(d)`**, **`WithMultiplier(m)`**,
  **`WithMaxDelay(d)`** – functional option helpers.

### Error semantics

- **Success:** returns `nil`.
- **Exhaustion:** returns the last error from `fn`.
- **Cancellation:** returns `ctx.Err()` (`context.Canceled` or
  `context.DeadlineExceeded`). Use `errors.Is` to distinguish.

## Performance

- O(MaxTries) calls to `fn`.
- Backoff sleeping uses `time.After` (interruptible via context).
- Per-call `*rand.Rand` source avoids global lock contention.
- No allocations outside the per-call RNG source.
