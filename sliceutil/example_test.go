package sliceutil_test

import (
	"fmt"

	"github.com/azghr/forge/sliceutil"
)

func ExampleMap() {
	ints := []int{1, 2, 3, 4}
	sqs := sliceutil.Map(ints, func(x int) int { return x * x })
	fmt.Println(sqs)
	// Output: [1 4 9 16]
}

func ExampleFilter() {
	ints := []int{1, 2, 3, 4}
	evens := sliceutil.Filter(ints, func(x int) bool { return x%2 == 0 })
	fmt.Println(evens)
	// Output: [2 4]
}

func ExampleReduce() {
	ints := []int{1, 2, 3, 4}
	sum := sliceutil.Reduce(ints, 0, func(acc, x int) int { return acc + x })
	fmt.Println(sum)
	// Output: 10
}

func ExampleAll() {
	ints := []int{1, 2, 3, 4}
	allPos := sliceutil.All(ints, func(x int) bool { return x > 0 })
	fmt.Println(allPos)
	// Output: true
}

func ExampleAny() {
	ints := []int{1, 2, 3, 4}
	hasThree := sliceutil.Any(ints, func(x int) bool { return x == 3 })
	fmt.Println(hasThree)
	// Output: true
}

func ExampleChunk() {
	ints := []int{1, 2, 3, 4, 5}
	chunks := sliceutil.Chunk(ints, 3)
	fmt.Println(chunks)
	// Output: [[1 2 3] [4 5]]
}
