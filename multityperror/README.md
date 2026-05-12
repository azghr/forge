# multityperror

Aggregate multiple error values into one.

## Problem

When collecting errors from multiple operations (e.g. validation, batch
processing), you need a single error value that combines all failures.
Go 1.20's `errors.Join` is static; `multityperror.MultiError` supports
incremental accumulation with nil-skipping, custom separators, and full
`errors.Is`/`errors.As` support.

Features:
- **Incremental** — append errors as they occur.
- **Nil-safe** — nil errors are silently skipped.
- **Custom separator** — change the message glue (default `"; "`).
- **`errors.Is` / `errors.As`** — works through Go's error chain.
- **Concurrency-safe** — all methods are goroutine-safe.

## Quick start

```go
var me multityperror.MultiError
me.Append(fmt.Errorf("first"))
me.Append(nil)              // skipped
me.Append(fmt.Errorf("second"))

if !me.IsEmpty() {
    fmt.Println(me.Error()) // "first; second"
}
```

## API

### Types

- **`MultiError`** — accumulates errors.

### Functions

- **`New(opts ...Option) *MultiError`** — creates a MultiError with options.

### Methods

- **`Append(err error)`** — adds an error; nil values are skipped.
- **`Error() string`** — returns all error messages joined by the separator.
- **`IsEmpty() bool`** — reports whether no non-nil errors exist.
- **`Len() int`** — number of non-nil errors appended.
- **`Unwrap() []error`** — returns the error list for `errors.Is`/`errors.As`.
- **`Errors() []error`** — returns a copy of the error slice.

### Options

- **`WithSeparator(sep string)`** — sets the separator between error messages
  (default `"; "`).

## Performance

- `Append`: O(1) amortized, no allocations beyond slice growth.
- `Error`: O(n) message concatenation.
- All operations guarded by `sync.Mutex`.
