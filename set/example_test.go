package set_test

import (
	"fmt"

	"github.com/azghr/forge/set"
)

func Example() {
	s := set.New("apple", "banana", "cherry")
	s.Add("banana")
	fmt.Println(s.Contains("banana"))
	fmt.Println(s.Len())
	// Output:
	// true
	// 3
}

func Example_union() {
	a := set.New(1, 2, 3)
	b := set.New(3, 4, 5)
	u := a.Union(b)
	fmt.Println(u.Len())
	// Output: 5
}
