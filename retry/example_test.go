package retry_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/azghr/forge/retry"
)

func ExampleRetryContext() {
	attempt := 0
	err := retry.RetryContext(context.Background(), retry.RetryConfig{
		MaxTries:   3,
		InitDelay:  10 * time.Millisecond,
		Multiplier: 2,
	}, func() error {
		attempt++
		if attempt < 2 {
			return fmt.Errorf("transient error")
		}
		return nil
	})
	if err != nil {
		fmt.Println("all attempts failed:", err)
		return
	}
	fmt.Println("succeeded after", attempt, "attempt(s)")
	// Output: succeeded after 2 attempt(s)
}

func ExampleRetryContext_http() {
	// Example of retrying an HTTP request that returns server errors.
	err := retry.RetryContext(context.Background(), retry.RetryConfig{
		MaxTries:   3,
		InitDelay:  50 * time.Millisecond,
		Multiplier: 2.0,
	}, func() error {
		resp, err := http.Get("http://example.com")
		if err != nil {
			return fmt.Errorf("network error: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: %d", resp.StatusCode)
		}
		return nil
	})
	if err != nil {
		fmt.Println("request eventually failed:", err)
	}
}

func ExampleRetryContext_cancellation() {
	// Context timeout interrupts retries early.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err := retry.RetryContext(ctx, retry.RetryConfig{
		MaxTries:   5,
		InitDelay:  10 * time.Millisecond,
		Multiplier: 2,
	}, func() error {
		return fmt.Errorf("slow")
	})
	fmt.Println(err == context.DeadlineExceeded)
	// Output: true
}
