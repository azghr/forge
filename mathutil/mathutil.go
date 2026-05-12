// Package mathutil provides small helper functions missing from the standard
// math package: Clamp, Sign, Lerp, GCD, and ApproxEqual.
//
// All functions are pure, concurrency-safe, and have no external dependencies
// beyond the standard library.
package mathutil

// Clamp confines x to the inclusive range [lo, hi].
// If lo > hi the result is undefined (the implementation swaps them).
func Clamp(x, lo, hi float64) float64 {
	if lo > hi {
		lo, hi = hi, lo
	}
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// Sign returns -1 if x is negative, +1 if x is positive, or 0 if x is zero.
// The sign of -0 is 0.
func Sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

// Lerp performs linear interpolation between a and b by t in [0,1].
// When t=0 the result is a; when t=1 the result is b. t is not clamped.
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// GCD returns the greatest common divisor of a and b using the Euclidean
// algorithm. GCD(0,0) returns 0.
func GCD(a, b int64) int64 {
	a, b = abs(a), abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// ApproxEqual returns true if a and b are within epsilon of each other.
// opts can specify a custom tolerance via WithEpsilon (defaults to 1e-9).
func ApproxEqual(a, b float64, opts ...Option) bool {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= cfg.epsilon
}
