# atomicfile

Atomic file operations for Go: writing or replacing files without leaving
partial data on failure.

## Problem

Writing to a file directly risks leaving partial or corrupted data if the
write is interrupted (crash, power loss, disk full). This package wraps the
standard Unix pattern — write to a temporary file, fsync, then atomically
rename over the target — so that the destination is either fully updated or
fully unchanged.

On POSIX, `os.Rename` is atomic. On Windows, `os.Rename` is not fully
atomic; consider platform-specific calls (`MoveFileEx`) if strict atomicity
is needed there.

## Quick start

```go
import "github.com/azghr/forge/atomicfile"

data := bytes.NewBufferString("important")
if err := atomicfile.WriteFile("/tmp/config.txt", data); err != nil {
    log.Fatal(err)
}
// config.txt is fully written or untouched.

// Atomic replacement:
atomicfile.ReplaceFile("/tmp/new.txt", "/tmp/config.txt")
```

## API

### Functions

| Function | Description |
|----------|-------------|
| `WriteFile(filename, r, opts...)` | Atomically write `io.Reader` to `filename`. |
| `ReplaceFile(source, dest)` | Atomically replace `dest` with `source` file. |

### Options

- **`WithFileMode(mode)`** — set file permission bits (default: `0644`).
- **`WithoutFSync()`** — skip fsync before rename (faster, less durable).

### Error semantics

- On success: `nil`.
- On I/O failure: `*WriteError` wrapping the underlying error. `Op` identifies
  the failing step (`"create"`, `"write"`, `"fsync"`, `"close"`, `"rename"`,
  `"sync-dir"`). Use `errors.As` to inspect.
- `WriteFile` removes the temp file on failure, leaving the original intact.

## Performance

- **WriteFile** — O(n) in file size. Data is copied once (temp write) then
  renamed (metadata only).
- **ReplaceFile** — O(1) (just a rename + optional dir fsync).

Benchmarks (4 KB data on Apple M1 Max):
- With fsync: ~10 ms
- Without fsync: ~5 ms

Concurrency: safe for concurrent writes to different paths. Concurrent writes
to the same path are safe but one writer's data will win (last rename wins).

## Cross-platform

On POSIX, rename is atomic. On Windows, `os.Rename` is not fully atomic;
consider Windows-specific syscalls if strict atomicity is required.
