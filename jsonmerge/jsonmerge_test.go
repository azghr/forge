package jsonmerge_test

import (
	"sort"
	"testing"

	"github.com/azghr/forge/jsonmerge"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("flat override", func(t *testing.T) {
		a := map[string]interface{}{"a": 1, "b": 2}
		b := map[string]interface{}{"b": 3, "c": 4}
		jsonmerge.Merge(a, b)
		if a["a"] != 1 || a["b"] != 3 || a["c"] != 4 {
			t.Errorf("flat merge failed: %v", a)
		}
	})

	t.Run("nested merge", func(t *testing.T) {
		a := map[string]interface{}{"a": 1, "nested": map[string]interface{}{"x": 10}}
		b := map[string]interface{}{"nested": map[string]interface{}{"x": 20}, "c": 3}
		jsonmerge.Merge(a, b)
		if a["c"] != 3 {
			t.Errorf("expected c=3, got %v", a["c"])
		}
		n := a["nested"].(map[string]interface{})
		if n["x"] != 20 {
			t.Errorf("expected nested.x=20, got %v", n["x"])
		}
	})

	t.Run("deep nested", func(t *testing.T) {
		a := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{"v": 1},
			},
		}
		b := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{"v": 2, "w": 3},
			},
		}
		jsonmerge.Merge(a, b)
		l1 := a["level1"].(map[string]interface{})
		l2 := l1["level2"].(map[string]interface{})
		if l2["v"] != 2 || l2["w"] != 3 {
			t.Errorf("deep nested merge failed: %v", l2)
		}
	})

	t.Run("new key added", func(t *testing.T) {
		a := map[string]interface{}{"x": 1}
		b := map[string]interface{}{"y": 2}
		jsonmerge.Merge(a, b)
		if a["x"] != 1 || a["y"] != 2 {
			t.Errorf("new key merge failed: %v", a)
		}
	})

	t.Run("non-map override", func(t *testing.T) {
		a := map[string]interface{}{"k": map[string]interface{}{"v": 1}}
		b := map[string]interface{}{"k": 42}
		jsonmerge.Merge(a, b)
		if a["k"] != 42 {
			t.Errorf("expected k=42, got %v", a["k"])
		}
	})

	t.Run("nil values", func(t *testing.T) {
		a := map[string]interface{}{"a": 1}
		b := map[string]interface{}{"a": nil}
		jsonmerge.Merge(a, b)
		if a["a"] != nil {
			t.Errorf("expected a=nil, got %v", a["a"])
		}
	})

	t.Run("empty maps", func(t *testing.T) {
		a := map[string]interface{}{}
		b := map[string]interface{}{"x": 1}
		jsonmerge.Merge(a, b)
		if a["x"] != 1 {
			t.Errorf("empty dst merge failed: %v", a)
		}
	})

	t.Run("slice default replace", func(t *testing.T) {
		a := map[string]interface{}{"s": []interface{}{1, 2}}
		b := map[string]interface{}{"s": []interface{}{3, 4}}
		jsonmerge.Merge(a, b)
		s := a["s"].([]interface{})
		if len(s) != 2 || s[0] != 3 || s[1] != 4 {
			t.Errorf("slice replace failed: %v", s)
		}
	})

	t.Run("slice append mode", func(t *testing.T) {
		a := map[string]interface{}{"s": []interface{}{1, 2}}
		b := map[string]interface{}{"s": []interface{}{3, 4}}
		jsonmerge.Merge(a, b, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
		s := a["s"].([]interface{})
		if len(s) != 4 || s[0] != 1 || s[1] != 2 || s[2] != 3 || s[3] != 4 {
			t.Errorf("slice append failed: %v", s)
		}
	})

	t.Run("slice mode non-slice values", func(t *testing.T) {
		a := map[string]interface{}{"k": 1}
		b := map[string]interface{}{"k": 2}
		jsonmerge.Merge(a, b, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
		if a["k"] != 2 {
			t.Errorf("non-slice append mode override failed: %v", a["k"])
		}
	})
}

func TestDiff(t *testing.T) {
	t.Parallel()

	t.Run("identical", func(t *testing.T) {
		a := map[string]interface{}{"x": 1, "y": 2}
		b := map[string]interface{}{"x": 1, "y": 2}
		got := jsonmerge.Diff(a, b)
		if len(got) != 0 {
			t.Errorf("expected no diffs, got %v", got)
		}
	})

	t.Run("different values", func(t *testing.T) {
		a := map[string]interface{}{"x": 1, "y": 2}
		b := map[string]interface{}{"x": 1, "y": 99}
		got := jsonmerge.Diff(a, b)
		if len(got) != 1 || got[0] != "y" {
			t.Errorf("expected [y], got %v", got)
		}
	})

	t.Run("nested diff", func(t *testing.T) {
		a := map[string]interface{}{"y": map[string]interface{}{"v": 2, "w": 3}}
		b := map[string]interface{}{"y": map[string]interface{}{"v": 3}, "z": 4}
		got := jsonmerge.Diff(a, b)
		sort.Strings(got)
		want := []string{"y.v", "y.w"}
		if !stringSliceEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("key missing in b", func(t *testing.T) {
		a := map[string]interface{}{"x": 1, "y": 2}
		b := map[string]interface{}{"x": 1}
		got := jsonmerge.Diff(a, b)
		if len(got) != 1 || got[0] != "y" {
			t.Errorf("expected [y], got %v", got)
		}
	})

	t.Run("key missing in a not reported", func(t *testing.T) {
		a := map[string]interface{}{"x": 1}
		b := map[string]interface{}{"x": 1, "y": 2}
		got := jsonmerge.Diff(a, b)
		if len(got) != 0 {
			t.Errorf("expected no diffs, got %v", got)
		}
	})

	t.Run("empty maps", func(t *testing.T) {
		got := jsonmerge.Diff(map[string]interface{}{}, map[string]interface{}{})
		if len(got) != 0 {
			t.Errorf("expected no diffs, got %v", got)
		}
	})

	t.Run("nil vs value", func(t *testing.T) {
		a := map[string]interface{}{"x": nil}
		b := map[string]interface{}{"x": 1}
		got := jsonmerge.Diff(a, b)
		if len(got) != 1 || got[0] != "x" {
			t.Errorf("expected [x], got %v", got)
		}
	})

	t.Run("slices differ", func(t *testing.T) {
		a := map[string]interface{}{"s": []interface{}{1, 2, 3}}
		b := map[string]interface{}{"s": []interface{}{1, 2, 4}}
		got := jsonmerge.Diff(a, b)
		if len(got) != 1 || got[0] != "s" {
			t.Errorf("expected [s], got %v", got)
		}
	})

	t.Run("all differ", func(t *testing.T) {
		a := map[string]interface{}{"a": 1, "b": 2}
		b := map[string]interface{}{"a": 10, "b": 20}
		got := jsonmerge.Diff(a, b)
		sort.Strings(got)
		want := []string{"a", "b"}
		if !stringSliceEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	run := make(chan struct{})
	done := make(chan struct{}, 30)

	for range 30 {
		go func() {
			<-run
			a := map[string]interface{}{"x": 1, "nested": map[string]interface{}{"v": 2}}
			b := map[string]interface{}{"nested": map[string]interface{}{"v": 3}, "z": 4}
			jsonmerge.Merge(a, b)
			jsonmerge.Merge(a, b, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
			jsonmerge.Diff(a, b)
			done <- struct{}{}
		}()
	}

	close(run)

	for range 30 {
		<-done
	}
}

func BenchmarkMerge(b *testing.B) {
	src := map[string]interface{}{
		"a": 1,
		"nested": map[string]interface{}{
			"x": 20,
			"y": 30,
		},
		"c": 3,
	}

	b.ResetTimer()
	for b.Loop() {
		dst := map[string]interface{}{
			"a": 10,
			"nested": map[string]interface{}{
				"x": 1,
			},
			"b": 2,
		}
		jsonmerge.Merge(dst, src)
	}
}

func BenchmarkMergeWithOptions(b *testing.B) {
	src := map[string]interface{}{
		"s": []interface{}{3, 4},
	}

	b.ResetTimer()
	for b.Loop() {
		dst := map[string]interface{}{
			"s": []interface{}{1, 2},
		}
		jsonmerge.Merge(dst, src, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
	}
}

func BenchmarkDiff(b *testing.B) {
	a := map[string]interface{}{
		"a": 1,
		"nested": map[string]interface{}{
			"x": 10,
			"y": 20,
		},
		"c": 3,
	}
	b2 := map[string]interface{}{
		"a": 1,
		"nested": map[string]interface{}{
			"x": 99,
			"y": 20,
		},
		"d": 4,
	}

	b.ResetTimer()
	for b.Loop() {
		jsonmerge.Diff(a, b2)
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
