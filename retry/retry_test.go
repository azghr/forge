package retry_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/azghr/forge/retry"
)

func TestRetryContext(t *testing.T) {
	t.Parallel()

	t.Run("success on first attempt", func(t *testing.T) {
		err := retry.RetryContext(context.Background(), retry.RetryConfig{MaxTries: 3}, func() error {
			return nil
		})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("success after retry", func(t *testing.T) {
		var attempts int
		err := retry.RetryContext(context.Background(), retry.RetryConfig{
			MaxTries: 3, InitDelay: 1, Multiplier: 2,
		}, func() error {
			attempts++
			if attempts < 2 {
				return fmt.Errorf("transient error")
			}
			return nil
		})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if attempts != 2 {
			t.Errorf("expected 2 attempts, got %d", attempts)
		}
	})

	t.Run("all attempts fail", func(t *testing.T) {
		var attempts int
		err := retry.RetryContext(context.Background(), retry.RetryConfig{
			MaxTries: 3, InitDelay: 1, Multiplier: 1,
		}, func() error {
			attempts++
			return fmt.Errorf("fail %d", attempts)
		})
		if err == nil || err.Error() != "fail 3" {
			t.Errorf("expected 'fail 3', got %v", err)
		}
		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("context cancelled before call", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()
		start := time.Now()
		err := retry.RetryContext(ctx, retry.RetryConfig{
			MaxTries: 3, InitDelay: 1, Multiplier: 2,
		}, func() error { return fmt.Errorf("nope") })
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected DeadlineExceeded, got %v", err)
		}
		if time.Since(start) > time.Millisecond {
			t.Errorf("should return immediately, took %v", time.Since(start))
		}
	})

	t.Run("context cancelled midway through retries", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		var attempts int
		err := retry.RetryContext(ctx, retry.RetryConfig{
			MaxTries: 5, InitDelay: 1, Multiplier: 1,
		}, func() error {
			attempts++
			if attempts == 2 {
				cancel()
			}
			return fmt.Errorf("fail")
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected Canceled, got %v", err)
		}
	})
}

func TestRetryEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  retry.RetryConfig
		wantMin int
		wantMax int
	}{
		{
			name:    "zero max tries defaults to 1",
			config:  retry.RetryConfig{MaxTries: 0, InitDelay: 1},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name:    "negative max tries defaults to 1",
			config:  retry.RetryConfig{MaxTries: -5, InitDelay: 1},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name:    "single try",
			config:  retry.RetryConfig{MaxTries: 1, InitDelay: 1},
			wantMin: 1,
			wantMax: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attempts int
			retry.RetryContext(context.Background(), tt.config, func() error {
				attempts++
				return fmt.Errorf("fail")
			})
			if attempts < tt.wantMin || attempts > tt.wantMax {
				t.Errorf("expected %d-%d attempts, got %d",
					tt.wantMin, tt.wantMax, attempts)
			}
		})
	}
}

func TestRetryContextConcurrent(t *testing.T) {
	ctx := context.Background()
	config := retry.RetryConfig{MaxTries: 3, InitDelay: 1, Multiplier: 1}

	run := make(chan struct{})
	done := make(chan struct{}, 10)

	for range 10 {
		go func() {
			<-run
			retry.RetryContext(ctx, config, func() error {
				return fmt.Errorf("fail")
			})
			done <- struct{}{}
		}()
	}

	close(run)

	timeout := time.After(5 * time.Second)
	for range 10 {
		select {
		case <-done:
		case <-timeout:
			t.Fatal("timed out waiting for goroutines")
		}
	}
}

func BenchmarkRetryContext(b *testing.B) {
	config := retry.RetryConfig{MaxTries: 3, InitDelay: 1, Multiplier: 2}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		attempt := 0
		retry.RetryContext(ctx, config, func() error {
			attempt++
			if attempt < 3 {
				return fmt.Errorf("fail")
			}
			return nil
		})
	}
}
