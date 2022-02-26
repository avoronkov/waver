package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

func Chain(fn types.Signaler, modifiers ...types.Modifier) types.Signaler {
	res := fn
	for _, md := range modifiers {
		res = md(res)
	}
	return res
}

func Sig(signal string) types.SignalFn {
	return func(bit int64, ctx types.Context) []string {
		return []string{signal}
	}
}

func Every(n int64) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []string {
			if bit%n == 0 {
				return fn.Eval(bit, ctx)
			}
			return nil
		}
		return types.SignalFn(f)
	}
}

func Shift(n int64) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []string {
			return fn.Eval(bit-n, ctx)
		}
		return types.SignalFn(f)
	}
}
