package utils

type Set[T comparable] struct {
	values map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		values: map[T]bool{},
	}
}

func (s *Set[T]) Add(values ...T) {
	for _, value := range values {
		s.values[value] = true
	}
}

func (s *Set[T]) Values() []T {
	values := make([]T, 0, len(s.values))
	for v := range s.values {
		values = append(values, v)
	}
	return values
}
