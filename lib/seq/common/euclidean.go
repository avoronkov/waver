package common

import (
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/seq/types"
)

func IsEuclidianBitFirst(pulses, steps, bit int64) (result bool) {
	var i, bucket int64
	for ; i <= bit; i++ {
		if bucket >= 0 {
			bucket -= steps
			result = true
		} else {
			result = false
		}
		bucket += pulses
	}
	return result
}

func IsEuclidianBitLast(pulses, steps, bit int64) (result bool) {
	var i, bucket int64
	for ; i <= bit; i++ {
		bucket += pulses
		if bucket >= steps {
			bucket -= steps
			result = true
		} else {
			result = false
		}
	}
	return result
}

func EuclideanFirst(pulsesVal, stepsVal types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			pulses := int64(pulsesVal.Val(bit, ctx).(Num))
			steps := int64(stepsVal.Val(bit, ctx).(Num))
			step := bit % steps
			if !IsEuclidianBitFirst(pulses, steps, step) {
				return nil
			}
			return fn.Eval(bit, ctx)
		}
		return types.SignalFn(f)
	}
}

func EuclideanLast(pulsesVal, stepsVal types.ValueFn) types.Modifier {
	return func(fn types.Signaler) types.Signaler {
		f := func(bit int64, ctx types.Context) []signals.Signal {
			pulses := int64(pulsesVal.Val(bit, ctx).(Num))
			steps := int64(stepsVal.Val(bit, ctx).(Num))
			step := bit % steps
			if !IsEuclidianBitLast(pulses, steps, step) {
				return nil
			}
			return fn.Eval(bit, ctx)
		}
		return types.SignalFn(f)
	}
}
