package rate_test

import (
	"context"
	"testing"
	"time"

	"github.com/azghr/forge/rate"
)

func TestAllow(t *testing.T) {
	l := rate.New(10, 5)

	// First 5 should be allowed (burst).
	for range 5 {
		if !l.Allow() {
			t.Error("expected allow within burst")
		}
	}
}

func TestAllowExceedBurst(t *testing.T) {
	l := rate.New(100, 3)

	for range 3 {
		l.Allow()
	}

	// Fourth should be blocked (exceeded burst).
	if l.Allow() {
		t.Error("expected deny after burst exhausted")
	}
}

func TestWait(t *testing.T) {
	l := rate.New(1000, 1)

	// Consume the initial token.
	l.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Wait failed: %v", err)
	}
}

func TestWaitContextCancelled(t *testing.T) {
	l := rate.New(1, 1)

	// Consume the only token.
	l.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx); err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestLimit(t *testing.T) {
	l := rate.New(10, 5)
	if l.Limit() != 10 {
		t.Errorf("expected limit 10, got %d", l.Limit())
	}
}

func TestBurst(t *testing.T) {
	l := rate.New(10, 5)
	if l.Burst() != 5 {
		t.Errorf("expected burst 5, got %d", l.Burst())
	}
}

func TestZeroRate(t *testing.T) {
	l := rate.New(0, 1)
	if l.Limit() <= 0 {
		t.Error("expected positive rate even when initialized with 0")
	}
}
