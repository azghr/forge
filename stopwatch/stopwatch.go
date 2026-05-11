// Package stopwatch provides a simple stopwatch for benchmarking code blocks
// by measuring elapsed time with start/stop/reset controls.
package stopwatch

import "time"

// Stopwatch measures elapsed time. A zero-value Stopwatch is ready to use.
//
// Stopwatch is not concurrency-safe; concurrent access must be synchronized
// externally by the caller.
type Stopwatch struct {
	start   time.Time
	elapsed time.Duration
	running bool
}

// Start begins or restarts the stopwatch. If the stopwatch is already
// running, the timer is reset from this point forward.
func (sw *Stopwatch) Start() {
	sw.start = time.Now()
	sw.running = true
}

// Stop halts the stopwatch and adds the measured duration to the total
// elapsed time. Calling Stop when the stopwatch is not running is a no-op.
func (sw *Stopwatch) Stop() {
	if !sw.running {
		return
	}
	sw.elapsed += time.Since(sw.start)
	sw.running = false
}

// Reset clears the elapsed time and stops the stopwatch if it is running.
func (sw *Stopwatch) Reset() {
	sw.start = time.Time{}
	sw.elapsed = 0
	sw.running = false
}

// Elapsed returns the total elapsed time measured so far. If the stopwatch
// is currently running, the time since the last Start call is included.
func (sw *Stopwatch) Elapsed() time.Duration {
	if sw.running {
		return sw.elapsed + time.Since(sw.start)
	}
	return sw.elapsed
}
