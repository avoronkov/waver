package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

type ValueHolder struct {
	Value types.Value
}

func Repeat(idx *Index, h *ValueHolder, times, fn types.ValueFn) types.ValueFn {
	return &repeatImpl{
		times: times,
		fn:    fn,
		idx:   idx,
		value: h,
	}
}

type repeatImpl struct {
	times types.ValueFn
	fn    types.ValueFn

	idx   *Index
	value *ValueHolder
}

func (s *repeatImpl) Val(bit int64, ctx types.Context) types.Value {
	nTimes := s.times.Val(bit, ctx)
	intTimes, ok := nTimes.(Num)
	if !ok {
		panic(fmt.Errorf("repeat expects integer as first argument, found: %v (%T)", nTimes, nTimes))
	}
	if s.idx.N >= int(intTimes) {
		s.idx.N = 0
	}
	if s.idx.N == 0 {
		s.value.Value = s.fn.Val(bit, ctx)
	}
	s.idx.N++

	return s.value.Value
}
