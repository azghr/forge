package rate_test

import (
	"context"
	"fmt"
	"time"

	"github.com/azghr/forge/rate"
)

func Example() {
	l := rate.New(2, 1)

	fmt.Println(l.Allow())
	fmt.Println(l.Allow())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	fmt.Println(l.Wait(ctx))
	// Output:
	// true
	// false
	// context deadline exceeded
}
