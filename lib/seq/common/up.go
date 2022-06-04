package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

func Up(shift, value types.ValueFn, invert bool) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		val := value.Val(bit, ctx)
		sh := shift.Val(bit, ctx)
		shInt, ok := sh.(Num)
		if !ok {
			panic(fmt.Errorf("up expects first argument to be number, found: %v (%T)", sh, sh))
		}
		if invert {
			shInt = -shInt
		}
		switch v := val.(type) {
		case Num:
			return Num(v + shInt)
		case List:
			res := []types.ValueFn{}
			for _, item := range v {
				res = append(res, Const(int64(item.Val(bit, ctx).(Num)+shInt)))
			}
			return List(res)
		default:
			panic(fmt.Errorf("up expects second argument to be number of list of numbers, found: %v (%T)", val, val))
		}
	}
	return types.ValueFunc(f)
}
