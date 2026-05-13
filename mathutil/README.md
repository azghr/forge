# mathutil

Small helper functions missing from the standard `math` package.

## Problem

Go's `math` package provides trigonometric, logarithmic, and other
transcendental functions, but omits simple everyday helpers like clamping a
value to a range, determining the sign of a number, or linear interpolation.
This package fills those gaps with zero dependencies.

## Quick start

```go
import "github.com/azghr/forge/mathutil"

c := mathutil.Clamp(10, 0, 5)   // 5
s := mathutil.Sign(-3.2)         // -1
l := mathutil.Lerp(0, 10, 0.5)  // 5
g := mathutil.GCD(8, 12)        // 4
```

## API

### Functions

- **`Clamp(x, lo, hi float64) float64`** — confine `x` to `[lo, hi]`.
- **`Sign(x float64) float64`** — return `-1`, `0`, or `+1`.
- **`Lerp(a, b, t float64) float64`** — linear interpolation (`t` in `[0,1]`).
- **`GCD(a, b int64) int64`** — greatest common divisor (Euclidean algorithm).
- **`ApproxEqual(a, b float64) bool`** — tolerant floating-point comparison
  using `DefaultEpsilon` (1e-9).
- **`ApproxEqualEpsilon(a, b, eps float64) bool`** — tolerant comparison with
  a custom epsilon. If `eps <= 0`, `DefaultEpsilon` is used.

### Constants

- **`DefaultEpsilon`** = `1e-9` — default tolerance for `ApproxEqual`.

## Performance

All functions are pure arithmetic — O(1), no allocations, no branching beyond
the minimum needed for correctness. `GCD` uses the Euclidean algorithm (O(log
min(a,b))). `ApproxEqual` is a single subtraction and comparison.
