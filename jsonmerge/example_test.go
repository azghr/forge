package jsonmerge_test

import (
	"fmt"
	"sort"

	"github.com/azghr/forge/jsonmerge"
)

func ExampleMerge() {
	a := map[string]any{"x": 1, "y": map[string]any{"v": 2}}
	b := map[string]any{"y": map[string]any{"v": 3}, "z": 4}
	jsonmerge.Merge(a, b)
	fmt.Println(a)
	// Output:
	// map[x:1 y:map[v:3] z:4]
}

func ExampleMerge_sliceAppend() {
	a := map[string]any{"items": []any{1, 2}}
	b := map[string]any{"items": []any{3, 4}}
	jsonmerge.Merge(a, b, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
	fmt.Println(a)
	// Output:
	// map[items:[1 2 3 4]]
}

func ExampleDiff() {
	a := map[string]any{"x": 1, "y": map[string]any{"v": 2, "w": 3}}
	b := map[string]any{"y": map[string]any{"v": 3}, "z": 4}
	got := jsonmerge.Diff(a, b)
	sort.Strings(got)
	fmt.Println(got)
	// Output:
	// [x y.v y.w]
}
