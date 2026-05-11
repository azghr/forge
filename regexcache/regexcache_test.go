package regexcache_test

import (
	"sync"
	"testing"

	"github.com/azghr/forge/regexcache"
)

func TestCompileCache(t *testing.T) {
	t.Parallel()

	cache := regexcache.New()

	r1 := cache.MustCompile("a.*z")
	r2 := cache.MustCompile("a.*z")
	if r1 != r2 {
		t.Error("Expected same pointer from cache")
	}

	if _, err := cache.Compile("("); err == nil {
		t.Error("Expected error on invalid regex")
	}
}

func TestCompile(t *testing.T) {
	t.Parallel()

	cache := regexcache.New()

	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{name: "simple", pattern: "a.*z", wantErr: false},
		{name: "anchored", pattern: "^[a-z]+$", wantErr: false},
		{name: "named groups", pattern: `(?P<name>\w+)`, wantErr: false},
		{name: "invalid parens", pattern: "(", wantErr: true},
		{name: "invalid bracket", pattern: "[a-z", wantErr: true},
		{name: "empty", pattern: "", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := cache.Compile(tt.pattern)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Compile(%q) expected error, got %v", tt.pattern, r)
				}
				return
			}
			if err != nil {
				t.Errorf("Compile(%q) unexpected error: %v", tt.pattern, err)
			}
			if r == nil {
				t.Error("expected non-nil *Regexp")
			}
		})
	}
}

func TestMustCompile(t *testing.T) {
	t.Parallel()

	cache := regexcache.New()

	t.Run("valid", func(t *testing.T) {
		r := cache.MustCompile("^a.*z$")
		if !r.MatchString("abz") {
			t.Error("expected match")
		}
	})

	t.Run("panics on invalid", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("expected panic")
			}
		}()
		cache.MustCompile("(")
	})
}

func TestCompileReturnsCached(t *testing.T) {
	t.Parallel()

	cache := regexcache.New()
	r1, _ := cache.Compile("^[a-z]+$")
	r2, _ := cache.Compile("^[a-z]+$")
	if r1 != r2 {
		t.Error("Compile should return cached value")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := regexcache.New()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			pattern := "a"
			if n%2 == 0 {
				pattern = "[a-z]+"
			}
			r, err := cache.Compile(pattern)
			if err != nil {
				t.Errorf("Compile(%q) error: %v", pattern, err)
				return
			}
			if r == nil {
				t.Error("expected non-nil result")
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentMustCompile(t *testing.T) {
	cache := regexcache.New()
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.MustCompile("\\w+")
		}()
	}

	wg.Wait()
}

func TestWithMaxSize(t *testing.T) {
	t.Parallel()

	cache := regexcache.New(regexcache.WithMaxSize(2))

	cache.MustCompile("a")
	cache.MustCompile("b")
	cache.MustCompile("c")

	got := cache.MustCompile("c")
	if !got.MatchString("c") {
		t.Error("expected c to be cached")
	}

	got = cache.MustCompile("b")
	if !got.MatchString("b") {
		t.Error("expected b to be cached")
	}
}

func BenchmarkCompile(b *testing.B) {
	cache := regexcache.New()
	for b.Loop() {
		cache.MustCompile("^[a-z]+$")
	}
}

func BenchmarkCompileCacheHit(b *testing.B) {
	cache := regexcache.New()
	cache.MustCompile("^[a-z]+$")
	b.ResetTimer()
	for b.Loop() {
		cache.MustCompile("^[a-z]+$")
	}
}

func BenchmarkCompileUniquePatterns(b *testing.B) {
	cache := regexcache.New()
	for b.Loop() {
		cache.MustCompile("^[a-z]+\\d+$")
	}
}
