package flagsub_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/azghr/forge/flagsub"
)

func TestSubcommands(t *testing.T) {
	t.Parallel()

	t.Run("dispatch to known subcommand", func(t *testing.T) {
		flagsub.Reset()

		var result int
		var serve *flagsub.Sub
		serve = flagsub.AddSub("serve", "Start server", func() {
			p := serve.Flags.Int("port", 8080, "listen port")
			serve.Flags.Parse(os.Args[2:])
			result = *p
		})

		oldArgs := os.Args
		os.Args = []string{"myapp", serve.Name, "--port=1234"}
		defer func() { os.Args = oldArgs }()

		flagsub.Parse()

		if result != 1234 {
			t.Errorf("got port %d, want 1234", result)
		}
	})

	t.Run("unknown subcommand returns error via ParseArgs", func(t *testing.T) {
		flagsub.Reset()
		flagsub.AddSub("foo", "Foo command", func() {})

		err := flagsub.ParseArgs([]string{"bogus"})
		if err == nil {
			t.Fatal("expected error for unknown subcommand")
		}
	})

	t.Run("empty args returns error", func(t *testing.T) {
		flagsub.Reset()

		err := flagsub.ParseArgs(nil)
		if err == nil {
			t.Fatal("expected error for empty args")
		}
	})

	t.Run("Parse prints usage and exits on unknown", func(t *testing.T) {
		flagsub.Reset()
		flagsub.AddSub("foo", "Foo command", func() {})

		var exited int
		var buf bytes.Buffer

		flagsub.SetExit(func(code int) { exited = code; panic("exit") })
		flagsub.SetStderr(&buf)

		oldArgs := os.Args
		os.Args = []string{"myapp", "bogus"}
		defer func() { os.Args = oldArgs }()

		func() {
			defer func() { recover() }()
			flagsub.Parse()
		}()

		if exited != 1 {
			t.Errorf("expected exit(1), got exit(%d)", exited)
		}
		if !bytes.Contains(buf.Bytes(), []byte("Usage:")) {
			t.Errorf("expected usage output, got %q", buf.String())
		}

		flagsub.SetExit(nil)
		flagsub.SetStderr(nil)
	})

	t.Run("multiple subcommands", func(t *testing.T) {
		flagsub.Reset()

		var out []string
		flagsub.AddSub("a", "Command A", func() { out = append(out, "a") })
		flagsub.AddSub("b", "Command B", func() { out = append(out, "b") })

		if err := flagsub.ParseArgs([]string{"b"}); err != nil {
			t.Fatal(err)
		}
		if len(out) != 1 || out[0] != "b" {
			t.Errorf("expected [b], got %v", out)
		}
	})
}

func TestReset(t *testing.T) {
	flagsub.Reset()
	flagsub.AddSub("x", "X", func() {})

	if err := flagsub.ParseArgs([]string{"x"}); err != nil {
		t.Fatal("expected subcommand to exist")
	}

	flagsub.Reset()
	if err := flagsub.ParseArgs([]string{"x"}); err == nil {
		t.Fatal("expected error after reset")
	}
}

func TestSub(t *testing.T) {
	t.Parallel()

	t.Run("subcommand fields are populated", func(t *testing.T) {
		flagsub.Reset()

		sub := flagsub.AddSub("greet", "Print a greeting", func() {})

		if sub.Name != "greet" {
			t.Errorf("got name %q, want %q", sub.Name, "greet")
		}
		if sub.Description != "Print a greeting" {
			t.Errorf("got desc %q, want %q", sub.Description, "Print a greeting")
		}
		if sub.Flags == nil {
			t.Error("Flags should not be nil")
		}
		if sub.Run == nil {
			t.Error("Run should not be nil")
		}
	})
}
