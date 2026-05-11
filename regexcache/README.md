# regexcache

A concurrency-safe cache for compiled Go regex patterns. Avoids the overhead
of `regexp.Compile` when the same patterns are used repeatedly.

## Problem

Calling `regexp.Compile` or `regexp.MustCompile` repeatedly with the same
pattern is wasteful — compilation involves parsing the pattern and building
the internal NFA/DFA. In hot paths, this adds unnecessary latency and
allocation. `regexcache.Cache` stores compiled `*regexp.Regexp` values so
compilation happens once per unique pattern.

## Quick start

```go
import "github.com/azghr/forge/regexcache"

cache := regexcache.New()

// Compile once, cache forever
r := cache.MustCompile("^a.*z$")
fmt.Println(r.MatchString("abz")) // true

// Second lookup returns the cached value (same pointer)
r2 := cache.MustCompile("^a.*z$")
fmt.Println(r == r2) // true

// Handle invalid patterns without panicking
_, err := cache.Compile("(")
fmt.Println(err) // error parsing regexp
```

## API

### Types

```go
type Cache struct { ... }
```

### Functions

| Constructor | Description |
|-------------|-------------|
| `New(opts ...Option) *Cache` | Create a new cache with optional configuration. |

### Methods

| Method | Description |
|--------|-------------|
| `(*Cache) Compile(pattern) (*Regexp, error)` | Compile or return cached pattern. |
| `(*Cache) MustCompile(pattern) *Regexp` | Like Compile but panics on error. |

### Options

- **`WithMaxSize(n)`** — limit the cache to `n` entries (default: unlimited).

### Error semantics

- `Compile` returns a `*regexp.Regexp` and an error. The error is non-nil when
  the pattern is invalid (same errors as `regexp.Compile`).
- `MustCompile` panics on invalid patterns, mirroring `regexp.MustCompile`.

## Performance

| Operation | Complexity | Notes |
|-----------|------------|-------|
| First `Compile` | O(pattern compilation) | Compiles and stores once. |
| Subsequent lookups | O(1) | Map lookup with read lock. |
| `WithMaxSize` eviction | O(1) | Evicts one entry when limit reached. |

Concurrent reads do not block each other (sync.RWMutex). Concurrent first-time
compilation for the same pattern is serialized on write lock.

All cache operations are safe for concurrent use by multiple goroutines.
