# orderedset

Generic insertion-ordered set for Go with set operations (Union, Intersect).

## Problem

Go has `map[T]struct{}` for set semantics but iteration order is undefined.
When you need deduplication **and** predictable insertion order, you must
manage both a map and a slice yourself. This package wraps that pattern in a
small, reusable generic type with a clean API.

## Quick start

```go
import "github.com/azghr/forge/orderedset"

s := orderedset.New[int]()
s.Add(1); s.Add(2); s.Add(1)
fmt.Println(s.Values()) // [1 2]

s.Remove(1)
fmt.Println(s.Values()) // [2]
fmt.Println(s.Contains(2)) // true

a := orderedset.New([]int{1, 2, 3}...)
b := orderedset.New([]int{2, 3, 4}...)
a.Union(b)
fmt.Println(a.Values()) // [1 2 3 4]

a = orderedset.New([]int{1, 2, 3}...)
a.Intersect(b)
fmt.Println(a.Values()) // [2 3]
```

## API

### Type

- **`Set[T comparable]`** — insertion-ordered set of comparable elements.

### Functions

- **`New[T](elems ...T) *Set[T]`** — create a set, optionally pre-populated
  with initial elements (duplicates dropped, first-occurrence order).

### Methods

| Method | Description |
|--------|-------------|
| `Add(v T)` | Insert `v` if not already present. |
| `Remove(v T)` | Delete `v`; order of others preserved. |
| `Contains(v T) bool` | Membership test. |
| `Values() []T` | Return copy of elements in insertion order. |
| `Len() int` | Number of elements. |
| `Union(other *Set[T])` | Modify set to union with `other`. |
| `Intersect(other *Set[T])` | Modify set to intersect with `other`. |

### Error / nil semantics

- No methods return errors.
- `Union` / `Intersect` treat a nil `other` as an empty set.
- `Values` returns a copy so callers may safely mutate it.

## Performance

| Operation | Time   | Notes |
|-----------|--------|-------|
| Add       | O(1)   | Map insert / lookup |
| Contains  | O(1)   | Map lookup |
| Remove    | O(n)   | Slice shift |
| Values    | O(n)   | Copy iteration |
| Len       | O(1)   | Slice length |
| Union     | O(m)   | m = other.Len(), Append new elements |
| Intersect | O(n+m) | Iterate + lookup |

Memory: O(n) for the map and slice.

Concurrency: Safe for concurrent reads and writes through the mutex-protected
API. Cross-calls (e.g. `a.Union(b)` and `b.Intersect(a)` concurrently) may
deadlock; serialise such operations externally.
