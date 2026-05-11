package orderedset_test

import (
	"fmt"

	"github.com/azghr/forge/orderedset"
)

func ExampleSet_Add() {
	s := orderedset.NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(1)
	fmt.Println(s.Values())
	// Output: [1 2]
}

func ExampleSet_Remove() {
	s := orderedset.NewSet(1, 2, 3)
	s.Remove(2)
	fmt.Println(s.Values())
	// Output: [1 3]
}

func ExampleSet_Contains() {
	s := orderedset.NewSet("a", "b", "c")
	fmt.Println(s.Contains("b"))
	fmt.Println(s.Contains("z"))
	// Output:
	// true
	// false
}

func ExampleSet_Union() {
	a := orderedset.NewSet([]int{1, 2, 3}...)
	b := orderedset.NewSet([]int{2, 3, 4}...)
	a.Union(b)
	fmt.Println(a.Values())
	// Output: [1 2 3 4]
}

func ExampleSet_Intersect() {
	a := orderedset.NewSet([]int{1, 2, 3}...)
	b := orderedset.NewSet([]int{2, 3, 4}...)
	a.Intersect(b)
	fmt.Println(a.Values())
	// Output: [2 3]
}

func ExampleNewSet() {
	s := orderedset.NewSet(3, 1, 2, 1, 3)
	fmt.Println(s.Values())
	fmt.Println(s.Len())
	// Output:
	// [3 1 2]
	// 3
}
