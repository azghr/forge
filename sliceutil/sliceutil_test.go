package sliceutil_test

import (
	"reflect"
	"testing"

	"github.com/azghr/forge/sliceutil"
)

func TestMap(t *testing.T) {
	t.Parallel()

	t.Run("double ints", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Map(data, func(x int) int { return x * 2 })
		want := []int{2, 4, 6, 8, 10}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("string length", func(t *testing.T) {
		data := []string{"a", "ab", "abc"}
		got := sliceutil.Map(data, func(s string) int { return len(s) })
		want := []int{1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := sliceutil.Map([]int{}, func(x int) int { return x * 2 })
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("type conversion", func(t *testing.T) {
		data := []int{1, 2, 3}
		got := sliceutil.Map(data, func(x int) float64 { return float64(x) })
		want := []float64{1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("evens", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Filter(data, func(x int) bool { return x%2 == 0 })
		want := []int{2, 4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("greater than 3", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Filter(data, func(x int) bool { return x > 3 })
		want := []int{4, 5}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := sliceutil.Filter([]int{}, func(x int) bool { return true })
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("all match", func(t *testing.T) {
		data := []int{1, 2, 3}
		got := sliceutil.Filter(data, func(x int) bool { return x > 0 })
		if !reflect.DeepEqual(got, data) {
			t.Errorf("got %v, want %v", got, data)
		}
	})

	t.Run("none match", func(t *testing.T) {
		data := []int{1, 2, 3}
		got := sliceutil.Filter(data, func(x int) bool { return false })
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Parallel()

	t.Run("sum", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Reduce(data, 0, func(acc, x int) int { return acc + x })
		if got != 15 {
			t.Errorf("got %d, want 15", got)
		}
	})

	t.Run("product", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Reduce(data, 1, func(acc, x int) int { return acc * x })
		if got != 120 {
			t.Errorf("got %d, want 120", got)
		}
	})

	t.Run("string concatenation", func(t *testing.T) {
		data := []string{"a", "b", "c"}
		got := sliceutil.Reduce(data, "", func(acc string, x string) string { return acc + x })
		if got != "abc" {
			t.Errorf("got %q, want %q", got, "abc")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := sliceutil.Reduce([]int{}, 42, func(acc, x int) int { return acc + x })
		if got != 42 {
			t.Errorf("got %d, want 42", got)
		}
	})

	t.Run("type conversion", func(t *testing.T) {
		data := []int{1, 2, 3}
		got := sliceutil.Reduce(data, "", func(acc string, x int) string { return acc + string(rune('0'+x)) })
		if got != "123" {
			t.Errorf("got %q, want %q", got, "123")
		}
	})
}

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run("all positive", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		if !sliceutil.All(data, func(x int) bool { return x > 0 }) {
			t.Errorf("expected true")
		}
	})

	t.Run("one negative", func(t *testing.T) {
		data := []int{1, -2, 3}
		if sliceutil.All(data, func(x int) bool { return x > 0 }) {
			t.Errorf("expected false")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		if !sliceutil.All([]int{}, func(x int) bool { return false }) {
			t.Errorf("expected true for empty slice")
		}
	})

	t.Run("short circuit", func(t *testing.T) {
		data := []int{1, 0, 3}
		calls := 0
		sliceutil.All(data, func(x int) bool {
			calls++
			return x > 0
		})
		if calls != 2 {
			t.Errorf("expected 2 calls, got %d", calls)
		}
	})
}

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("has three", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		if !sliceutil.Any(data, func(x int) bool { return x == 3 }) {
			t.Errorf("expected true")
		}
	})

	t.Run("no match", func(t *testing.T) {
		data := []int{1, 2, 3}
		if sliceutil.Any(data, func(x int) bool { return x == 99 }) {
			t.Errorf("expected false")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		if sliceutil.Any([]int{}, func(x int) bool { return true }) {
			t.Errorf("expected false for empty slice")
		}
	})

	t.Run("short circuit", func(t *testing.T) {
		data := []int{1, 2, 3}
		calls := 0
		sliceutil.Any(data, func(x int) bool {
			calls++
			return x == 2
		})
		if calls != 2 {
			t.Errorf("expected 2 calls, got %d", calls)
		}
	})
}

func TestChunk(t *testing.T) {
	t.Parallel()

	t.Run("exact division", func(t *testing.T) {
		data := []int{1, 2, 3, 4}
		got := sliceutil.Chunk(data, 2)
		want := [][]int{{1, 2}, {3, 4}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("non-exact", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5}
		got := sliceutil.Chunk(data, 2)
		want := [][]int{{1, 2}, {3, 4}, {5}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("chunk size 1", func(t *testing.T) {
		data := []int{1, 2, 3}
		got := sliceutil.Chunk(data, 1)
		want := [][]int{{1}, {2}, {3}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("chunk size larger than slice", func(t *testing.T) {
		data := []int{1, 2}
		got := sliceutil.Chunk(data, 10)
		want := [][]int{{1, 2}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := sliceutil.Chunk([]int{}, 3)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("n equals 0 returns nil", func(t *testing.T) {
		got := sliceutil.Chunk([]int{1, 2}, 0)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("n negative returns nil", func(t *testing.T) {
		got := sliceutil.Chunk([]int{1, 2, 3}, -1)
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestConcurrentSafety(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}

	run := make(chan struct{})
	done := make(chan struct{}, 20)

	fn := func(x int) int { return x * 2 }

	for range 20 {
		go func() {
			<-run
			sliceutil.Map(data, fn)
			sliceutil.Filter(data, func(x int) bool { return x > 0 })
			sliceutil.Reduce(data, 0, func(acc, x int) int { return acc + x })
			sliceutil.All(data, func(x int) bool { return x > 0 })
			sliceutil.Any(data, func(x int) bool { return x == 1 })
			sliceutil.Chunk(data, 2)
			done <- struct{}{}
		}()
	}

	close(run)

	for range 20 {
		<-done
	}
}

func TestInputNotModified(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	orig := []int{1, 2, 3, 4, 5}

	sliceutil.Map(data, func(x int) int { return x * 2 })
	if !reflect.DeepEqual(data, orig) {
		t.Errorf("Map modified input: got %v", data)
	}

	sliceutil.Filter(data, func(x int) bool { return x > 3 })
	if !reflect.DeepEqual(data, orig) {
		t.Errorf("Filter modified input: got %v", data)
	}

	sliceutil.Reduce(data, 0, func(acc, x int) int { return acc + x })
	if !reflect.DeepEqual(data, orig) {
		t.Errorf("Reduce modified input: got %v", data)
	}

	sliceutil.Chunk(data, 2)
	if !reflect.DeepEqual(data, orig) {
		t.Errorf("Chunk modified input: got %v", data)
	}
}

func TestNamedSliceType(t *testing.T) {
	type IDs []int

	data := IDs{1, 2, 3, 4, 5}

	doubled := sliceutil.Map(data, func(x int) int { return x * 2 })
	want := []int{2, 4, 6, 8, 10}
	if !reflect.DeepEqual(doubled, want) {
		t.Errorf("Map with named type: got %v, want %v", doubled, want)
	}

	evens := sliceutil.Filter(data, func(x int) bool { return x%2 == 0 })
	if !reflect.DeepEqual(evens, IDs{2, 4}) {
		t.Errorf("Filter with named type: got %v", evens)
	}

	chunks := sliceutil.Chunk(data, 2)
	if len(chunks) != 3 {
		t.Errorf("Chunk with named type: expected 3 chunks, got %d", len(chunks))
	}
}

func BenchmarkMap(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	fn := func(x int) int { return x * 2 }

	b.ResetTimer()
	for b.Loop() {
		sliceutil.Map(data, fn)
	}
}

func BenchmarkFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	fn := func(x int) bool { return x%2 == 0 }

	b.ResetTimer()
	for b.Loop() {
		sliceutil.Filter(data, fn)
	}
}

func BenchmarkReduce(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	fn := func(acc, x int) int { return acc + x }

	b.ResetTimer()
	for b.Loop() {
		sliceutil.Reduce(data, 0, fn)
	}
}
