package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

type ValueHolder struct {
	Value types.Value
}

func Repeat(idxName, valueName string, times, fn types.ValueFn) types.ValueFn {
	return &repeatImpl{
		times:     times,
		fn:        fn,
		idxName:   idxName,
		valueName: valueName,
	}
}

type repeatImpl struct {
	times types.ValueFn
	fn    types.ValueFn

	idxName   string
	valueName string
}

func (s *repeatImpl) Val(bit int64, ctx types.Context) types.Value {
	nTimes := s.times.Val(bit, ctx)
	intTimes, ok := nTimes.(Num)
	if !ok {
		panic(fmt.Errorf("repeat expects integer as first argument, found: %v (%T)", nTimes, nTimes))
	}
	index := 0
	if idx, ok := ctx.GlobalGet(s.idxName); ok {
		index = idx.(int)
	}
	if index >= int(intTimes) {
		index = 0
	}
	var value types.Value
	val, ok := ctx.GlobalGet(s.valueName)
	if index == 0 || !ok {
		value = s.fn.Val(bit, ctx)
		ctx.GlobalPut(s.valueName, value)
	} else {
		value = val.(types.Value)
	}
	index++
	ctx.GlobalPut(s.idxName, index)

	return value
}
