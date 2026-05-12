package cache_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/azghr/forge/cache"
)

func TestCacheSetGet(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	c.Set("a", 1)
	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Errorf("expected (1, true), got (%d, %v)", v, ok)
	}
}

func TestCacheGetMissing(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	_, ok := c.Get("nope")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestCacheExpire(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](10 * time.Millisecond)

	c.Set("a", 1)
	time.Sleep(15 * time.Millisecond)

	if _, ok := c.Get("a"); ok {
		t.Error("expected expired entry to be removed")
	}
}

func TestCacheNoExpire(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, string](0)

	c.Set("a", "val")
	v, ok := c.Get("a")
	if !ok || v != "val" {
		t.Errorf("no-expire cache: expected (val, true), got (%q, %v)", v, ok)
	}
}

func TestCacheOverwrite(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](time.Hour)

	c.Set("k", 1)
	c.Set("k", 2)

	v, ok := c.Get("k")
	if !ok || v != 2 {
		t.Errorf("expected (2, true), got (%d, %v)", v, ok)
	}
}

func TestCacheDelete(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	c.Set("k", 1)
	c.Delete("k")

	if _, ok := c.Get("k"); ok {
		t.Error("expected key to be deleted")
	}
}

func TestCacheDeleteMissing(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	c.Delete("nope")
}

func TestCacheLen(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	if n := c.Len(); n != 0 {
		t.Errorf("expected 0, got %d", n)
	}

	c.Set("a", 1)
	c.Set("b", 2)

	if n := c.Len(); n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestCacheExpiredEntryRemovedOnGet(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](10 * time.Millisecond)

	c.Set("a", 1)
	time.Sleep(15 * time.Millisecond)

	c.Get("a")

	if n := c.Len(); n != 0 {
		t.Errorf("expected 0 after lazy cleanup, got %d", n)
	}
}

func TestCacheGetOrLoad(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	loader := func(ctx context.Context) (int, error) {
		return 42, nil
	}

	v, err := c.GetOrLoad(context.Background(), "k", loader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}

	v, ok := c.Get("k")
	if !ok || v != 42 {
		t.Errorf("expected cached (42, true), got (%d, %v)", v, ok)
	}
}

func TestCacheGetOrLoadError(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](0)

	loader := func(ctx context.Context) (int, error) {
		return 0, fmt.Errorf("load failed")
	}

	_, err := c.GetOrLoad(context.Background(), "k", loader)
	if err == nil {
		t.Fatal("expected error from loader")
	}

	if _, ok := c.Get("k"); ok {
		t.Error("expected key NOT to be cached on loader error")
	}
}

func TestCacheGetOrLoadCached(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](time.Hour)

	c.Set("k", 99)

	loader := func(ctx context.Context) (int, error) {
		return 42, nil
	}

	v, err := c.GetOrLoad(context.Background(), "k", loader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 99 {
		t.Errorf("expected cached 99, got %d", v)
	}
}

func TestCacheGetOrLoadExpired(t *testing.T) {
	t.Parallel()
	c := cache.NewCache[string, int](10 * time.Millisecond)

	c.Set("k", 1)
	time.Sleep(15 * time.Millisecond)

	loader := func(ctx context.Context) (int, error) {
		return 2, nil
	}

	v, err := c.GetOrLoad(context.Background(), "k", loader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 2 {
		t.Errorf("expected freshly loaded 2, got %d", v)
	}
}

func TestCacheConcurrentSetGet(t *testing.T) {
	c := cache.NewCache[int, int](0)

	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Set(n, n)
		}(i)
	}
	wg.Wait()

	for i := range 20 {
		v, ok := c.Get(i)
		if !ok || v != i {
			t.Errorf("concurrent set/get: expected (%d, true), got (%d, %v)", i, v, ok)
		}
	}
}

func TestCacheConcurrentGetOrLoad(t *testing.T) {
	c := cache.NewCache[int, int](0)

	var wg sync.WaitGroup
	var loadCount atomic.Int64

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			loader := func(ctx context.Context) (int, error) {
				loadCount.Add(1)
				return 42, nil
			}
			c.GetOrLoad(context.Background(), 1, loader)
		}()
	}
	wg.Wait()
}

func TestCacheConcurrentReadWrite(t *testing.T) {
	c := cache.NewCache[int, int](0)

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				c.Set(i, i)
			}
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				c.Get(i)
			}
		}()
	}

	wg.Wait()
}

