# forge

Monorepo for Go utility and helper packages.

## CLI

`forge` is a CLI tool for working with Forge packages.

```
forge doctor               Check Go version, config validity, env issues
forge example <type>       Generate working example scaffold
forge new <type> <name>    Scaffold production-ready project
```

### Install

```sh
go install github.com/azghr/forge/cmd/forge@latest
```

Or run from source:

```sh
go run ./cmd/forge <command>
```

## Examples

Ready-to-run examples showing Forge packages working together:

- [cli](examples/cli/) — Task manager CLI using `flagsub`, `envconfig`, `tablewriter`, `sliceutil`, `shellquote`, `option`, `regexcache`
- [server](examples/server/) — Task API server using `flagsub`, `envconfig`, `validator`, `pathsafe`, `stringutil`, `atomicfile`, `stopwatch`
- [config](examples/config/) — JSON config merger using `jsonmerge`, `atomicfile`, `pathsafe`
- [pipeline](examples/pipeline/) — Data processing pipeline using `orderedset`, `priorityqueue`, `queue`, `workerpool`, `cache`, `retry`, `lockutil`, `multityperror`, `mathutil`

Run any example:

```sh
cd examples/<name>
go run .
```

## Packages

- [atomicfile](atomicfile/) — Atomic file writes without leaving partial data on failure.
- [cache](cache/) — Generic in-memory TTL cache.
- [envconfig](envconfig/) — Load environment variables into structs via tags.
- [flagsub](flagsub/) — Subcommand support for the standard `flag` package.
- [jsonmerge](jsonmerge/) — Recursively merge and diff JSON-like data.
- [lockutil](lockutil/) — Non-blocking `TryLock` and context-aware `Lock` for Mutex/RWMutex.
- [mathutil](mathutil/) — Small math helpers: Clamp, Sign, Lerp, GCD, ApproxEqual.
- [multityperror](multityperror/) — Aggregate multiple errors into one.
- [option](option/) — Generic Option (Maybe) type for Go.
- [orderedset](orderedset/) — Insertion-ordered set with Union, Intersect operations.
- [pathsafe](pathsafe/) — Safe path joining to prevent directory traversal.
- [priorityqueue](priorityqueue/) — Generic binary heap (min/max) with concurrency-safe push/pop.
- [queue](queue/) — Generic FIFO queue (ring-buffer, concurrency-safe, blocking dequeue).
- [regexcache](regexcache/) — Concurrency-safe cache for compiled regex patterns.
- [retry](retry/) — Retry operations with exponential backoff and full-jitter.
- [shellquote](shellquote/) — Shell-safe string quoting for POSIX and Windows.
- [sliceutil](sliceutil/) — Generic slice operations: Map, Filter, Reduce, All, Any, Chunk.
- [stopwatch](stopwatch/) — Simple stopwatch for benchmarking code blocks.
- [stringutil](stringutil/) — String transformations: Title, Slug, RemoveAccents.
- [tablewriter](tablewriter/) — Format tabular data as ASCII tables.
- [validator](validator/) — Struct field validation via tags (nonzero, email).
- [workerpool](workerpool/) — Fixed-size worker goroutine pool.
