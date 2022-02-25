package main

import "gitlab.com/avoronkov/waver/lib/seq"

type Modifier = func(seq.Signaler) seq.Signaler

func Every(n int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			if bit%n == 0 {
				return fn.Eval(bit, ctx)
			}
			return nil
		}
		return seq.SignalFn(f)
	}
}

func Shift(n int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			return fn.Eval(bit-n, ctx)
		}
		return seq.SignalFn(f)
	}
}

func Sig(signal string) seq.SignalFn {
	return func(bit int64, ctx seq.Context) []string {
		return []string{signal}
	}
}

func OnBits(loop int64, bits ...int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			mod := bit % loop
			for _, b := range bits {
				if mod == b {
					return fn.Eval(bit, ctx)
				}
			}
			return nil
		}
		return seq.SignalFn(f)
	}
}

func Chain(fn seq.Signaler, modifiers ...Modifier) seq.Signaler {
	res := fn
	for _, md := range modifiers {
		res = md(res)
	}
	return res
}
