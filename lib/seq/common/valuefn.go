package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

func Const(n int64) types.ValueFn {
	return types.ValueFunc(func(int64, types.Context) types.Value {
		return Num(n)
	})
}

func Lst(values ...types.ValueFn) types.ValueFn {
	return types.ValueFunc(func(bit int64, ctx types.Context) types.Value {
		var res List
		for _, x := range values {
			res = append(res, x.Val(bit, ctx))
		}
		return res
	})
}
