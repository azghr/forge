// Package orderedset provides a generic insertion-ordered set.
//
// It solves deduplication while preserving insertion order and supports
// set operations (Union, Intersect). An internal map[T]struct{} provides
// O(1) lookups while a slice maintains insertion order.
//
// A Set is safe for concurrent use when accessed through its methods.
// Cross-calls (e.g. a.Union(b) and b.Intersect(a) concurrently) may
// deadlock; callers should serialise such operations.
package orderedset

import "sync"

// Set is a generic ordered set of comparable elements.
// The zero value is not usable; use NewSet to construct a set.
type Set[T comparable] struct {
	mu    sync.RWMutex
	elems []T
	index map[T]struct{}
}

// NewSet returns an empty Set pre-populated with the given elements
// (duplicates are silently dropped, preserving first-occurrence order).
func NewSet[T comparable](elems ...T) *Set[T] {
	s := &Set[T]{
		index: make(map[T]struct{}, len(elems)),
	}
	if len(elems) > 0 {
		s.elems = make([]T, 0, len(elems))
		for _, v := range elems {
			if _, ok := s.index[v]; !ok {
				s.index[v] = struct{}{}
				s.elems = append(s.elems, v)
			}
		}
	}
	return s
}

// Add inserts v into the set if not already present, preserving insertion
// order. If v is already present the set is unchanged.
func (s *Set[T]) Add(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.add(v)
}

// add is the lock-free internal version of Add.
func (s *Set[T]) add(v T) {
	if _, ok := s.index[v]; ok {
		return
	}
	s.index[v] = struct{}{}
	s.elems = append(s.elems, v)
}

// Remove deletes v from the set, preserving the order of remaining elements.
// If v is not present the set is unchanged.
func (s *Set[T]) Remove(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.index[v]; !ok {
		return
	}
	delete(s.index, v)
	for i, e := range s.elems {
		if e == v {
			s.elems = append(s.elems[:i], s.elems[i+1:]...)
			return
		}
	}
}

// Contains returns true if v is in the set.
func (s *Set[T]) Contains(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.index[v]
	return ok
}

// Values returns the elements in insertion order. The returned slice is a
// copy so callers may mutate it freely.
func (s *Set[T]) Values() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]T, len(s.elems))
	copy(out, s.elems)
	return out
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.elems)
}

// Union modifies s to be the union of s and other. Elements from other are
// appended to the end of s in their order within other. If other is nil it
// is treated as empty.
func (s *Set[T]) Union(other *Set[T]) {
	if other == nil {
		return
	}

	var addrs []T
	other.mu.RLock()
	addrs = make([]T, len(other.elems))
	copy(addrs, other.elems)
	other.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range addrs {
		s.add(v)
	}
}

// Intersect modifies s to keep only elements that are also present in other.
// The relative order of survivors is preserved. If other is nil the result
// is an empty set.
func (s *Set[T]) Intersect(other *Set[T]) {
	var otherIdx map[T]struct{}
	if other != nil {
		other.mu.RLock()
		otherIdx = make(map[T]struct{}, len(other.index))
		for k := range other.index {
			otherIdx[k] = struct{}{}
		}
		other.mu.RUnlock()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(otherIdx) == 0 {
		s.elems = s.elems[:0]
		s.index = make(map[T]struct{})
		return
	}

	newElems := make([]T, 0, len(s.elems))
	for _, v := range s.elems {
		if _, ok := otherIdx[v]; ok {
			newElems = append(newElems, v)
		}
	}

	newIdx := make(map[T]struct{}, len(newElems))
	for _, v := range newElems {
		newIdx[v] = struct{}{}
	}

	s.elems = newElems
	s.index = newIdx
}
