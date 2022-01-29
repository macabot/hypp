package hypp

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	s := Set[T]{}
	for _, v := range values {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[T]) Has(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s[v]
	return ok
}

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Has(k K) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

func (m *Map[K, V]) Set(k K, v V) {
	if *m == nil {
		*m = Map[K, V]{}
	}
	n := *m
	n[k] = v
}
