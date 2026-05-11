# option

A generic **Option (Maybe)** type for Go — a container that may hold a value or none.

## Problem

In Go, optional values are typically represented with pointers (nil means absent)
or with multi-return `(value, ok)`. Both approaches have downsides:

- **Pointers** are ambiguous: nil can mean "not found" or "not set", and you
  need to remember to nil-check everywhere.
- **Multi-return** is stateless — you can't pass an optional value around as a
  single value, compose them, or chain operations.

`option.Option[T]` wraps both the value and the presence bool into one
immutable, concurrency-safe value.

## Quick start

```go
import "github.com/azghr/forge/option"

func find(m map[string]int, key string) option.Option[int] {
    if v, ok := m[key]; ok {
        return option.Some(v)
    }
    return option.None[int]()
}

o := find(map[string]int{"a": 1}, "a")
if v, ok := o.Unwrap(); ok {
    fmt.Println(v) // 1
}
```

## API

### Constructors

- **`Some[T any](v T) Option[T]`** — wrap a value.
- **`None[T any]() Option[T]`** — empty option.

### Methods

- **`(o Option[T]) IsSome() bool`** — true if the option holds a value.
- **`(o Option[T]) Unwrap() (T, bool)`** — returns `(value, true)` if Some,
  `(zero, false)` if None. Never panics.
- **`(o Option[T]) Must() T`** — returns the contained value, **panics** if
  None. Use only when you are certain a value exists.

### Error semantics

`Unwrap` never panics or errors; it returns false for None. `Must` panics on
None. No other method panics.

## Performance

| Method  | Time | Space | Notes       |
|---------|------|-------|-------------|
| Some    | O(1) | O(1)  | Stack-local |
| None    | O(1) | O(1)  | Stack-local |
| IsSome  | O(1) | O(1)  | Inlineable  |
| Unwrap  | O(1) | O(1)  | Inlineable  |
| Must    | O(1) | O(1)  | Inlineable  |

The type stores one `bool` and one `T`, no heap allocation. All methods are
concurrency-safe (pure reads on an immutable struct).
