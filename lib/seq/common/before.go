package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/types"
)

func Before(n types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			val := n.Val(bit, ctx)
			nVal, ok := val.(Num)
			if !ok {
				panic(fmt.Errorf("Cannot use item as argument for Before(`<`): %v (%T)", val, val))
			}
			if bit >= int64(nVal) {
				return nil
			}
			return fn.Eval(bit, ctx)
		}
		return types.SignalFn(f)
	}
}
