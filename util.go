package hypp

// Option represents an optional value.
// If OK is true, then V contains the values that can be used.
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
		s[v] = struct{}{}
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

// Map is a wrapper around the build-in map type.
type Map[K comparable, V any] map[K]V

// Has returns true if the Map contains the given key.
// If does not panic if the Map is nil.
func (m Map[K, V]) Has(k K) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

// // Set sets the value v for the given key k.
// // It does not panic if the Map is nil.
// // Instead it initializes the map and sets the key-value pair.
// func (m *Map[K, V]) Set(k K, v V) {
// 	if *m == nil {
// 		*m = Map[K, V]{}
// 	}
// 	n := *m
// 	n[k] = v
// }

// Get returns the value for the given key.
// It does not panic if the Map is nil.
// Instead it returns the empty value of type V.
func (m Map[K, V]) Get(k K) V {
	if m == nil {
		var v V
		return v
	}
	return m[k]
}

// GetOption returns an Option for the given key.
// It does not panic if the Map is nil.
// Instead it returns an Option whose value should not be used.
func (m Map[K, V]) GetOption(k K) Option[V] {
	if m == nil {
		return Option[V]{}
	}
	v, ok := m[k]
	return Option[V]{V: v, OK: ok}
}

func (m Map[K, V]) Copy() Map[K, V] {
	if m == nil {
		return nil
	}
	out := Map[K, V]{}
	for k, v := range m {
		out[k] = v
	}
	return out
}
