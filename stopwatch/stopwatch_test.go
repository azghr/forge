package stopwatch_test

import (
	"testing"
	"time"

	"github.com/azghr/forge/stopwatch"
)

func TestStopwatch(t *testing.T) {
	t.Parallel()

	t.Run("start stop measures positive elapsed", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Start()
		time.Sleep(time.Millisecond)
		sw.Stop()
		if sw.Elapsed() <= 0 {
			t.Error("expected positive elapsed")
		}
	})

	t.Run("cumulative across multiple start-stop cycles", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Start()
		time.Sleep(time.Millisecond)
		sw.Stop()
		prev := sw.Elapsed()

		sw.Start()
		time.Sleep(time.Millisecond)
		sw.Stop()
		if sw.Elapsed() <= prev {
			t.Error("elapsed did not increase after second cycle")
		}
	})

	t.Run("reset clears elapsed", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Start()
		time.Sleep(time.Millisecond)
		sw.Stop()
		sw.Reset()
		if sw.Elapsed() != 0 {
			t.Errorf("expected 0 after reset, got %v", sw.Elapsed())
		}
	})

	t.Run("stop without start is no-op", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Stop()
		if sw.Elapsed() != 0 {
			t.Errorf("expected 0, got %v", sw.Elapsed())
		}
	})

	t.Run("start restarts timer", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Start()
		time.Sleep(2 * time.Millisecond)
		sw.Start() // restart
		time.Sleep(time.Millisecond)
		sw.Stop()
		if sw.Elapsed() <= 0 {
			t.Error("expected positive after restart")
		}
		// elapsed should be ~1ms, not ~3ms
		if sw.Elapsed() >= 3*time.Millisecond {
			t.Errorf("expected ~1ms after restart, got %v", sw.Elapsed())
		}
	})

	t.Run("elapsed includes running time", func(t *testing.T) {
		var sw stopwatch.Stopwatch
		sw.Start()
		time.Sleep(time.Millisecond)
		if sw.Elapsed() <= 0 {
			t.Error("expected positive elapsed while running")
		}
		sw.Stop()
	})
}

func TestStopwatchResetWhileRunning(t *testing.T) {
	var sw stopwatch.Stopwatch
	sw.Start()
	time.Sleep(time.Millisecond)
	sw.Reset()
	if sw.Elapsed() != 0 {
		t.Errorf("expected 0 after reset while running, got %v", sw.Elapsed())
	}
	// should be able to start again cleanly
	sw.Start()
	time.Sleep(time.Millisecond)
	sw.Stop()
	if sw.Elapsed() <= 0 {
		t.Error("expected positive after start after reset")
	}
}

func TestStopwatchZeroValue(t *testing.T) {
	var sw stopwatch.Stopwatch
	if sw.Elapsed() != 0 {
		t.Errorf("expected zero-value Stopwatch to have 0 elapsed, got %v", sw.Elapsed())
	}
}

func BenchmarkStopwatch(b *testing.B) {
	var sw stopwatch.Stopwatch
	b.ResetTimer()
	for b.Loop() {
		sw.Start()
		sw.Stop()
	}
}

func BenchmarkStopwatchElapsed(b *testing.B) {
	var sw stopwatch.Stopwatch
	sw.Start()
	b.ResetTimer()
	for b.Loop() {
		sw.Elapsed()
	}
}
