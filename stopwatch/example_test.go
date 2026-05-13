package stopwatch_test

import (
	"fmt"
	"time"

	"github.com/azghr/forge/stopwatch"
)

func ExampleStopwatch() {
	var sw stopwatch.Stopwatch
	sw.Start()
	time.Sleep(100 * time.Millisecond)
	sw.Stop()
	fmt.Printf("elapsed: %v\n", sw.Elapsed().Round(100*time.Millisecond))
	// Output: elapsed: 100ms
}

func ExampleStopwatch_reset() {
	var sw stopwatch.Stopwatch
	sw.Start()
	time.Sleep(5 * time.Millisecond)
	sw.Stop()
	sw.Reset()
	fmt.Println(sw.Elapsed())
	// Output: 0s
}

func ExampleStopwatch_elapsed() {
	var sw stopwatch.Stopwatch
	if sw.Elapsed() == 0 {
		fmt.Println("zero value ready to use")
	}
	// Output: zero value ready to use
}

func ExampleStopwatch_stopWithoutStart() {
	var sw stopwatch.Stopwatch
	sw.Stop() // no-op
	fmt.Println(sw.Elapsed())
	// Output: 0s
}
