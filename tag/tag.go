package tag

type Set map[string]struct{}

func NewSet(tags ...string) Set {
	s := Set{}
	for _, t := range tags {
		s[t] = struct{}{}
	}
	return s
}

func (s Set) Has(t string) bool {
	if s == nil {
		return false
	}
	_, ok := s[t]
	return ok
}
