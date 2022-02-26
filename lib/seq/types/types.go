package types

import "gitlab.com/avoronkov/waver/lib/midisynth/signals"

type Context = map[string]interface{}

// Signaler
type Signaler interface {
	Eval(bit int64, ctx Context) []signals.Signal
}

type SignalFn func(bit int64, ctx Context) []signals.Signal

func (f SignalFn) Eval(bit int64, ctx Context) []signals.Signal {
	return f(bit, ctx)
}

// Modifier
type Modifier = func(Signaler) Signaler

// Value function

type Value interface {
	ToInt64List() []int64
}

type ValueFn interface {
	Val(bit int64, ctx Context) Value
}

type ValueFunc func(bit int64, ctx Context) Value

func (f ValueFunc) Val(bit int64, ctx Context) Value {
	return f(bit, ctx)
}
