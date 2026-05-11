# sliceutil

Generic slice operations for Go — Map, Filter, Reduce, All, Any, and Chunk.

## Problem

Go's type system supports generics (since 1.18) but the standard library does
not include common higher-order slice functions. Every project ends up
reimplementing `map`, `filter`, and `reduce` by hand. This package provides a
small, stdlib-only set of generic slice utilities with a clean, idiomatic API.

## Quick start

```go
import "github.com/azghr/forge/sliceutil"

ints := []int{1, 2, 3, 4}

sqs := sliceutil.Map(ints, func(x int) int { return x * x })
// sqs == []int{1, 4, 9, 16}

evens := sliceutil.Filter(ints, func(x int) bool { return x%2 == 0 })
// evens == []int{2, 4}

sum := sliceutil.Reduce(ints, 0, func(acc, x int) int { return acc + x })
// sum == 10
```

## API

### Functions

All functions accept a slice `S ~[]E` (any concrete type whose underlying type
is `[]E`) and never mutate the input.

- **`Map[S ~[]E, E, R any](s S, f func(E) R) []R`** — apply `f` to each
  element, return new slice of results.
- **`Filter[S ~[]E, E any](s S, f func(E) bool) S`** — return new slice with
  elements where `f` returns true.
- **`Reduce[S ~[]E, E, R any](s S, init R, f func(R, E) R) R`** — accumulate
  elements left-to-right, starting with `init`.
- **`All[S ~[]E, E any](s S, f func(E) bool) bool`** — true if `f` holds for
  every element (vacuously true for empty slices). Short-circuits.
- **`Any[S ~[]E, E any](s S, f func(E) bool) bool`** — true if `f` holds for
  at least one element. Short-circuits.
- **`Chunk[S ~[]E, E any](s S, n int) []S`** — split into chunks of size `n`
  (last may be smaller). Returns nil when `n <= 0` or `s` is empty.

### Error semantics

None of the functions return errors. If the provided `f` panics, the panic
propagates to the caller.

## Performance

| Function | Time | Extra space | Notes |
|----------|------|-------------|-------|
| Map      | O(n) | O(n)        | Allocates result slice |
| Filter   | O(n) | O(n)        | Pre-allocates capacity |
| Reduce   | O(n) | O(1)        | In-place accumulator |
| All      | O(n) | O(1)        | Short-circuits on false |
| Any      | O(n) | O(1)        | Short-circuits on true |
| Chunk    | O(n) | O(n)        | Shares backing with input |

All functions are concurrency-safe (no shared state).
