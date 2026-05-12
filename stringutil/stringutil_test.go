package stringutil_test

import (
	"testing"

	"github.com/azghr/forge/stringutil"
)

func TestTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "hello world", in: "hello world", want: "Hello World"},
		{name: "single word", in: "hello", want: "Hello"},
		{name: "empty", in: "", want: ""},
		{name: "already title", in: "Hello World", want: "Hello World"},
		{name: "numbers", in: "hello 123 world", want: "Hello 123 World"},
		{name: "mixed case", in: "hELLO wORLD", want: "HELLO WORLD"},
		{name: "with punctuation", in: "hello, world!", want: "Hello, World!"},
		{name: "unicode", in: "héllo wörld", want: "Héllo Wörld"},
		{name: "multiple spaces", in: "hello    world", want: "Hello    World"},
		{name: "leading space", in: " hello world", want: " Hello World"},
		{name: "trailing space", in: "hello world ", want: "Hello World "},
		{name: "single letter", in: "a", want: "A"},
		{name: "non-letters only", in: "123 !@#", want: "123 !@#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringutil.Title(tt.in)
			if got != tt.want {
				t.Errorf("Title(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSlug(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			in   string
			want string
		}{
			{name: "basic", in: "Hello World", want: "hello-world"},
			{name: "already slug", in: "hello-world", want: "hello-world"},
			{name: "leading/trailing spaces", in: "  hello world  ", want: "hello-world"},
			{name: "multiple spaces", in: "hello   world", want: "hello-world"},
			{name: "empty", in: "", want: ""},
			{name: "only special chars", in: "@#$%", want: ""},
			{name: "with punctuation", in: "Hello, Go!", want: "hello-go"},
			{name: "underscores", in: "hello_world", want: "hello-world"},
			{name: "mixed case", in: "Go Lang Library", want: "go-lang-library"},
			{name: "unicode preserved", in: "héllo wörld", want: "héllo-wörld"},
			{name: "numbers", in: "version 2.0", want: "version-2-0"},
			{name: "single word", in: "Hello", want: "hello"},
			{name: "dashes and dots", in: "a--b..c", want: "a-b-c"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := stringutil.Slug(tt.in)
				if got != tt.want {
					t.Errorf("Slug(%q) = %q, want %q", tt.in, got, tt.want)
				}
			})
		}
	})

	t.Run("custom separator", func(t *testing.T) {
		t.Parallel()

		got := stringutil.Slug("Hello World", stringutil.WithSeparator("_"))
		if got != "hello_world" {
			t.Errorf("Slug with underscore = %q, want %q", got, "hello_world")
		}

		got2 := stringutil.Slug("Hello World", stringutil.WithSeparator(""))
		if got2 != "hello-world" {
			t.Errorf("Slug with empty separator = %q, want %q", got2, "hello-world")
		}
	})

	t.Run("max length", func(t *testing.T) {
		t.Parallel()

		got := stringutil.Slug("hello world foo bar", stringutil.WithMaxLength(11))
		if got != "hello-world" {
			t.Errorf("Slug with max length = %q, want %q", got, "hello-world")
		}

		got2 := stringutil.Slug("hello", stringutil.WithMaxLength(10))
		if got2 != "hello" {
			t.Errorf("Slug with max > length = %q, want %q", got2, "hello")
		}
	})

	t.Run("combined options", func(t *testing.T) {
		t.Parallel()

		got := stringutil.Slug("Hello World Foo",
			stringutil.WithSeparator("_"),
			stringutil.WithMaxLength(11),
		)
		if got != "hello_world" {
			t.Errorf("Slug with combined options = %q, want %q", got, "hello_world")
		}
	})

	t.Run("max length zero", func(t *testing.T) {
		t.Parallel()

		got := stringutil.Slug("hello world foo bar", stringutil.WithMaxLength(0))
		if got != "hello-world-foo-bar" {
			t.Errorf("Slug with max length 0 = %q, want %q", got, "hello-world-foo-bar")
		}
	})
}

func TestRemoveAccents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "café", in: "café", want: "cafe"},
		{name: "naïve", in: "naïve", want: "naive"},
		{name: "no accents", in: "hello", want: "hello"},
		{name: "empty", in: "", want: ""},
		{name: "crème brûlée", in: "crème brûlée", want: "creme brulee"},
		{name: "piñata", in: "piñata", want: "pinata"},
		{name: "über cool", in: "über cool", want: "uber cool"},
		{name: "résumé", in: "résumé", want: "resume"},
		{name: "élève", in: "élève", want: "eleve"},
		{name: "garçon", in: "garçon", want: "garcon"},
		{name: "Ångström", in: "Ångström", want: "Angstrom"},
		{name: "Æsop", in: "Æsop", want: "Asop"},
		{name: "mañana", in: "mañana", want: "manana"},
		{name: "déjà vu", in: "déjà vu", want: "deja vu"},
		{name: "mixed", in: "Héllò Wörld", want: "Hello World"},
		{name: "Cyrillic preserved", in: "Привет", want: "Привет"},
		{name: "Chinese preserved", in: "你好", want: "你好"},
		{name: "soft hyphen", in: "soft\u00ADhyphen", want: "soft hyphen"},
		{name: "no-break space", in: "no\u00A0break", want: "no break"},
		{name: "sharp s", in: "straße", want: "strase"},
		{name: "combined", in: "façade", want: "facade"},
		{name: "Đông", in: "Đông", want: "Dong"},
		{name: "Łódź", in: "Łódź", want: "Lodz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringutil.RemoveAccents(tt.in)
			if got != tt.want {
				t.Errorf("RemoveAccents(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	run := make(chan struct{})
	done := make(chan struct{}, 30)

	for range 30 {
		go func() {
			<-run
			stringutil.Title("hello world")
			stringutil.Slug("Hello World")
			stringutil.Slug("Hello World", stringutil.WithSeparator("_"), stringutil.WithMaxLength(10))
			stringutil.RemoveAccents("café")
			done <- struct{}{}
		}()
	}

	close(run)

	for range 30 {
		<-done
	}
}

func BenchmarkTitle(b *testing.B) {
	for b.Loop() {
		stringutil.Title("hello world, this is a test")
	}
}

func BenchmarkSlug(b *testing.B) {
	for b.Loop() {
		stringutil.Slug("Hello World, This is a Test!")
	}
}

func BenchmarkSlugWithOptions(b *testing.B) {
	for b.Loop() {
		stringutil.Slug("Hello World, This is a Test!",
			stringutil.WithSeparator("_"),
			stringutil.WithMaxLength(20),
		)
	}
}

func BenchmarkRemoveAccents(b *testing.B) {
	for b.Loop() {
		stringutil.RemoveAccents("crème brûlée naïve café")
	}
}

func BenchmarkRemoveAccentsASCII(b *testing.B) {
	for b.Loop() {
		stringutil.RemoveAccents("hello world")
	}
}
