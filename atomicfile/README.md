# atomicfile

Atomic file writes for Go using a temp-file + rename pattern.

## Problem

Writing to a file directly risks leaving partial or corrupted data if the
write is interrupted (crash, power loss, disk full). This package ensures
that a file's content is always **fully written or unchanged** — never
in-between. It uses the standard Unix pattern: write to a temporary file in
the same directory, fsync, then atomically rename over the target.

## Quick start

```go
import "github.com/azghr/forge/atomicfile"

err := atomicfile.Write("/tmp/data.txt", []byte("hello world"))
if err != nil {
    log.Fatal(err)
}
```

## API

### Functions

| Function | Description |
|----------|-------------|
| `Write(path, data, opts...)` | Write `data` atomically to `path`. |
| `WriteContext(ctx, path, data, opts...)` | Like `Write` with context cancellation. |
| `WriteReader(ctx, path, r, opts...)` | Atomically write from an `io.Reader`. |

### Types

- **`WriteError`** — returned on write failures; fields `Op` (string) and
  `Err` (underlying error). Use `errors.As` to inspect.

### Options

- **`WithFileMode(mode)`** — set file permission bits (default: `0644`).
- **`WithoutFSync()`** — skip the fsync before rename (faster, less durable).

### Error semantics

- On success: `nil`.
- On context cancellation before rename: `ErrCancelled`.
- On I/O failure: `*WriteError` wrapping the underlying error.

## Performance

| Operation | Cost | Notes |
|-----------|------|-------|
| Write     | O(n) | Full data copy to temp file + rename + optional dir fsync |
| WriteReader | O(n) | Streaming write from reader |

Benchmarks (4 KB data on Apple M1 Max):
- With fsync: ~75 µs
- Without fsync: ~30 µs

Concurrency: safe for concurrent writes to different paths. Concurrent writes
to the same path are safe but one writer's data will win (last rename wins).
