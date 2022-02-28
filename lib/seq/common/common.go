package common

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/midisynth/udp"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func Chain(fn types.Signaler, modifiers ...types.Modifier) types.Signaler {
	res := fn
	for _, md := range modifiers {
		res = md(res)
	}
	return res
}

func Sig(signal string) (types.SignalFn, error) {
	sig, err := udp.ParseMessage([]byte(signal))
	if err != nil {
		return nil, err
	}
	return func(bit int64, ctx types.Context) []signals.Signal {
		return []signals.Signal{*sig}
	}, nil
}

func Shift(n types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			nVal := n.Val(bit, ctx).(Num)
			return fn.Eval(bit-int64(nVal), ctx)
		}
		return types.SignalFn(f)
	}
}

func Var(name string) types.ValueFn {
	return types.ValueFunc(func(n int64, ctx types.Context) types.Value {
		return ctx.Get(name, n)
	})
}
