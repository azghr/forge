package queue_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/azghr/forge/queue"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	t.Run("basic FIFO", func(t *testing.T) {
		q := queue.New[string]()
		q.Enqueue("a")
		q.Enqueue("b")

		v, ok := q.Dequeue()
		if !ok || v != "a" {
			t.Errorf("expected a, got %q (ok=%v)", v, ok)
		}
		v, ok = q.Dequeue()
		if !ok || v != "b" {
			t.Errorf("expected b, got %q (ok=%v)", v, ok)
		}
		_, ok = q.Dequeue()
		if ok {
			t.Error("expected empty")
		}
	})

	t.Run("int queue", func(t *testing.T) {
		q := queue.New[int]()
		q.Enqueue(10)
		q.Enqueue(20)
		x, ok := q.Dequeue()
		if !ok || x != 10 {
			t.Errorf("expected 10, got %d (ok=%v)", x, ok)
		}
		x, ok = q.Dequeue()
		if !ok || x != 20 {
			t.Errorf("expected 20, got %d (ok=%v)", x, ok)
		}
	})

	t.Run("empty queue", func(t *testing.T) {
		q := queue.New[int]()
		_, ok := q.Dequeue()
		if ok {
			t.Error("expected false from empty queue")
		}
	})

	t.Run("single element", func(t *testing.T) {
		q := queue.New[int]()
		q.Enqueue(42)
		if q.Len() != 1 {
			t.Errorf("Len = %d, want 1", q.Len())
		}
		v, ok := q.Dequeue()
		if !ok || v != 42 {
			t.Errorf("expected 42, got %d", v)
		}
		if q.Len() != 0 {
			t.Errorf("Len after dequeue = %d, want 0", q.Len())
		}
	})

	t.Run("enqueue dequeue interleaved", func(t *testing.T) {
		q := queue.New[int]()
		for i := range 100 {
			q.Enqueue(i)
			v, ok := q.Dequeue()
			if !ok || v != i {
				t.Fatalf("expected %d, got %d", i, v)
			}
		}
	})

	t.Run("many elements", func(t *testing.T) {
		n := 1000
		q := queue.New[int](queue.WithCapacity(n / 10))
		for i := range n {
			q.Enqueue(i)
		}
		if q.Len() != n {
			t.Errorf("Len = %d, want %d", q.Len(), n)
		}
		for i := range n {
			v, ok := q.Dequeue()
			if !ok || v != i {
				t.Fatalf("at %d: expected %d, got %d (ok=%v)", i, i, v, ok)
			}
		}
		if q.Len() != 0 {
			t.Errorf("final Len = %d, want 0", q.Len())
		}
	})

	t.Run("with capacity option", func(t *testing.T) {
		q := queue.New[int](queue.WithCapacity(64))
		for i := range 100 {
			q.Enqueue(i)
		}
		if q.Len() != 100 {
			t.Errorf("Len = %d, want 100", q.Len())
		}
	})

	t.Run("zero capacity option", func(t *testing.T) {
		q := queue.New[int](queue.WithCapacity(0))
		q.Enqueue(1)
		v, ok := q.Dequeue()
		if !ok || v != 1 {
			t.Errorf("expected 1, got %d", v)
		}
	})
}

func TestDequeueContext(t *testing.T) {
	t.Parallel()

	t.Run("immediate value", func(t *testing.T) {
		ctx := context.Background()
		q := queue.New[int]()
		q.Enqueue(99)
		v, ok := q.DequeueContext(ctx)
		if !ok || v != 99 {
			t.Errorf("expected 99, got %d (ok=%v)", v, ok)
		}
	})

	t.Run("block until enqueue", func(t *testing.T) {
		ctx := context.Background()
		q := queue.New[int]()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			q.Enqueue(42)
		}()

		v, ok := q.DequeueContext(ctx)
		if !ok || v != 42 {
			t.Errorf("expected 42, got %d (ok=%v)", v, ok)
		}
		wg.Wait()
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // already cancelled

		q := queue.New[int]()
		v, ok := q.DequeueContext(ctx)
		if ok {
			t.Errorf("expected false, got %d", v)
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		q := queue.New[int]()
		v, ok := q.DequeueContext(ctx)
		if ok {
			t.Errorf("expected false on timeout, got %d", v)
		}
	})

	t.Run("value before cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		q := queue.New[int]()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			q.Enqueue(77)
		}()

		v, ok := q.DequeueContext(ctx)
		if !ok || v != 77 {
			t.Errorf("expected 77 before cancel, got %d (ok=%v)", v, ok)
		}
		wg.Wait()
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	q := queue.New[int]()

	var wg sync.WaitGroup
	n := 100

	// Concurrent enqueues.
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range n {
				q.Enqueue(i)
			}
		}()
	}

	// Concurrent dequeues.
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range n {
				q.Dequeue()
			}
		}()
	}

	// Concurrent Len.
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range n {
				q.Len()
			}
		}()
	}

	wg.Wait()
}

func BenchmarkEnqueueDequeue(b *testing.B) {
	q := queue.New[int]()
	for b.Loop() {
		q.Enqueue(42)
		q.Dequeue()
	}
}

func BenchmarkEnqueueDequeueBatched(b *testing.B) {
	q := queue.New[int]()
	for b.Loop() {
		for range 100 {
			q.Enqueue(1)
		}
		for range 100 {
			q.Dequeue()
		}
	}
}
