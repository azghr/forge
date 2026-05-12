package stringutil_test

import (
	"fmt"

	"github.com/azghr/forge/stringutil"
)

func ExampleTitle() {
	fmt.Println(stringutil.Title("hello world"))
	fmt.Println(stringutil.Title("go is awesome"))
	fmt.Println(stringutil.Title("123 test"))
	// Output:
	// Hello World
	// Go Is Awesome
	// 123 Test
}

func ExampleSlug() {
	fmt.Println(stringutil.Slug("Hello World"))
	fmt.Println(stringutil.Slug("Go Lang Library"))
	fmt.Println(stringutil.Slug("Hello, Go!"))
	// Output:
	// hello-world
	// go-lang-library
	// hello-go
}

func ExampleSlug_options() {
	fmt.Println(stringutil.Slug("Hello World", stringutil.WithSeparator("_")))
	fmt.Println(stringutil.Slug("Hello World Foo", stringutil.WithMaxLength(11)))
	// Output:
	// hello_world
	// hello-world
}

func ExampleRemoveAccents() {
	fmt.Println(stringutil.RemoveAccents("café"))
	fmt.Println(stringutil.RemoveAccents("naïve"))
	fmt.Println(stringutil.RemoveAccents("crème brûlée"))
	// Output:
	// cafe
	// naive
	// creme brulee
}
