package common

import (
	"fmt"
	"log"

	"github.com/avoronkov/waver/lib/seq/types"
)

func Repeat(idx *Index, times, fn types.ValueFn) types.ValueFn {
	return &repeatImpl{
		times: times,
		fn:    fn,
		idx:   idx,
	}
}

type repeatImpl struct {
	times types.ValueFn
	fn    types.ValueFn

	idx   *Index
	value types.Value
}

func (s *repeatImpl) Val(bit int64, ctx types.Context) types.Value {
	nTimes := s.times.Val(bit, ctx)
	intTimes, ok := nTimes.(Num)
	if !ok {
		panic(fmt.Errorf("repeat expects integer as first argument, found: %v (%T)", nTimes, nTimes))
	}
	log.Printf("Repeat idx: %v, times: %v", s.idx.N, intTimes)
	if s.idx.N >= int(intTimes) {
		s.idx.N = 0
	}
	if s.idx.N == 0 {
		s.value = s.fn.Val(bit, ctx)
	}
	s.idx.N++

	return s.value
}
