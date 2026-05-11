# stopwatch

A simple stopwatch for benchmarking Go code blocks with start/stop/reset controls.

## Problem

Measuring elapsed time in Go typically requires manually tracking `time.Now()`
calls and subtracting timestamps. For repeated measurements (e.g., aggregating
loop iterations), you need to accumulate durations yourself. `stopwatch.Stopwatch`
packages this into a reusable, zero-allocation value type.

## Quick start

```go
import "github.com/azghr/forge/stopwatch"

var sw stopwatch.Stopwatch
sw.Start()
time.Sleep(10 * time.Millisecond)
sw.Stop()
fmt.Println(sw.Elapsed()) // ~10ms

// Cumulative
sw.Start()
time.Sleep(5 * time.Millisecond)
sw.Stop()
fmt.Println(sw.Elapsed()) // ~15ms total

// Reset
sw.Reset()
fmt.Println(sw.Elapsed()) // 0s
```

## API

### Methods

- **`(*Stopwatch) Start()`** — begin or restart the timer. Consecutive calls
  reset the start time.
- **`(*Stopwatch) Stop()`** — halt and accumulate elapsed time. No-op if not
  running.
- **`(*Stopwatch) Reset()`** — zero elapsed time and stop. Safe to call in any
  state.
- **`(*Stopwatch) Elapsed() time.Duration`** — return total elapsed time.
  Includes partial time if currently running.

### Error semantics

No methods return errors or panic. Calling Stop on a non-running stopwatch is
a no-op.

## Performance

| Method  | Allocations | Notes                      |
|---------|-------------|----------------------------|
| Start   | 0           | Reads `time.Now`           |
| Stop    | 0           | Reads `time.Since`         |
| Elapsed | 0           | Duration addition only     |
| Reset   | 0           | Zeroes three fields        |

All methods are O(1) with zero heap allocation. Precision is limited by
`time.Now` resolution (~ns on modern hardware).

Stopwatch is **not concurrency-safe**; external synchronization is required
for concurrent use.
