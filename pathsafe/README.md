# pathsafe

Safely join paths to prevent directory traversal in Go.

## Problem

When accepting user-controlled relative paths and joining them to a base
directory, naive `filepath.Join(base, rel)` can produce a path outside base
via `../` traversal. This is a common security vulnerability.

`pathsafe.SafeJoin` resolves this by:

- Resolving `base` to an absolute, cleaned path.
- Joining `rel` and cleaning the result.
- Verifying the result is within base (equal to base or a subpath).
- Returning `ErrOutsideBase` if traversal is detected.

## Quick start

```go
import "github.com/azghr/forge/pathsafe"

safe, err := pathsafe.SafeJoin("/home/user", "docs/report.pdf")
// safe == "/home/user/docs/report.pdf"

_, err = pathsafe.SafeJoin("/home/user", "../etc/passwd")
// err == pathsafe.ErrOutsideBase
```

## API

### Functions

- **`SafeJoin(base, rel string, opts ...Option) (string, error)`** – joins base
  and rel, ensuring the result is within base. Returns cleaned absolute path or
  `ErrOutsideBase`. Options can enable symlink resolution.

### Options

- **`AllowSymlinkFollow()`** – resolves symlinks in both base and joined
  paths before the containment check. Use this to prevent symlink-based
  traversal attacks.

### Errors

- **`ErrOutsideBase`** – returned when the joined path is outside the base
  directory. Use `errors.Is(err, pathsafe.ErrOutsideBase)` to check.

## Performance

- O(path length) – uses `filepath.Abs`, `filepath.Clean`, and optionally
  `filepath.EvalSymlinks`.
- No allocations beyond the input string lengths and filesystem calls.
- Concurrency-safe: zero global state, no shared mutexes.
