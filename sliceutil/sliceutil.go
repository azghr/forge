// Package sliceutil provides generic collection functions for slices.
//
// It addresses Go's lack of built-in higher-order slice operations with a
// small set of stateless functions: Map, Filter, Reduce, All, Any, and Chunk.
// No side-effects are performed on input slices.
package sliceutil

// Map applies function f to each element of input slice s, returning a new
// slice of results. The input slice is not modified.
//
// Performance: O(n) time, O(n) allocated space.
func Map[S ~[]E, E any, R any](s S, f func(E) R) []R {
	if len(s) == 0 {
		return nil
	}
	out := make([]R, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

// Filter returns a new slice containing only elements of s for which f
// returns true. The input slice is not modified.
//
// Performance: O(n) time, O(n) worst-case allocated space.
func Filter[S ~[]E, E any](s S, f func(E) bool) S {
	if len(s) == 0 {
		return nil
	}
	out := make(S, 0, len(s))
	for _, v := range s {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reduce aggregates elements of s by applying f cumulatively, starting with
// init. For each element x in s, f(acc, x) produces the next accumulator
// value, which is returned after all elements have been processed.
//
// Performance: O(n) time, O(1) space.
func Reduce[S ~[]E, E any, R any](s S, init R, f func(acc R, x E) R) R {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// All returns true if f(x) is true for every element x in s, or true if s
// is empty. Short-circuits on the first false result.
//
// Performance: O(n) time worst-case, O(1) space.
func All[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

// Any returns true if f(x) is true for at least one element x in s.
// Short-circuits on the first true result.
//
// Performance: O(n) time worst-case, O(1) space.
func Any[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// Chunk splits slice s into chunks of size n, returned as a slice of slices.
// The last chunk may be smaller than n. If n <= 0 or s is empty, nil is
// returned.
//
// Performance: O(n) time, O(n) allocated space.
func Chunk[S ~[]E, E any](s S, n int) []S {
	if n <= 0 || len(s) == 0 {
		return nil
	}
	size := (len(s) + n - 1) / n
	out := make([]S, 0, size)
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
	}
	return out
}
