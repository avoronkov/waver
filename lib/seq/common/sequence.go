package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

type Index struct {
	N int
}

func Sequence(idx *Index, values types.ValueFn) types.ValueFn {
	return &sequenceImpl{
		fn:  values,
		idx: idx,
	}
}

type sequenceImpl struct {
	fn  types.ValueFn
	idx *Index
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
	if s.idx.N >= l {
		s.idx.N = 0
	}
	res := list[s.idx.N]
	s.idx.N = (s.idx.N + 1) % l
	return res.Val(bit, ctx)
}
