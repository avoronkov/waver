package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/types"
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
			} else if nList, ok := val.(EvaluatedList); ok {
				var loop int64
				l := nList.Len()
				for i := range l {
					item := nList.Get(i)
					if it, ok := item.(Num); ok {
						loop += int64(it)
					} else {
						panic(fmt.Errorf("Cannot use item as argument for Every: %v (%T)", item, item))
					}
				}
				x := bit % loop
				var s int64
				for i := range l {
					item := nList.Get(i)
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
