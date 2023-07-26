package forth

type Stack struct {
	stack []int
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(value int) {
	s.stack = append(s.stack, value)
}

func (s *Stack) Pop() (int, bool) {
	if l := len(s.stack); l > 0 {
		value := s.stack[l-1]
		s.stack = s.stack[0 : l-1]
		return value, true
	}
	return 0, false
}

func (s *Stack) Top() (int, bool) {
	if l := len(s.stack); l > 0 {
		value := s.stack[l-1]
		return value, true
	}
	return 0, false
}
