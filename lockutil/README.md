# lockutil

Non-blocking `TryLock` and context-aware blocking `Lock` for `sync.Mutex` and `sync.RWMutex`.

## Problem

Go 1.18 added `Mutex.TryLock()` and `RWMutex.TryRLock()`, but the standard
library does not provide a way to block on lock acquisition with a deadline or
cancellation. This package exposes both through a uniform API and adds
context-cancellable lock acquisition for callers that need to bound wait time.

## Quick start

```go
import "github.com/azghr/forge/lockutil"

var mu sync.Mutex
if lockutil.TryLockMutex(&mu) {
    defer mu.Unlock()
    // critical section
} else {
    fmt.Println("Mutex busy")
}
```

Context-aware blocking:

```go
var mu sync.Mutex
mu.Lock()

ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

if lockutil.LockMutex(ctx, &mu) {
    defer mu.Unlock()
    // acquired after wait
}
```

## API

### Functions

- **`TryLockMutex(mu *sync.Mutex) bool`** — non-blocking mutex lock.
- **`TryLockRW(rw *sync.RWMutex) bool`** — non-blocking read lock.
- **`LockMutex(ctx, mu, opts...) bool`** — blocking mutex lock, cancellable via `ctx`.
- **`LockRW(ctx, rw, opts...) bool`** — blocking read lock, cancellable via `ctx`.

### Options

- **`WithPollInterval(d time.Duration)`** — interval between `TryLock` retries in
  the context-aware functions (default: 10 µs).

## Performance

- `TryLockMutex` / `TryLockRW` are direct delegations to the stdlib — no
  allocations, O(1).
- `LockMutex` / `LockRW` poll in a tight loop with `time.After`. Poll interval
  is configurable. The default (10 µs) balances CPU usage with latency.
- All operations avoid global state and are safe for concurrent use.
