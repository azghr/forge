package shellquote_test

import (
	"strings"
	"testing"

	"github.com/azghr/forge/shellquote"
)

func TestQuote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input string
		want string
	}{
		{name: "empty", input: "", want: "''"},
		{name: "plain", input: "hello", want: "'hello'"},
		{name: "with spaces", input: "hello world", want: "'hello world'"},
		{name: "single quote", input: "a'b c", want: "'a'\"'\"'b c'"},
		{name: "multiple quotes", input: `it's "fine"`, want: `'it'"'"'s "fine"'`},
		{name: "special chars", input: "file; rm -rf /", want: "'file; rm -rf /'"},
		{name: "newline", input: "a\nb", want: "'a\nb'"},
		{name: "only quotes", input: "'''", want: "''\"'\"''\"'\"''\"'\"''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellquote.Quote(tt.input)
			if got != tt.want {
				t.Errorf("Quote(%q) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestQuoteCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "empty", args: nil, want: ""},
		{name: "single", args: []string{"echo"}, want: "'echo'"},
		{name: "multiple", args: []string{"echo", "hello world", "foo"}, want: "'echo' 'hello world' 'foo'"},
		{name: "with quotes", args: []string{"echo", "it's"}, want: "'echo' 'it'\"'\"'s'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellquote.QuoteCommand(tt.args)
			if got != tt.want {
				t.Errorf("QuoteCommand(%v) = %s, want %s", tt.args, got, tt.want)
			}
		})
	}
}

func TestQuoteCommandQuotesEachArg(t *testing.T) {
	args := []string{"echo", "hello world", "foo"}
	cmd := shellquote.QuoteCommand(args)
	if !strings.Contains(cmd, "'hello world'") {
		t.Errorf("QuoteCommand should quote args with spaces: %s", cmd)
	}
}

func TestQuoteWindows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input string
		want string
	}{
		{name: "empty", input: "", want: `""`},
		{name: "plain", input: "hello", want: `"hello"`},
		{name: "with spaces", input: "hello world", want: `"hello world"`},
		{name: "with backslash", input: `C:\path`, want: `"C:\path"`},
		{name: "trailing backslash", input: `C:\path\`, want: `"C:\path\\"`},
		{name: "with double quote", input: `a"b`, want: `"a\"b"`},
		{name: "backslash before quote", input: `a\"b`, want: `"a\\\"b"`},
		{name: "multiple backslashes before quote", input: `a\\"b`, want: `"a\\\\\"b"`},
		{name: "only backslashes", input: `\\`, want: `"\\\\"`},
		{name: "mixed", input: `C:\"path with spaces"\`, want: `"C:\\\"path with spaces\"\\"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellquote.QuoteWindows(tt.input)
			if got != tt.want {
				t.Errorf("QuoteWindows(%q) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkQuote(b *testing.B) {
	for b.Loop() {
		shellquote.Quote("a'b c")
	}
}

func BenchmarkQuotePlain(b *testing.B) {
	for b.Loop() {
		shellquote.Quote("hello")
	}
}

func BenchmarkQuoteCommand(b *testing.B) {
	args := []string{"ls", "-l", "my file", "it's"}
	for b.Loop() {
		shellquote.QuoteCommand(args)
	}
}

func BenchmarkQuoteWindows(b *testing.B) {
	for b.Loop() {
		shellquote.QuoteWindows(`C:\path with spaces\`)
	}
}

func BenchmarkQuoteWindowsPlain(b *testing.B) {
	for b.Loop() {
		shellquote.QuoteWindows("hello")
	}
}
