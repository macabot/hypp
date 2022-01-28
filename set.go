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
