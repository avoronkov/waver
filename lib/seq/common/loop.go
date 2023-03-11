package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

func Loop(size, fn types.ValueFn) types.ValueFn {
	return &loopImpl{
		size: size,
		fn:   fn,
	}
}

type loopImpl struct {
	size   types.ValueFn
	fn     types.ValueFn
	values []types.Value
	index  int
}

func (s *loopImpl) Val(bit int64, ctx types.Context) types.Value {
	nTimes := s.size.Val(bit, ctx)
	size, ok := nTimes.(Num)
	if !ok {
		panic(fmt.Errorf("loop expects integer as first argument, found: %v (%T)", nTimes, nTimes))
	}
	if len(s.values) < int(size) {
		// get next value and store in array
		val := s.fn.Val(bit, ctx)
		s.values = append(s.values, val)
		return val
	}

	// get stored index value
	if s.index >= len(s.values) {
		s.index = 0
	}
	val := s.values[s.index]
	s.index++
	return val
}
