package common

import (
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func Sequence(values ...types.ValueFn) types.ValueFn {
	return &sequenceImpl{fns: values}
}

type sequenceImpl struct {
	fns []types.ValueFn
	idx int
}

func (s *sequenceImpl) Val(bit int64, ctx types.Context) types.Value {
	res := s.fns[s.idx].Val(bit, ctx)
	s.idx = (s.idx + 1) % len(s.fns)
	return res
}
