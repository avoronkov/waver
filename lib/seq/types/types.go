package types

import "gitlab.com/avoronkov/waver/lib/midisynth/signals"

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
	IsValue()
}

type ValueFn interface {
	Val(bit int64, ctx Context) Value
}

type ValueFunc func(bit int64, ctx Context) Value

func (f ValueFunc) Val(bit int64, ctx Context) Value {
	return f(bit, ctx)
}

type Context interface {
	Put(name string, fn ValueFn) error
	Get(name string, bit int64) Value
}
