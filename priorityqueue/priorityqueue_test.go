package priorityqueue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/azghr/forge/priorityqueue"
)

func TestMinHeap(t *testing.T) {
	t.Parallel()

	t.Run("basic ordering", func(t *testing.T) {
		pq := priorityqueue.New[int]()
		pq.Push(priorityqueue.Item[int]{Value: 1, Priority: 2})
		pq.Push(priorityqueue.Item[int]{Value: 2, Priority: 1})

		it, ok := pq.Pop()
		if !ok || it.Value != 2 {
			t.Errorf("expected 2 first (priority 1), got %+v", it)
		}
		it, ok = pq.Pop()
		if !ok || it.Value != 1 {
			t.Errorf("expected 1 second, got %+v", it)
		}
		_, ok = pq.Pop()
		if ok {
			t.Error("expected empty")
		}
	})

	t.Run("FIFO for equal priority", func(t *testing.T) {
		pq := priorityqueue.New[string]()
		pq.Push(priorityqueue.Item[string]{Value: "first", Priority: 1})
		pq.Push(priorityqueue.Item[string]{Value: "second", Priority: 1})

		it, _ := pq.Pop()
		if it.Value != "first" {
			t.Errorf("expected first, got %s", it.Value)
		}
	})

	t.Run("many items ascending", func(t *testing.T) {
		n := 1000
		pq := priorityqueue.New[int]()
		for i := range n {
			pq.Push(priorityqueue.Item[int]{Value: i, Priority: i})
		}
		for i := range n {
			it, ok := pq.Pop()
			if !ok || it.Value != i {
				t.Fatalf("at %d: expected %d, got %d", i, i, it.Value)
			}
		}
	})

	t.Run("many items descending", func(t *testing.T) {
		n := 1000
		pq := priorityqueue.New[int]()
		for i := n - 1; i >= 0; i-- {
			pq.Push(priorityqueue.Item[int]{Value: i, Priority: i})
		}
		for i := range n {
			it, ok := pq.Pop()
			if !ok || it.Value != i {
				t.Fatalf("at %d: expected %d, got %d", i, i, it.Value)
			}
		}
	})

	t.Run("pop empty returns false", func(t *testing.T) {
		pq := priorityqueue.New[int]()
		_, ok := pq.Pop()
		if ok {
			t.Error("expected false from empty queue")
		}
	})

	t.Run("len", func(t *testing.T) {
		pq := priorityqueue.New[string]()
		if pq.Len() != 0 {
			t.Errorf("initial Len = %d, want 0", pq.Len())
		}
		pq.Push(priorityqueue.Item[string]{Value: "a", Priority: 1})
		if pq.Len() != 1 {
			t.Errorf("after push Len = %d, want 1", pq.Len())
		}
		pq.Pop()
		if pq.Len() != 0 {
			t.Errorf("after pop Len = %d, want 0", pq.Len())
		}
	})
}

func TestMaxHeap(t *testing.T) {
	t.Parallel()

	t.Run("max ordering", func(t *testing.T) {
		pq := priorityqueue.New[string](priorityqueue.WithMaxHeap())
		pq.Push(priorityqueue.Item[string]{Value: "low", Priority: 1})
		pq.Push(priorityqueue.Item[string]{Value: "high", Priority: 10})
		pq.Push(priorityqueue.Item[string]{Value: "mid", Priority: 5})

		it, _ := pq.Pop()
		if it.Value != "high" {
			t.Errorf("expected high first (priority 10), got %s", it.Value)
		}
		it, _ = pq.Pop()
		if it.Value != "mid" {
			t.Errorf("expected mid second (priority 5), got %s", it.Value)
		}
		it, _ = pq.Pop()
		if it.Value != "low" {
			t.Errorf("expected low last (priority 1), got %s", it.Value)
		}
	})

	t.Run("same priorities", func(t *testing.T) {
		pq := priorityqueue.New[string](priorityqueue.WithMaxHeap())
		pq.Push(priorityqueue.Item[string]{Value: "a", Priority: 1})
		pq.Push(priorityqueue.Item[string]{Value: "b", Priority: 1})

		it, _ := pq.Pop()
		if it.Value != "a" {
			t.Errorf("expected a, got %s", it.Value)
		}
	})

	t.Run("many items max", func(t *testing.T) {
		n := 1000
		pq := priorityqueue.New[int](priorityqueue.WithMaxHeap())
		for i := range n {
			pq.Push(priorityqueue.Item[int]{Value: i, Priority: i})
		}
		for i := n - 1; i >= 0; i-- {
			it, ok := pq.Pop()
			if !ok || it.Value != i {
				t.Fatalf("expected %d, got %d", i, it.Value)
			}
		}
	})
}

func TestPopContext(t *testing.T) {
	t.Parallel()

	t.Run("immediate value", func(t *testing.T) {
		ctx := context.Background()
		pq := priorityqueue.New[int]()
		pq.Push(priorityqueue.Item[int]{Value: 99, Priority: 1})
		it, ok := pq.PopContext(ctx)
		if !ok || it.Value != 99 {
			t.Errorf("expected 99, got %+v (ok=%v)", it, ok)
		}
	})

	t.Run("block until push", func(t *testing.T) {
		ctx := context.Background()
		pq := priorityqueue.New[int]()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			pq.Push(priorityqueue.Item[int]{Value: 42, Priority: 1})
		}()

		it, ok := pq.PopContext(ctx)
		if !ok || it.Value != 42 {
			t.Errorf("expected 42, got %+v (ok=%v)", it, ok)
		}
		wg.Wait()
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		pq := priorityqueue.New[int]()
		_, ok := pq.PopContext(ctx)
		if ok {
			t.Error("expected false on cancelled context")
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		pq := priorityqueue.New[int]()
		_, ok := pq.PopContext(ctx)
		if ok {
			t.Error("expected false on timeout")
		}
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	pq := priorityqueue.New[int]()

	var wg sync.WaitGroup
	n := 100

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range n {
				pq.Push(priorityqueue.Item[int]{Value: i, Priority: i})
			}
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range n {
				pq.Pop()
			}
		}()
	}

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range n {
				pq.Len()
			}
		}()
	}

	wg.Wait()
}

func BenchmarkPush(b *testing.B) {
	pq := priorityqueue.New[int]()
	for b.Loop() {
		pq.Push(priorityqueue.Item[int]{Value: 1, Priority: 1})
	}
}

func BenchmarkPushPop(b *testing.B) {
	pq := priorityqueue.New[int]()
	for b.Loop() {
		pq.Push(priorityqueue.Item[int]{Value: 1, Priority: 1})
		pq.Pop()
	}
}

func BenchmarkPushPopMany(b *testing.B) {
	pq := priorityqueue.New[int]()
	for b.Loop() {
		for range 100 {
			pq.Push(priorityqueue.Item[int]{Value: 1, Priority: 1})
		}
		for range 100 {
			pq.Pop()
		}
	}
}
