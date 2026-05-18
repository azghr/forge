# Forge Philosophy

Forge is a different kind of Go ecosystem. It doesn't compete with frameworks. It competes with complexity.

## Small Packages Rule

Each package solves exactly 1-3 related problems and no more. If a package description needs more than three paragraphs, the package is too big.

This is UNIX philosophy applied to Go: do one thing well, and compose freely. A developer should understand what a package does from reading its public API — no guesswork, no surprises.

> A package's README fits in 3 paragraphs. If it doesn't, split the package.

## Stdlib-First Design

The Go standard library is the best dependency you'll never have to update. We prefer it over everything else:

- `net/http` over frameworks
- `encoding/json` over JSON libraries
- `flag` over cobra (unless subcommands demand `flagsub`)
- `log/slog` over logging frameworks
- `container/heap`, `sync`, `context` over re-inventing wheels

Every external dependency is a liability. It must be maintained, audited, and justified. In Forge, the default answer to "should we add a dependency?" is no.

This doesn't mean we never add dependencies. It means each one must earn its place.

## Zero Magic

Forge packages have no init functions, no global state, no reflection unless it's strictly necessary, and no "it just works" behavior that requires reading the source to understand.

What the code does should be obvious from reading it. When you call `retry.RetryContext(ctx, config, fn)`, you know exactly what happens: retry with exponential backoff. No surprises. No hidden goroutines. No magic.

The opposite of magic is clarity. We choose clarity.

## Explicit Over Implicit

Good APIs are boring. They return errors. They don't panic. They don't swallow failures.

```go
// Good: explicit, handles errors
val, err := DoSomething(ctx, opts)
if err != nil {
    return err
}

// Bad: implicit, panics on failure
val := DoSomething(opts)
```

Implicit behavior is the enemy of production systems. When something goes wrong at 3 AM, you want the code to tell you exactly what happened — not silently recover, not panic, not log and continue.

## Composability

Packages compose like UNIX pipes. Each package is a building block, not a framework:

- `retry` wraps any function — it doesn't care what the function does
- `workerpool` takes any task — it doesn't know about queues or databases
- `validator` validates structs — it doesn't know about HTTP or forms
- `envconfig` loads config into structs — it doesn't know about your application

This means you use Forge packages together, not as a platform. You're always in control of the architecture. Forge provides the tools; you provide the design.

## Concurrency Safety by Default

All exported types are safe for concurrent use unless explicitly documented otherwise. We use `sync.RWMutex`, atomic operations, and channel-based designs — whatever fits the problem.

Concurrency bugs are the hardest to find and fix. Forge prevents them at the API level.

## Generics, Not Interfaces

When a package works with any type, we use Go generics (`[T any]`) rather than `interface{}` or dynamic dispatch. This preserves type safety at zero runtime cost.

```go
// Good: generic, type-safe, zero allocation
func Map[T, U any](input []T, fn func(T) U) []U

// Bad: interface{}, no type safety
func Map(input interface{}, fn interface{}) interface{}
```

Generics make Forge packages both safe and fast. You get compile-time type checking without sacrificing performance.

## Why This Matters

Most Go libraries fall into one of two traps:

1. **Too small** — solve one problem but don't compose well
2. **Too big** — become frameworks that dictate your architecture

Forge avoids both. Each package is independently useful but designed to compose. You can use one package or all of them — the experience is the same: clear APIs, no magic, production-safe defaults.

We believe the best Go code reads like the standard library. Forge is our attempt to build a modern standard library extension that follows the same principles.

---

*"Perfection is achieved not when there is nothing more to add, but when there is nothing left to take away." — Antoine de Saint-Exupéry*
