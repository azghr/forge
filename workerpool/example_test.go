package workerpool_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/azghr/forge/workerpool"
)

func ExamplePool() {
	pool := workerpool.NewPool(3)
	for i := 0; i < 5; i++ {
		i := i
		pool.Submit(func() interface{} {
			time.Sleep(time.Duration(i) * 10 * time.Millisecond)
			return i * 2
		})
	}
	pool.Close()

	var results []int
	for r := range pool.Results {
		results = append(results, r.(int))
	}
	sort.Ints(results)
	fmt.Println(results)
	// Output: [0 2 4 6 8]
}

func ExamplePool_noTasks() {
	pool := workerpool.NewPool(2)
	pool.Close()

	_, ok := <-pool.Results
	fmt.Println(ok)
	// Output: false
}
