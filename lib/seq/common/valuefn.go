package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

func Const(n int64) types.ValueFn {
	return types.ValueFunc(func(int64, types.Context) types.Value {
		return Num(n)
	})
}

func Lst(values ...int64) types.ValueFn {
	var res List
	for _, v := range values {
		res = append(res, Num(v))
	}
	return types.ValueFunc(func(int64, types.Context) types.Value {
		return res
	})
}
