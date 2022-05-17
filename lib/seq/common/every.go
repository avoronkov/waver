package common

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func Every(n types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			val := n.Val(bit, ctx)
			if nVal, ok := val.(Num); ok {
				if bit%int64(nVal) == 0 {
					_ = ctx.Put("_dur", Const(int64(nVal)))
					return fn.Eval(bit, ctx)
				}
			} else if nList, ok := val.(List); ok {
				var loop int64
				for _, item := range nList {
					if it, ok := item.(Num); ok {
						loop += int64(it)
					} else {
						panic(fmt.Errorf("Cannot use item as argument for Every: %v (%v)", item, item))
					}
				}
				x := bit % loop
				var s int64
				for _, item := range nList {
					cur := int64(item.(Num))
					if x == s {
						_ = ctx.Put("_dur", Const(cur))
						return fn.Eval(bit, ctx)
					}
					s += cur
				}
			} else {
				panic(fmt.Errorf("Unknown type: %v (%T)", val, val))
			}
			return nil
		}
		return types.SignalFn(f)
	}
}
