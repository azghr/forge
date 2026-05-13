package orderedset_test

import (
	"reflect"
	"testing"

	"github.com/azghr/forge/orderedset"
)

func TestAddRemove(t *testing.T) {
	t.Parallel()

	s := orderedset.New[string]()
	s.Add("a")
	s.Add("b")
	s.Add("a")

	got := s.Values()
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Values = %v; want %v", got, want)
	}

	s.Remove("a")
	got = s.Values()
	want = []string{"b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Remove failed: got %v", got)
	}

	if s.Len() != 1 {
		t.Errorf("Len = %d; want 1", s.Len())
	}
}

func TestUnionIntersect(t *testing.T) {
	t.Parallel()

	t.Run("union", func(t *testing.T) {
		a := orderedset.New([]int{1, 2}...)
		b := orderedset.New([]int{2, 3}...)
		a.Union(b)
		want := []int{1, 2, 3}
		if !reflect.DeepEqual(a.Values(), want) {
			t.Errorf("Union = %v; want %v", a.Values(), want)
		}
	})

	t.Run("intersect", func(t *testing.T) {
		a := orderedset.New([]int{1, 2}...)
		b := orderedset.New([]int{2, 3}...)
		a.Intersect(b)
		want := []int{2}
		if !reflect.DeepEqual(a.Values(), want) {
			t.Errorf("Intersect = %v; want %v", a.Values(), want)
		}
	})
}

func TestNewWithVariadic(t *testing.T) {
	t.Parallel()

	s := orderedset.New(3, 1, 2, 1, 3)
	want := []int{3, 1, 2}
	got := s.Values()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("New = %v; want %v", got, want)
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	s := orderedset.New("a", "b", "c")

	tests := []struct {
		elem string
		want bool
	}{
		{"a", true},
		{"b", true},
		{"c", true},
		{"d", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.elem, func(t *testing.T) {
			got := s.Contains(tt.elem)
			if got != tt.want {
				t.Errorf("Contains(%q) = %v; want %v", tt.elem, got, tt.want)
			}
		})
	}
}

func TestRemoveNonexistent(t *testing.T) {
	t.Parallel()

	s := orderedset.New(1, 2, 3)
	s.Remove(99)
	if s.Len() != 3 {
		t.Errorf("Len = %d; want 3", s.Len())
	}
}

func TestEmptySet(t *testing.T) {
	t.Parallel()

	s := orderedset.New[int]()

	if s.Len() != 0 {
		t.Errorf("empty Len = %d; want 0", s.Len())
	}

	if s.Contains(1) {
		t.Errorf("empty Contains should be false")
	}

	vals := s.Values()
	if len(vals) != 0 {
		t.Errorf("empty Values = %v; want []", vals)
	}
}

func TestUnionNilOther(t *testing.T) {
	t.Parallel()

	s := orderedset.New(1, 2, 3)
	s.Union(nil)
	want := []int{1, 2, 3}
	got := s.Values()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Union with nil = %v; want %v", got, want)
	}
}

func TestIntersectNilOther(t *testing.T) {
	t.Parallel()

	s := orderedset.New(1, 2, 3)
	s.Intersect(nil)
	if s.Len() != 0 {
		t.Errorf("Intersect with nil should be empty, got %v", s.Values())
	}
}

func TestIntersectDisjoint(t *testing.T) {
	t.Parallel()

	a := orderedset.New(1, 2, 3)
	b := orderedset.New(4, 5, 6)
	a.Intersect(b)
	if a.Len() != 0 {
		t.Errorf("Intersect disjoint should be empty, got %v", a.Values())
	}
}

func TestUnionPreservesOrder(t *testing.T) {
	t.Parallel()

	a := orderedset.New(1, 3, 5)
	b := orderedset.New(2, 4, 6)
	a.Union(b)
	want := []int{1, 3, 5, 2, 4, 6}
	got := a.Values()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Union order = %v; want %v", got, want)
	}
}

func TestIntersectPreservesOrder(t *testing.T) {
	t.Parallel()

	a := orderedset.New(1, 2, 3, 4, 5)
	b := orderedset.New(5, 3, 1)
	a.Intersect(b)
	want := []int{1, 3, 5}
	got := a.Values()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect order = %v; want %v", got, want)
	}
}

func TestConcurrentSafety(t *testing.T) {
	s := orderedset.New[int]()

	done := make(chan struct{}, 10)

	add := func() {
		for i := 0; i < 100; i++ {
			s.Add(i)
		}
		done <- struct{}{}
	}

	read := func() {
		for i := 0; i < 100; i++ {
			s.Contains(i)
			s.Values()
			s.Len()
		}
		done <- struct{}{}
	}

	for range 5 {
		go add()
		go read()
	}

	for range 10 {
		<-done
	}
}

func BenchmarkAdd(b *testing.B) {
	s := orderedset.New[int]()
	b.ResetTimer()
	for b.Loop() {
		s.Add(1)
	}
}

func BenchmarkContains(b *testing.B) {
	s := orderedset.New(1, 2, 3)
	b.ResetTimer()
	for b.Loop() {
		s.Contains(1)
		s.Contains(99)
	}
}

func BenchmarkUnion(b *testing.B) {
	a := orderedset.New(1, 2, 3)
	other := orderedset.New(3, 4, 5)
	b.ResetTimer()
	for b.Loop() {
		a.Union(other)
	}
}

func BenchmarkIntersect(b *testing.B) {
	a := orderedset.New(1, 2, 3)
	other := orderedset.New(3, 4, 5)
	b.ResetTimer()
	for b.Loop() {
		a.Intersect(other)
	}
}
