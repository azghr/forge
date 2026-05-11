# envconfig

Load environment variables into a Go struct using struct tags — with support for required fields and default values.

## Problem

Reading configuration from environment variables in Go typically requires
boilerplate: call `os.Getenv` for each key, parse strings to the right type,
and handle missing values manually. `envconfig.Load` replaces this with a
single call and struct tags.

## Quick start

```go
import "github.com/azghr/forge/envconfig"

type Config struct {
    Port int    `env:"PORT,default=8080"`
    Host string `env:"HOST,required"`
}

os.Setenv("HOST", "localhost")

var cfg Config
if err := envconfig.Load(&cfg); err != nil {
    log.Fatal(err)
}
fmt.Println(cfg.Port, cfg.Host) // 8080 localhost
```

## API

### Functions

- **`Load(dst interface{}, opts ...Option) error`** — parse environment
  variables into the struct pointed to by `dst`.

### Options

- **`WithPrefix(prefix string)`** — prepend a prefix to every environment
  variable name (e.g. `WithPrefix("MYAPP_")` looks up `MYAPP_*`).

### Tag format

Struct tags use the key `env` with a comma-separated value:

| Tag | Behaviour |
|-----|-----------|
| `env:"VAR"` | Read from env var `VAR`. |
| `env:"VAR,required"` | Error if `VAR` is unset and has no default. |
| `env:"VAR,default=value"` | Use `value` when `VAR` is unset. |

Supported field types: `string`, `bool`, all `int`/`uint` sizes, `float32`/`float64`.

### Error semantics

- `*MissingError` — returned when a required variable is missing.
- `*ParseError` — returned when a value cannot be converted to the field
  type. Wraps the underlying conversion error (unwraps via `Unwrap()`).

## Performance

Uses reflection to scan struct fields once per `Load` call. Suitable for
init-time configuration. Time is O(fields). No allocations beyond the input
struct and temporary strings. Concurrency-safe (no global state).
