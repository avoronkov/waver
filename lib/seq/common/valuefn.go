package common

import "github.com/avoronkov/waver/lib/seq/types"

func Const(n int64) types.ValueFn {
	return types.ValueFunc(func(int64, types.Context) types.Value {
		return Num(n)
	})
}

func Lst(values ...types.ValueFn) types.ValueFn {
	return types.ValueFunc(func(bit int64, ctx types.Context) types.Value {
		return List(values)
	})
}

func FloatConst(f float64) types.ValueFn {
	return types.ValueFunc(func(bit int64, ctx types.Context) types.Value {
		return Float(f)
	})
}

func StrConst(s string) types.ValueFn {
	return types.ValueFunc(func(bit int64, ctx types.Context) types.Value {
		return Str(s)
	})
}
