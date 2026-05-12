# validator

Struct field validation via tags. Supports `nonzero` and `email` rules.

## Problem

Go applications frequently need to validate user input, config values, or API
payloads before processing them. The standard library has no built-in struct
validation — every project ends up writing boilerplate `if` checks. This
package provides a small, tag-based validator with zero external dependencies.

## Quick start

```go
import "github.com/azghr/forge/validator"

type User struct {
    Name  string `validate:"nonzero"`
    Email string `validate:"email"`
}

u := User{Name: "Alice", Email: "x@x"}
if err := validator.ValidateStruct(u); err != nil {
    log.Fatal(err)
}
```

## API

### Functions

- **`ValidateStruct(v interface{}, opts ...Option) error`** — validate all
  exported fields of a struct by their `validate` tags. Returns a
  `ValidationError` on failure, or `nil` if all fields pass.

### Options

- **`WithTagName(name string) Option`** — set the struct tag key (default
  `"validate"`).

### Supported tags

| Tag       | Applies to             | Description                                   |
|-----------|------------------------|-----------------------------------------------|
| `nonzero` | string, int, bool, ptr, slice, map, etc. | Value must not be the zero value for its type |
| `email`   | string                 | Value must match a basic email pattern        |

### Error types

- **`ValidationError`** — `[]FieldError` that implements `error`. Each
  `FieldError` contains `Field` (name), `Tag` (failed rule), and `Value`.

## Performance

Uses reflection once per call. Benchmark on Apple M1 Max:
- Valid struct: ~500 ns/op
- Invalid struct: ~600 ns/op
