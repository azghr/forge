# shellquote

Shell-safe string quoting for Go. Escape strings for safe use as shell
arguments on POSIX or Windows without injection risks.

## Problem

Passing unsanitized strings to shell commands (e.g., via `os/exec`) can lead
to argument injection or command execution. A filename like `file; rm -rf /`
would be interpreted as multiple tokens by a POSIX shell. This package escapes
such strings so they are treated as a single, literal argument.

## Quick start

```go
import "github.com/azghr/forge/shellquote"

// POSIX single-quote wrapping
safe := shellquote.Quote("file; rm -rf /")
// safe = 'file; rm -rf /'

// Build a command line from args
cmd := shellquote.QuoteCommand([]string{"ls", "-l", "my file"})
// cmd = ls -l 'my file'

// Windows cmd.exe quoting
win := shellquote.QuoteWindows(`C:\path with spaces\`)
// win = "C:\path with spaces\\"
```

## API

### Functions

| Function | Description |
|----------|-------------|
| `Quote(s) string` | Escape `s` for POSIX shells (single-quote wrapping). |
| `QuoteCommand(args) string` | Quote each arg and join into a command line. |
| `QuoteWindows(s) string` | Escape `s` for Windows cmd.exe. |

### Error semantics

No functions return errors. All inputs are valid strings; the output is always
a safe, quoted representation.

## Performance

- **Quote** — O(n). Single pass with `strings.ReplaceAll`.
- **QuoteCommand** — O(n) across all args. Builder-based join avoids extra
  allocations.
- **QuoteWindows** — O(n). Single pass with backslash-counting logic.

All functions are safe for concurrent use (no shared state).

## Cross-platform

- **Quote / QuoteCommand** — produce POSIX single-quote syntax. Compatible
  with `sh`, `bash`, `zsh`, `dash`, and similar shells.
- **QuoteWindows** — produces double-quote syntax compatible with Windows
  `cmd.exe` and `CommandLineToArgv`.
