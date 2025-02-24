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
			bitsList, ok := bitsVal.(EvaluatedList)
			if !ok {
				panic(fmt.Errorf("Cannot use non-list item as argument for Bits: %v (%T)", bitsVal, bitsVal))
			}
			l := bitsList.Len()
			totalBits := int64(bitsList.Get(l - 1).(Num))
			currentBit := bit % totalBits
			for i := range l - 1 {
				it := int64(bitsList.Get(i).(Num))
				if it == currentBit {
					dur := int64(bitsList.Get(i+1).(Num)) - it
					_ = ctx.Put("_dur", Const(dur))
					return fn.Eval(bit, ctx)
				}
			}
			return nil
		}
		return types.SignalFn(f)
	}
}
