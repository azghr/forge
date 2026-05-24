package set_test

import (
	"testing"

	"github.com/azghr/forge/set"
)

func TestNew(t *testing.T) {
	s := set.New(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("expected len 3, got %d", s.Len())
	}
}

func TestAdd(t *testing.T) {
	s := set.New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(1)
	if s.Len() != 2 {
		t.Errorf("expected len 2 after duplicate add, got %d", s.Len())
	}
}

func TestRemove(t *testing.T) {
	s := set.New(1, 2, 3)
	s.Remove(2)
	if s.Contains(2) {
		t.Error("expected 2 to be removed")
	}
	if s.Len() != 2 {
		t.Errorf("expected len 2, got %d", s.Len())
	}
}

func TestContains(t *testing.T) {
	s := set.New("a", "b")
	if !s.Contains("a") {
		t.Error("expected contains a")
	}
	if s.Contains("c") {
		t.Error("expected not contains c")
	}
}

func TestItems(t *testing.T) {
	s := set.New(1, 2, 3)
	items := s.Items()
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestUnion(t *testing.T) {
	a := set.New(1, 2)
	b := set.New(2, 3)
	u := a.Union(b)
	for _, v := range []int{1, 2, 3} {
		if !u.Contains(v) {
			t.Errorf("expected union to contain %d", v)
		}
	}
	if u.Len() != 3 {
		t.Errorf("expected union len 3, got %d", u.Len())
	}
}

func TestIntersect(t *testing.T) {
	a := set.New(1, 2, 3)
	b := set.New(2, 3, 4)
	i := a.Intersect(b)
	for _, v := range []int{2, 3} {
		if !i.Contains(v) {
			t.Errorf("expected intersection to contain %d", v)
		}
	}
	if i.Contains(1) || i.Contains(4) {
		t.Error("intersection should not contain non-overlapping items")
	}
}

func TestDifference(t *testing.T) {
	a := set.New(1, 2, 3)
	b := set.New(2, 3, 4)
	d := a.Difference(b)
	if !d.Contains(1) {
		t.Error("expected difference to contain 1")
	}
	if d.Contains(2) || d.Contains(4) {
		t.Error("difference should not contain 2 or 4")
	}
}

func TestEmpty(t *testing.T) {
	s := set.New[int]()
	if s.Len() != 0 {
		t.Errorf("expected empty set len 0, got %d", s.Len())
	}
	if s.Contains(0) {
		t.Error("empty set should not contain anything")
	}
}
