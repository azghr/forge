package mathutil_test

import (
	"math"
	"testing"

	"github.com/azghr/forge/mathutil"
)

func TestClamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		x, lo, hi, want float64
	}{
		{name: "inside", x: 2, lo: 0, hi: 5, want: 2},
		{name: "above", x: 6, lo: 0, hi: 5, want: 5},
		{name: "below", x: -1, lo: 0, hi: 5, want: 0},
		{name: "equal min", x: 0, lo: 0, hi: 5, want: 0},
		{name: "equal max", x: 5, lo: 0, hi: 5, want: 5},
		{name: "swapped range", x: 2, lo: 5, hi: 0, want: 2},
		{name: "negative range", x: -10, lo: -5, hi: -1, want: -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mathutil.Clamp(tt.x, tt.lo, tt.hi)
			if got != tt.want {
				t.Errorf("Clamp(%v,%v,%v) = %v, want %v", tt.x, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{name: "positive", x: 3.0, want: 1},
		{name: "negative", x: -0.1, want: -1},
		{name: "zero", x: 0.0, want: 0},
		{name: "neg zero", x: math.Copysign(0, -1), want: 0},
		{name: "large positive", x: 1e308, want: 1},
		{name: "large negative", x: -1e308, want: -1},
		{name: "small positive", x: 1e-300, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mathutil.Sign(tt.x)
			if got != tt.want {
				t.Errorf("Sign(%v) = %v, want %v", tt.x, got, tt.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		a, b, t float64
		want    float64
	}{
		{name: "t=0", a: 0, b: 10, t: 0, want: 0},
		{name: "t=1", a: 0, b: 10, t: 1, want: 10},
		{name: "t=0.5", a: 0, b: 10, t: 0.5, want: 5},
		{name: "negative a", a: -10, b: 10, t: 0.5, want: 0},
		{name: "t outside (low)", a: 0, b: 10, t: -0.5, want: -5},
		{name: "t outside (high)", a: 0, b: 10, t: 1.5, want: 15},
		{name: "same values", a: 5, b: 5, t: 0.5, want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mathutil.Lerp(tt.a, tt.b, tt.t)
			if got != tt.want {
				t.Errorf("Lerp(%v,%v,%v) = %v, want %v", tt.a, tt.b, tt.t, got, tt.want)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		a, b, want int64
	}{
		{name: "8 and 12", a: 8, b: 12, want: 4},
		{name: "zero and n", a: 0, b: 5, want: 5},
		{name: "n and zero", a: 7, b: 0, want: 7},
		{name: "both zero", a: 0, b: 0, want: 0},
		{name: "negative", a: -8, b: 12, want: 4},
		{name: "both negative", a: -8, b: -12, want: 4},
		{name: "coprime", a: 7, b: 13, want: 1},
		{name: "same", a: 9, b: 9, want: 9},
		{name: "one", a: 1, b: 100, want: 1},
		{name: "large", a: 1024, b: 256, want: 256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mathutil.GCD(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("GCD(%v,%v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestApproxEqual(t *testing.T) {
	t.Parallel()

	t.Run("default epsilon", func(t *testing.T) {
		if !mathutil.ApproxEqual(1.0, 1.0) {
			t.Error("identical values should be equal")
		}
		if !mathutil.ApproxEqual(1.0, 1.0+1e-10) {
			t.Error("within default epsilon should be equal")
		}
		if mathutil.ApproxEqual(1.0, 1.0+1e-8) {
			t.Error("beyond default epsilon should not be equal")
		}
	})

	t.Run("custom epsilon via ApproxEqualEpsilon", func(t *testing.T) {
		if !mathutil.ApproxEqualEpsilon(1.0, 1.5, 0.6) {
			t.Error("expected equal with wide epsilon")
		}
		if mathutil.ApproxEqualEpsilon(1.0, 1.5, 0.4) {
			t.Error("expected not equal with tight epsilon")
		}
	})

	t.Run("zero epsilon defaults to DefaultEpsilon", func(t *testing.T) {
		if !mathutil.ApproxEqualEpsilon(1.0, 1.0, 0) {
			t.Error("zero epsilon should use default and match identical values")
		}
	})
}

func TestConcurrentSafety(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	for range 10 {
		go func() {
			for range 100 {
				mathutil.Clamp(5, 0, 10)
				mathutil.Sign(-3)
				mathutil.Lerp(0, 10, 0.5)
				mathutil.GCD(12, 8)
				mathutil.ApproxEqual(1.0, 1.0)
			}
			done <- struct{}{}
		}()
	}
	for range 10 {
		<-done
	}
}

func BenchmarkClamp(b *testing.B) {
	for b.Loop() {
		mathutil.Clamp(5, 0, 10)
	}
}

func BenchmarkSign(b *testing.B) {
	for b.Loop() {
		mathutil.Sign(-3.5)
	}
}

func BenchmarkLerp(b *testing.B) {
	for b.Loop() {
		mathutil.Lerp(0, 10, 0.5)
	}
}

func BenchmarkGCD(b *testing.B) {
	for b.Loop() {
		mathutil.GCD(123456, 7890)
	}
}

func BenchmarkApproxEqual(b *testing.B) {
	for b.Loop() {
		mathutil.ApproxEqual(1.0, 1.0+1e-10)
	}
}
