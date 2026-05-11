package option_test

import (
	"fmt"

	"github.com/azghr/forge/option"
)

func ExampleSome() {
	find := func(m map[string]int, key string) option.Option[int] {
		if v, ok := m[key]; ok {
			return option.Some(v)
		}
		return option.None[int]()
	}

	o := find(map[string]int{"a": 1}, "a")
	v, ok := o.Unwrap()
	fmt.Println(v, ok)
	// Output: 1 true
}

func ExampleNone() {
	o := option.None[string]()
	_, ok := o.Unwrap()
	fmt.Println(ok)
	// Output: false
}

func ExampleOption_Must() {
	defer func() {
		fmt.Println(recover())
	}()

	_ = option.None[int]().Must()
	// Output: option: Must() called on None
}
