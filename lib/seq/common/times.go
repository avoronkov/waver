package common

import (
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/types"
)

func Times(n types.ValueFn, key string) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			totalTimes := int64(n.Val(bit, ctx).(Num))
			var currentTimes int64
			gt, ok := ctx.GlobalGet(key)
			if ok {
				currentTimes = gt.(int64)
			}
			if currentTimes >= totalTimes {
				return nil
			}
			currentTimes++
			ctx.GlobalPut(key, currentTimes)

			return fn.Eval(bit, ctx)
		}
		return types.SignalFn(f)
	}
}
