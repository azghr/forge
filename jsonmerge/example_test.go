package jsonmerge_test

import (
	"fmt"
	"sort"

	"github.com/azghr/forge/jsonmerge"
)

func ExampleMerge() {
	a := map[string]interface{}{"x": 1, "y": map[string]interface{}{"v": 2}}
	b := map[string]interface{}{"y": map[string]interface{}{"v": 3}, "z": 4}
	jsonmerge.Merge(a, b)
	fmt.Println(a)
	// Output:
	// map[x:1 y:map[v:3] z:4]
}

func ExampleMerge_sliceAppend() {
	a := map[string]interface{}{"items": []interface{}{1, 2}}
	b := map[string]interface{}{"items": []interface{}{3, 4}}
	jsonmerge.Merge(a, b, jsonmerge.WithSliceMode(jsonmerge.SliceAppend))
	fmt.Println(a)
	// Output:
	// map[items:[1 2 3 4]]
}

func ExampleDiff() {
	a := map[string]interface{}{"x": 1, "y": map[string]interface{}{"v": 2, "w": 3}}
	b := map[string]interface{}{"y": map[string]interface{}{"v": 3}, "z": 4}
	got := jsonmerge.Diff(a, b)
	sort.Strings(got)
	fmt.Println(got)
	// Output:
	// [x y.v y.w]
}
