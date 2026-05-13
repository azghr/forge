package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/azghr/forge/cache"
	"github.com/azghr/forge/lockutil"
	"github.com/azghr/forge/mathutil"
	"github.com/azghr/forge/multityperror"
	"github.com/azghr/forge/orderedset"
	"github.com/azghr/forge/priorityqueue"
	"github.com/azghr/forge/queue"
	"github.com/azghr/forge/retry"
	"github.com/azghr/forge/workerpool"
)

type JobResult struct {
	Value   int
	IsPrime bool
	IsEven  bool
	Err     error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: pipeline <num1> [num2 ...]\n")
		os.Exit(1)
	}

	raw := os.Args[1:]

	// 1. Deduplicate using orderedset
	unique := orderedset.New[int]()
	for _, s := range raw {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %q: not a number\n", s)
			continue
		}
		unique.Add(n)
	}

	vals := unique.Values()
	if len(vals) == 0 {
		fmt.Fprintln(os.Stderr, "no valid numbers")
		os.Exit(1)
	}
	fmt.Printf("Unique values: %v\n", vals)

	// 2. Enqueue into priority queue (larger = higher priority)
	pq := priorityqueue.New[int]()
	for _, v := range vals {
		pq.Push(priorityqueue.Item[int]{Value: v, Priority: v})
	}

	// 3. Drain into FIFO queue
	q := queue.New[int]()
	for pq.Len() > 0 {
		item, _ := pq.Pop()
		q.Enqueue(item.Value)
	}

	// 4 & 5. Worker pool with cache
	pool := workerpool.New[JobResult](4, workerpool.WithTaskBuffer(16))
	calcCache := cache.New[int, JobResult](5 * time.Minute)

	// 6. Collect errors
	errs := multityperror.New()

	// 7. Lockutil for safe result collection
	var mu sync.Mutex
	var results []JobResult

	// 8. RetryConfig for transient failures
	rc := retry.RetryConfig{
		MaxTries:  3,
		InitDelay: 5 * time.Millisecond,
	}

	count := q.Len()
	for i := 0; i < count; i++ {
		val, _ := q.Dequeue()
		v := val

		pool.Submit(func() JobResult {
			if r, ok := calcCache.Get(v); ok {
				return r
			}

			var jr JobResult
			err := retry.RetryContext(context.Background(), rc, func() error {
				r, err := compute(val)
				jr = r
				return err
			})
				jr.Value = v
			if err != nil {
				jr.Err = err
				errs.Append(fmt.Errorf("value %d: %w", v, err))
			}

			calcCache.Set(v, jr)

			if !lockutil.TryLockMutex(&mu) {
				mu.Lock()
			}
			results = append(results, jr)
			mu.Unlock()

			return jr
		})
	}

	pool.Close()

	for r := range pool.Results {
		if r.Err != nil {
			fmt.Printf("Error: %v\n", r.Err)
		} else {
			fmt.Printf("Value %d: prime=%v even=%v\n", r.Value, r.IsPrime, r.IsEven)
		}
	}

	// 9. mathutil demonstration: find pair with smallest GCD gap
	if len(results) >= 2 {
		bestA, bestB, bestGCD := results[0].Value, results[1].Value,
			mathutil.GCD(int64(results[0].Value), int64(results[1].Value))
		for i := 0; i < len(results); i++ {
			for j := i + 1; j < len(results); j++ {
				g := mathutil.GCD(int64(results[i].Value), int64(results[j].Value))
				if g > bestGCD {
					bestA, bestB, bestGCD = results[i].Value, results[j].Value, g
				}
			}
		}
		fmt.Printf("\nClosest pair by GCD: (%d, %d) with GCD=%d\n", bestA, bestB, bestGCD)
	}

	// 10. ApproxEqual demonstration
	if len(results) >= 2 {
		a := float64(results[0].Value)
		b := float64(results[0].Value) * 1.0000000001
		fmt.Printf("ApproxEqual(%f, %f) = %v\n", a, b, mathutil.ApproxEqual(a, b))
		fmt.Printf("GCD(%d, %d) = %d\n", results[0].Value, results[1].Value,
			mathutil.GCD(int64(results[0].Value), int64(results[1].Value)))
	}

	if !errs.IsEmpty() {
		fmt.Fprintf(os.Stderr, "\nErrors: %s\n", errs.Error())
	}
}

func compute(n int) (JobResult, error) {
	if n%7 == 0 {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		return JobResult{}, fmt.Errorf("transient failure for %d", n)
	}
	time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
	return JobResult{
		IsPrime: isPrime(n),
		IsEven:  n%2 == 0,
	}, nil
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
