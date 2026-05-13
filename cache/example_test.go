package cache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/azghr/forge/cache"
)

func Example() {
	c := cache.New[string, int](5 * time.Second)
	c.Set("x", 42)
	v, ok := c.Get("x")
	fmt.Println(v, ok)
	// Output: 42 true
}

func ExampleCache_ttl() {
	c := cache.New[string, int](10 * time.Millisecond)
	c.Set("a", 1)
	time.Sleep(15 * time.Millisecond)
	_, ok := c.Get("a")
	fmt.Println(ok)
	// Output: false
}

func ExampleCache_noExpiry() {
	c := cache.New[int, string](0)
	c.Set(1, "one")
	v, ok := c.Get(1)
	fmt.Println(v, ok)
	// Output: one true
}

func ExampleCache_GetOrLoad() {
	c := cache.New[string, int](time.Minute)

	v, err := c.GetOrLoad(context.Background(), "key", func(ctx context.Context) (int, error) {
		return 99, nil
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(v)
	// Output: 99
}
