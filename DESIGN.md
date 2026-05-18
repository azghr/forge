# Forge Design Philosophy

Forge is a collection of Go utility packages built around a consistent design language. This document defines that language.

## Principles

### Small Packages Rule

Each package solves exactly 1–3 related problems. If a package does more than that, split it.

**Signal:** A package's README fits in 3 paragraphs.

### Stdlib-First Design

Prefer the standard library over external dependencies. Every external dependency is a liability.

- Use `net/http` over frameworks
- Use `encoding/json` over JSON libraries
- Use `flag` over cobra (unless subcommands demand `flagsub`)
- Use `container/heap`, `sync`, `context` instead of re-inventing

### Zero Magic

No init functions. No global state. No reflection unless absolutely necessary.

What the code does should be obvious from reading it. A developer should be able to understand a package by reading its public API without guessing.

### Explicit APIs

Favor explicit over convenient. Return errors. Don't swallow them. Don't panic.

```go
// Good
val, err := DoSomething(ctx, opts)
if err != nil {
    return err
}

// Bad
val := DoSomething(opts) // panics on error
```

### Composability

Packages should compose like UNIX pipes. Each package is a building block, not a framework.

- `retry` wraps any function — it doesn't care what the function does
- `workerpool` takes any task — it doesn't know about queues
- `validator` validates structs — it doesn't know about HTTP

### Concurrency Safety by Default

All exported types are safe for concurrent use unless documented otherwise. Use `sync.RWMutex`, atomic operations, or channel-based designs.

### Generics, Not Interfaces

When a package can work with any type, use Go generics (`[T any]`) rather than `interface{}` or dynamic dispatch. This preserves type safety at zero runtime cost.

## Package Template

```
pkgname/
├── pkgname.go       # Core types and functions
├── options.go       # Functional options (if needed)
├── errors.go        # Custom error types (if needed)
└── *_test.go        # Tests + examples
```

## Checklist for Every Package

- [ ] Minimal exported API (fewer symbols is better)
- [ ] No panics in normal use
- [ ] Context-aware where blocking occurs
- [ ] Zero global state
- [ ] Thread-safe
- [ ] Sensible defaults (New() works out of the box)
- [ ] Clear, descriptive error messages
- [ ] Table-driven tests
- [ ] Runnable examples in `ExampleXxx` functions
- [ ] Benchmarks for hot paths
- [ ] No dependencies beyond stdlib (exceptions documented)

## Versioning

All packages share the monorepo version. We follow semantic versioning. Breaking changes require a major version bump and a migration guide.
