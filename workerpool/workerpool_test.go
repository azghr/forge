package workerpool_test

import (
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/azghr/forge/workerpool"
)

func TestPool(t *testing.T) {
	t.Parallel()

	t.Run("results", func(t *testing.T) {
		pool := workerpool.New[int](2)
		for i := 0; i < 4; i++ {
			val := i
			pool.Submit(func() int {
				return val * val
			})
		}
		pool.Close()

		var res []int
		for r := range pool.Results {
			res = append(res, r)
		}
		sort.Ints(res)

		want := []int{0, 1, 4, 9}
		if !reflect.DeepEqual(res, want) {
			t.Errorf("got %v, want %v", res, want)
		}
	})

	t.Run("single worker order", func(t *testing.T) {
		pool := workerpool.New[string](1)
		pool.Submit(func() string { return "a" })
		pool.Submit(func() string { return "b" })
		pool.Close()

		first := <-pool.Results
		second := <-pool.Results
		if first != "a" || second != "b" {
			t.Errorf("expected sequential a, b; got %v, %v", first, second)
		}
	})

	t.Run("no tasks", func(t *testing.T) {
		pool := workerpool.New[int](2)
		pool.Close()

		_, ok := <-pool.Results
		if ok {
			t.Error("expected closed Results channel")
		}
	})

	t.Run("single task", func(t *testing.T) {
		pool := workerpool.New[int](3)
		pool.Submit(func() int { return 42 })
		pool.Close()

		r := <-pool.Results
		if r != 42 {
			t.Errorf("got %d, want 42", r)
		}
	})

	t.Run("more tasks than workers", func(t *testing.T) {
		n := 100
		pool := workerpool.New[int](4)
		for i := 0; i < n; i++ {
			val := i
			pool.Submit(func() int { return val })
		}
		pool.Close()

		seen := make(map[int]bool)
		for r := range pool.Results {
			if seen[r] {
				t.Errorf("duplicate result %d", r)
			}
			seen[r] = true
		}
		if len(seen) != n {
			t.Errorf("got %d results, want %d", len(seen), n)
		}
	})
}

func TestPoolConcurrentSubmit(t *testing.T) {
	var wg sync.WaitGroup
	pool := workerpool.New[any](4)

	var count atomic.Int32
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				pool.Submit(func() any {
					count.Add(1)
					return nil
				})
			}
		}()
	}
	wg.Wait()
	pool.Close()

	var n int
	for range pool.Results {
		n++
	}
	if n != 200 {
		t.Errorf("expected 200 results, got %d", n)
	}
	if c := count.Load(); c != 200 {
		t.Errorf("expected 200 task executions, got %d", c)
	}
}

func BenchmarkPool(b *testing.B) {
	for b.Loop() {
		pool := workerpool.New[int](4)
		for i := 0; i < 100; i++ {
			val := i
			pool.Submit(func() int { return val })
		}
		pool.Close()
		for range pool.Results {
		}
	}
}
