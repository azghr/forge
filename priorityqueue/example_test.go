package priorityqueue_test

import (
	"fmt"

	"github.com/azghr/forge/priorityqueue"
)

func ExampleQueue() {
	pq := priorityqueue.New[string]()
	pq.Push(priorityqueue.Item[string]{Value: "a", Priority: 10})
	pq.Push(priorityqueue.Item[string]{Value: "b", Priority: 5})

	it, ok := pq.Pop()
	fmt.Println(it.Value, ok)

	it, ok = pq.Pop()
	fmt.Println(it.Value, ok)

	_, ok = pq.Pop()
	fmt.Println(ok)
	// Output:
	// b true
	// a true
	// false
}

func ExampleQueue_maxHeap() {
	pq := priorityqueue.New[string](priorityqueue.WithMaxHeap())
	pq.Push(priorityqueue.Item[string]{Value: "low", Priority: 1})
	pq.Push(priorityqueue.Item[string]{Value: "high", Priority: 10})

	it, _ := pq.Pop()
	fmt.Println(it.Value)

	it, _ = pq.Pop()
	fmt.Println(it.Value)
	// Output:
	// high
	// low
}
