// Package option provides a generic Option (Maybe) type that may hold a value
// or none. It offers a safer alternative to nullable pointers and multi-return
// patterns for optional values.
package option

// Option wraps an optional value of type T. An Option is either Some (holds a
// value) or None (empty). It is immutable and therefore concurrency-safe by
// design.
type Option[T any] struct {
	value T
	ok    bool
}

// Some constructs an Option containing v.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}

// None returns an empty Option representing the absence of a value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// IsSome reports whether the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// Unwrap returns (value, true) if the Option contains a value, or (zero, false)
// if it is None. It never panics.
func (o Option[T]) Unwrap() (T, bool) {
	return o.value, o.ok
}

// Must returns the contained value if present, or panics if the Option is None.
// Use only when you are certain a value exists.
func (o Option[T]) Must() T {
	if !o.ok {
		panic("option: Must() called on None")
	}
	return o.value
}
