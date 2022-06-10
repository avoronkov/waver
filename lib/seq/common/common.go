package common

import (
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/midisynth/udp"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/types"
)

var Scale notes.Scale

func Chain(fn types.Signaler, modifiers ...types.Modifier) types.Signaler {
	res := fn
	for _, md := range modifiers {
		res = md(res)
	}
	return res
}

func Sig(signal string) (types.SignalFn, error) {
	sig, err := udp.ParseMessage([]byte(signal), Scale)
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
	return types.ValueFunc(func(n int64, ctx types.Context) (res types.Value) {
		defer func() {
			if r := recover(); r != nil {
				res = Str(name)
			}
		}()
		return ctx.Get(name, n)
	})
}
