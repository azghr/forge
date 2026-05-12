package queue_test

import (
	"context"
	"fmt"
	"time"

	"github.com/azghr/forge/queue"
)

func ExampleQueue() {
	q := queue.New[int]()
	q.Enqueue(10)
	q.Enqueue(20)

	x, ok := q.Dequeue()
	fmt.Println(x, ok)

	x, ok = q.Dequeue()
	fmt.Println(x, ok)

	_, ok = q.Dequeue()
	fmt.Println(ok)
	// Output:
	// 10 true
	// 20 true
	// false
}

func ExampleQueue_capacity() {
	q := queue.New[string](queue.WithCapacity(64))
	q.Enqueue("hello")
	q.Enqueue("world")
	fmt.Println(q.Len())
	// Output: 2
}

func ExampleQueue_dequeueContext() {
	q := queue.New[int]()

	go func() {
		time.Sleep(10 * time.Millisecond)
		q.Enqueue(99)
	}()

	ctx := context.Background()
	v, ok := q.DequeueContext(ctx)
	fmt.Println(v, ok)
	// Output: 99 true
}
