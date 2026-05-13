# cache

Generic in-memory TTL cache for Go.

## Problem

Repeated computations or lookups for the same data waste resources. A cache
stores key/value pairs with optional expiration so that once a result is
computed it can be reused until it is stale.

`cache` provides:

- **Generic** `Cache[K, V]` — works with any comparable key and any value type.
- **TTL expiration** — values expire after a configurable duration (lazy
  removal on Get, optional periodic cleanup).
- **No-expiration** mode — set TTL to 0 for values that live forever.
- **GetOrLoad** — atomically check-then-load a value, with context support.

## Quick start

```go
import "github.com/azghr/forge/cache"

c := cache.New[string, int](5 * time.Second)
c.Set("x", 42)

v, ok := c.Get("x")   // 42, true

time.Sleep(6 * time.Second)
_, ok = c.Get("x")     // false (expired)
```

```go
// No-expiration mode
c := cache.New[int, string](0)
c.Set(1, "one")
v, ok := c.Get(1)      // "one", true
```

## API

### Types

- **`Cache[K comparable, V any]`** — the cache.

### Functions

- **`New[K, V](ttl time.Duration, opts ...Option) *Cache[K, V]`** —
  creates a cache. If ttl > 0, values expire after ttl.

### Methods

- **`Set(key K, value V)`** — stores a value (resets TTL if key exists).
- **`Get(key K) (V, bool)`** — retrieves a value; bool is false if missing
  or expired.
- **`Delete(key K)`** — removes a key.
- **`Len() int`** — number of items currently stored.
- **`GetOrLoad(ctx, key, loader func(context.Context) (V, error)) (V, error)`**
  — returns cached value or calls loader to produce and cache it.

### Options

- **`WithCleanupInterval(d time.Duration)`** — starts a background goroutine
  that evicts expired entries periodically.
- **`WithOnEvict(fn func(K, V))`** — callback fired when an entry is evicted
  (via expiration, deletion, or eviction from an option).

## Performance

- Get/Set/Delete are O(1) map operations guarded by `sync.RWMutex`.
- TTL uses `time.Now().UnixNano()` comparison on each Get.
- `Len()` acquires only a read lock.
- No goroutines are started unless `WithCleanupInterval` is provided.
- Memory: O(n) where n is the number of stored items.
