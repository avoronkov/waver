package common

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func Sequence(values types.ValueFn) types.ValueFn {
	return &sequenceImpl{fn: values}
}

type sequenceImpl struct {
	fn  types.ValueFn
	idx int
}

func (s *sequenceImpl) Val(bit int64, ctx types.Context) types.Value {
	values := s.fn.Val(bit, ctx)
	list, ok := values.(List)
	if !ok {
		panic(fmt.Errorf("seq expects list, found: %v", values))
	}
	l := len(list)
	if l == 0 {
		panic(fmt.Errorf("seq expects non-empty list"))
	}
	if s.idx >= l {
		s.idx = 0
	}
	res := list[s.idx]
	s.idx = (s.idx + 1) % l
	return res
}
