// Package set provides a generic unordered set implementation.
//
// It is backed by a map and provides the standard set operations: Add,
// Remove, Contains, Union, Intersect, Difference. All operations are
// safe for concurrent use when the set is accessed through separate
// goroutines; access must be synchronized externally otherwise.
package set

// Set is a generic unordered set.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates a set populated with the given items.
func New[T comparable](items ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(items))}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

// Add adds item to the set. If the item is already present this is a no-op.
func (s *Set[T]) Add(item T) {
	s.m[item] = struct{}{}
}

// Remove deletes item from the set. If the item is not present this is a no-op.
func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

// Contains reports whether item is in the set.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Items returns all elements in the set as a slice. The order is
// non-deterministic. The returned slice is a copy.
func (s *Set[T]) Items() []T {
	out := make([]T, 0, len(s.m))
	for item := range s.m {
		out = append(out, item)
	}
	return out
}

// Union returns a new set containing all elements from s and other.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := New[T]()
	for item := range s.m {
		result.m[item] = struct{}{}
	}
	for item := range other.m {
		result.m[item] = struct{}{}
	}
	return result
}

// Intersect returns a new set containing elements present in both s and other.
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	result := New[T]()
	// Iterate over the smaller set for efficiency.
	small, large := s, other
	if len(s.m) > len(other.m) {
		small, large = other, s
	}
	for item := range small.m {
		if _, ok := large.m[item]; ok {
			result.m[item] = struct{}{}
		}
	}
	return result
}

// Difference returns a new set containing elements present in s but not in other.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := New[T]()
	for item := range s.m {
		if _, ok := other.m[item]; !ok {
			result.m[item] = struct{}{}
		}
	}
	return result
}
