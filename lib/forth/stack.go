package forth

type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if l := len(s.stack); l > 0 {
		value := s.stack[l-1]
		s.stack = s.stack[0 : l-1]
		return value, true
	}
	var t T
	return t, false
}

func (s *Stack[T]) Top() (T, bool) {
	if l := len(s.stack); l > 0 {
		value := s.stack[l-1]
		return value, true
	}
	var t T
	return t, false
}
