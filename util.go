package hypp

// option represents an optional value.
// If OK is true, then V contains the value that can be used.
// If OK is false, then value V should not be used.
type option[T any] struct {
	V  T
	OK bool
}

// set is an unordered collections of items.
type set[T comparable] map[T]struct{}

// Has returns true if the set contains the given value.
func (s set[T]) Has(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s[v]
	return ok
}

// Add adds the given value to the set if not present already.
func (s set[T]) Add(v T) {
	s[v] = struct{}{}
}
