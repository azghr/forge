package lockutil_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/azghr/forge/lockutil"
)

func ExampleTryLockMutex() {
	var mu sync.Mutex
	if lockutil.TryLockMutex(&mu) {
		defer mu.Unlock()
		fmt.Println("Mutex acquired")
	} else {
		fmt.Println("Mutex busy")
	}
	// Output: Mutex acquired
}

func ExampleTryLockRW() {
	var rw sync.RWMutex
	if lockutil.TryLockRW(&rw) {
		defer rw.RUnlock()
		fmt.Println("Read lock acquired")
	} else {
		fmt.Println("Write lock held")
	}
	// Output: Read lock acquired
}

func ExampleLockMutex() {
	var mu sync.Mutex
	mu.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if lockutil.LockMutex(ctx, &mu) {
		mu.Unlock()
		fmt.Println("Lock acquired after wait")
	} else {
		fmt.Println("Context expired before lock acquired")
	}
	mu.Unlock()
	// Output: Context expired before lock acquired
}

func ExampleLockRW() {
	var rw sync.RWMutex
	rw.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if lockutil.LockRW(ctx, &rw) {
		rw.RUnlock()
		fmt.Println("Read lock acquired after wait")
	} else {
		fmt.Println("Context expired before read lock acquired")
	}
	rw.Unlock()
	// Output: Context expired before read lock acquired
}

func ExampleLockMutex_customInterval() {
	var mu sync.Mutex
	mu.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	// Use a faster poll interval for lower-latency lock detection.
	_ = lockutil.LockMutex(ctx, &mu, lockutil.WithPollInterval(time.Microsecond))
	mu.Unlock()
}
