package util

// Option represents an optional value.
// If OK is true, then V contains the value that can be used.
// If OK is false, then value V should not be used.
type Option[T any] struct {
	V  T
	OK bool
}

// Set is an unordered collections of items.
type Set[T comparable] map[T]struct{}

// NewSet creates a new Set.
func NewSet[T comparable](values ...T) Set[T] {
	s := Set[T]{}
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Has returns true if the Set contains the given value.
func (s Set[T]) Has(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s[v]
	return ok
}

// Add adds the given value in to the Set if not present already.
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}
