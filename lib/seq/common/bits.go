package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/types"
)

func Bits(bits types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			bitsVal := bits.Val(bit, ctx)
			bitsList, ok := bitsVal.(List)
			if !ok {
				panic(fmt.Errorf("Cannot use non-list item as argument for Bits: %v (%T)", bitsVal, bitsVal))
			}
			l := bitsList.Len()
			totalBits := int64(bitsList.Get(l-1, bit, ctx).(Num))
			currentBit := bit % totalBits
			for i := 0; i < l-1; i++ {
				it := int64(bitsList.Get(i, bit, ctx).(Num))
				if it == currentBit {
					dur := int64(bitsList.Get(i+1, bit, ctx).(Num)) - it
					_ = ctx.Put("_dur", Const(dur))
					return fn.Eval(bit, ctx)
				}
			}
			return nil
		}
		return types.SignalFn(f)
	}
}
