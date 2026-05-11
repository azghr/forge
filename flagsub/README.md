# flagsub

Extends the standard `flag` package with subcommand support. Each subcommand has its own `flag.FlagSet` and `Run` function.

## Problem

Go's `flag` package handles flat argument parsing well, but real CLIs often need subcommands (`serve`, `deploy`, `help`). Without a library, you end up switching on `os.Args[1]` and managing separate `flag.FlagSet` instances by hand.

## Quick start

```go
import "github.com/azghr/forge/flagsub"

serve := flagsub.AddSub("serve", "Start the server", func() {
    port := serve.Flags.Int("port", 8080, "listen port")
    serve.Flags.Parse(os.Args[2:])
    fmt.Println("Serving on", *port)
})

flagsub.Parse()
```

Run with `./app serve --port=9090`.

## API

### Functions

- **`AddSub(name, desc string, run func()) *Sub`** — register a subcommand.
- **`ParseArgs(args []string) error`** — dispatch to the matching subcommand
  using `args` (e.g. `os.Args[1:]`). Returns an error if unknown.
- **`Parse()`** — parses `os.Args` and dispatches. On unknown subcommand,
  prints usage to stderr and exits with status 1.
- **`Reset()`** — clears all registered subcommands (useful in tests).

### Type `Sub`

```go
type Sub struct {
    Name        string
    Description string
    Flags       *flag.FlagSet // use to define flags for this subcommand
    Run         func()
}
```

### Error semantics

`ParseArgs` returns an error when no subcommand matches. `Parse` exits the
process via `os.Exit(1)` in that case. Subcommand `Run` functions should
call `Flags.Parse` themselves with the remaining args.

## Performance

Negligible — CLI initialisation only. `AddSub` is O(1), `Parse` is O(n) in
registered subcommands.
