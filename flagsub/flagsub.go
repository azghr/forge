// Package flagsub extends the standard flag package with subcommand support.
// Each subcommand has its own flag.FlagSet and Run function.
package flagsub

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

// Sub represents a single subcommand with its name, description, flag set,
// and run function.
type Sub struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Run         func()
}

var (
	mu     sync.Mutex
	subs   []*Sub
	exitFn = os.Exit
	errWr  io.Writer = os.Stderr
)

// SetExit replaces the exit function used by Parse. Pass nil to restore the
// default (os.Exit).
func SetExit(fn func(int)) {
	mu.Lock()
	defer mu.Unlock()
	if fn == nil {
		exitFn = os.Exit
		return
	}
	exitFn = fn
}

// SetStderr replaces the writer used for usage output. Pass nil to restore
// the default (os.Stderr).
func SetStderr(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	if w == nil {
		errWr = os.Stderr
		return
	}
	errWr = w
}

// AddSub registers a new subcommand with the given name and description.
// The run function is called when this subcommand is dispatched.
func AddSub(name, desc string, run func()) *Sub {
	sub := &Sub{
		Name:        name,
		Description: desc,
		Flags:       flag.NewFlagSet(name, flag.ContinueOnError),
		Run:         run,
	}
	mu.Lock()
	subs = append(subs, sub)
	mu.Unlock()
	return sub
}

// ParseArgs dispatches to the matching subcommand based on args (which
// should be os.Args[1:], i.e. without the program name). It returns an
// error if args is empty or the subcommand is unknown.
func ParseArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("flagsub: no subcommand provided")
	}
	name := args[0]

	mu.Lock()
	for _, sub := range subs {
		if sub.Name == name {
			mu.Unlock()
			sub.Run()
			return nil
		}
	}
	mu.Unlock()

	return fmt.Errorf("flagsub: unknown subcommand %q", name)
}

// Parse parses os.Args and dispatches to the matching subcommand. On an
// unknown subcommand it prints help to stderr and exits with status 1.
func Parse() {
	if err := ParseArgs(os.Args[1:]); err != nil {
		printUsage()
		mu.Lock()
		fn := exitFn
		mu.Unlock()
		fn(1)
	}
}

// Reset clears all registered subcommands.
func Reset() {
	mu.Lock()
	subs = nil
	exitFn = os.Exit
	errWr = os.Stderr
	mu.Unlock()
}

func printUsage() {
	mu.Lock()
	w := errWr
	list := append([]*Sub(nil), subs...)
	mu.Unlock()

	fmt.Fprintf(w, "Usage: %s <command> [flags]\n\n", os.Args[0])
	fmt.Fprintln(w, "Commands:")
	for _, sub := range list {
		fmt.Fprintf(w, "  %-12s %s\n", sub.Name, sub.Description)
	}
}
