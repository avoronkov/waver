package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

type Index struct {
	N int
}

func Sequence(idxName string, values types.ValueFn) types.ValueFn {
	return &sequenceImpl{
		fn:        values,
		indexName: idxName,
	}
}

type sequenceImpl struct {
	fn        types.ValueFn
	indexName string
}

func (s *sequenceImpl) Val(bit int64, ctx types.Context) types.Value {
	values := s.fn.Val(bit, ctx)
	list, ok := values.(EvaluatedList)
	if !ok {
		panic(fmt.Errorf("seq expects list, found: %v", values))
	}
	l := list.Len()
	if l == 0 {
		panic(fmt.Errorf("seq expects non-empty list"))
	}
	index := 0
	if idx, ok := ctx.GlobalGet(s.indexName); ok {
		index = idx.(int)
	}
	if index >= l {
		index = 0
	}
	res := list.Get(index)
	index = (index + 1) % l
	ctx.GlobalPut(s.indexName, index)
	return res
}
