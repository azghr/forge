# Structured Error Handling in Go

Errors in production need to be clear, composable, and debuggable. Forge follows these principles:

## Return Errors, Don't Panic

Panic is for programmer mistakes, not runtime errors. Every fallible operation returns an error:

```go
// Good
val, err := DoSomething(ctx)
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// Bad
val := mustDoSomething(ctx) // panics on error
```

## Wrap Errors with Context

Always wrap errors with context about what was happening:

```go
if err := processOrder(ctx, order); err != nil {
    return fmt.Errorf("processing order %s: %w", order.ID, err)
}
```

Use `%w` (not `%v`) to preserve the error chain for `errors.Is` and `errors.As`.

## Aggregate Multiple Errors

When collecting errors from multiple operations, use `multityperror`:

```go
var errs multityperror.Error

for _, item := range items {
    if err := process(item); err != nil {
        errs.Add(err)
    }
}

return errs.Err()
```

This preserves all errors rather than returning only the first one.

## Custom Error Types

Define custom error types when callers need to distinguish error classes:

```go
type NotFoundError struct {
    Entity string
    ID     string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %s not found", e.Entity, e.ID)
}
```

Callers can use `errors.As` to handle specific error types.

## Validation Errors

Use `validator.ValidationError` for structured field-level validation:

```go
if err := validator.ValidateStruct(input); err != nil {
    var verr validator.ValidationError
    if errors.As(err, &verr) {
        for _, ferr := range verr {
            log.Error("validation failed", "field", ferr.Field, "rule", ferr.Tag)
        }
    }
}
```

## Key Rules

1. Never silence errors — `_ = doSomething()` is a code smell
2. Always wrap with context — the caller needs to know what failed and why
3. Use typed errors for actionable failures — let callers distinguish and recover
4. Log at the boundary — log errors where they cross system boundaries (HTTP handler, background worker)
