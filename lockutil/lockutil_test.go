package lockutil_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/azghr/forge/lockutil"
)

func TestTryLockMutex(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	if !lockutil.TryLockMutex(&mu) {
		t.Error("Expected lock to succeed on unlocked mutex")
	}

	if lockutil.TryLockMutex(&mu) {
		t.Error("Expected TryLock to fail when mutex is already locked")
	}
	mu.Unlock()

	if !lockutil.TryLockMutex(&mu) {
		t.Error("Expected lock to succeed after unlock")
	}
	mu.Unlock()
}

func TestTryLockRW(t *testing.T) {
	t.Parallel()

	var rw sync.RWMutex

	rw.Lock()
	if lockutil.TryLockRW(&rw) {
		t.Error("Expected TryLockRW to fail when write-locked")
	}
	rw.Unlock()

	if !lockutil.TryLockRW(&rw) {
		t.Error("Expected TryLockRW to succeed after write unlock")
	}
	rw.RUnlock()
}

func TestLockMutex(t *testing.T) {
	t.Parallel()

	t.Run("acquires lock", func(t *testing.T) {
		var mu sync.Mutex
		mu.Lock()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		got := lockutil.LockMutex(ctx, &mu)
		if got {
			t.Error("expected false when ctx already cancelled")
		}
		mu.Unlock()
	})

	t.Run("cancel while waiting", func(t *testing.T) {
		var mu sync.Mutex
		mu.Lock()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		got := lockutil.LockMutex(ctx, &mu)
		if got {
			t.Error("expected false when context times out")
		}
		mu.Unlock()
	})

	t.Run("acquires when freed", func(t *testing.T) {
		var mu sync.Mutex
		mu.Lock()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			time.Sleep(10 * time.Millisecond)
			mu.Unlock()
		}()

		got := lockutil.LockMutex(ctx, &mu)
		if !got {
			t.Error("expected true when lock becomes available")
		}
		mu.Unlock()
	})
}

func TestLockRW(t *testing.T) {
	t.Parallel()

	t.Run("acquires read lock", func(t *testing.T) {
		var rw sync.RWMutex
		rw.Lock()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		go func() {
			time.Sleep(10 * time.Millisecond)
			rw.Unlock()
		}()

		got := lockutil.LockRW(ctx, &rw)
		if !got {
			t.Error("expected true when read lock becomes available")
		}
		rw.RUnlock()
	})

	t.Run("cancel while waiting", func(t *testing.T) {
		var rw sync.RWMutex
		rw.Lock()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		got := lockutil.LockRW(ctx, &rw)
		if got {
			t.Error("expected false when context times out")
		}
		rw.Unlock()
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	var count int
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				if lockutil.TryLockMutex(&mu) {
					count++
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	if count == 0 {
		t.Error("expected at least one successful lock")
	}
}

func TestWithPollInterval(t *testing.T) {
	t.Parallel()

	var rw sync.RWMutex
	rw.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	got := lockutil.LockRW(ctx, &rw,
		lockutil.WithPollInterval(time.Microsecond),
	)
	if got {
		t.Error("expected false when write-locked with custom interval")
	}
	rw.Unlock()
}

func BenchmarkTryLockMutex(b *testing.B) {
	var mu sync.Mutex
	for b.Loop() {
		if lockutil.TryLockMutex(&mu) {
			mu.Unlock()
		}
	}
}

func BenchmarkTryLockRW(b *testing.B) {
	var rw sync.RWMutex
	for b.Loop() {
		if lockutil.TryLockRW(&rw) {
			rw.RUnlock()
		}
	}
}

func BenchmarkLockMutex(b *testing.B) {
	var mu sync.Mutex
	ctx := context.Background()
	for b.Loop() {
		if lockutil.LockMutex(ctx, &mu) {
			mu.Unlock()
		}
	}
}