func TestCacheWithCleanupInterval(t *testing.T) {
	c := cache.NewCache[string, int](
		10*time.Millisecond,
		cache.WithCleanupInterval(5*time.Millisecond),
	)
	defer c.Stop()

	c.Set("a", 1)
	time.Sleep(30 * time.Millisecond)

	if n := c.Len(); n != 0 {
		t.Errorf("expected cleanup goroutine to remove expired entries, got %d", n)
	}
}

func TestCacheStopCleanup(t *testing.T) {
	c := cache.NewCache[string, int](
		10*time.Millisecond,
		cache.WithCleanupInterval(5*time.Millisecond),
	)

	c.Set("a", 1)
	c.Stop()

	// After Stop, the goroutine is done. The entry should eventually expire
	// via lazy cleanup, but won't be cleaned up by the background goroutine.
	time.Sleep(30 * time.Millisecond)

	_, ok := c.Get("a")
	if ok {
		t.Error("expected entry to be expired")
	}
}

func TestCacheTableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ttl      time.Duration
		ops      func(c *cache.Cache[string, int])
		wantGet  map[string]int
		wantMiss []string
	}{
		{
			name: "set and get",
			ttl:  0,
			ops: func(c *cache.Cache[string, int]) {
				c.Set("x", 10)
			},
			wantGet:  map[string]int{"x": 10},
			wantMiss: nil,
		},
		{
			name: "overwrite",
			ttl:  0,
			ops: func(c *cache.Cache[string, int]) {
				c.Set("x", 10)
				c.Set("x", 20)
			},
			wantGet:  map[string]int{"x": 20},
			wantMiss: nil,
		},
		{
			name: "delete removes",
			ttl:  0,
			ops: func(c *cache.Cache[string, int]) {
				c.Set("x", 10)
				c.Delete("x")
			},
			wantGet:  nil,
			wantMiss: []string{"x"},
		},
		{
			name: "expired entry",
			ttl:  10 * time.Millisecond,
			ops: func(c *cache.Cache[string, int]) {
				c.Set("x", 10)
				time.Sleep(15 * time.Millisecond)
			},
			wantGet:  nil,
			wantMiss: []string{"x"},
		},
		{
			name: "no expiry when ttl zero",
			ttl:  0,
			ops: func(c *cache.Cache[string, int]) {
				c.Set("x", 10)
			},
			wantGet:  map[string]int{"x": 10},
			wantMiss: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cache.NewCache[string, int](tt.ttl)
			tt.ops(c)
			for k, wantV := range tt.wantGet {
				v, ok := c.Get(k)
				if !ok {
					t.Errorf("key %q: expected ok=true", k)
				}
				if v != wantV {
					t.Errorf("key %q: expected %d, got %d", k, wantV, v)
				}
			}
			for _, k := range tt.wantMiss {
				if _, ok := c.Get(k); ok {
					t.Errorf("key %q: expected ok=false", k)
				}
			}
		})
	}
}

func TestCacheDifferentTypes(t *testing.T) {
	t.Parallel()

	t.Run("int keys string values", func(t *testing.T) {
		c := cache.NewCache[int, string](0)
		c.Set(1, "one")
		v, ok := c.Get(1)
		if !ok || v != "one" {
			t.Errorf("expected (one, true), got (%q, %v)", v, ok)
		}
	})

	t.Run("struct keys", func(t *testing.T) {
		type key struct{ a, b int }
		c := cache.NewCache[key, string](0)
		k := key{a: 1, b: 2}
		c.Set(k, "val")
		v, ok := c.Get(k)
		if !ok || v != "val" {
			t.Errorf("expected (val, true), got (%q, %v)", v, ok)
		}
	})
}

func TestCacheStopIdempotent(t *testing.T) {
	c := cache.NewCache[string, int](
		time.Millisecond,
		cache.WithCleanupInterval(time.Millisecond),
	)
	c.Stop()
	c.Stop()
}

func BenchmarkCacheSet(b *testing.B) {
	c := cache.NewCache[int, int](0)
	b.ReportAllocs()
	for b.Loop() {
		c.Set(1, 1)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := cache.NewCache[int, int](0)
	c.Set(1, 1)
	b.ReportAllocs()
	for b.Loop() {
		c.Get(1)
	}
}

func BenchmarkCacheGetMiss(b *testing.B) {
	c := cache.NewCache[int, int](0)
	b.ReportAllocs()
	for b.Loop() {
		c.Get(999)
	}
}

func BenchmarkCacheGetOrLoad(b *testing.B) {
	c := cache.NewCache[int, int](0)
	ctx := context.Background()
	loader := func(ctx context.Context) (int, error) { return 42, nil }
	b.ReportAllocs()
	for b.Loop() {
		c.GetOrLoad(ctx, 1, loader)
	}
}
