package types

// Signal
type Signal = string

type Context = map[string]interface{}

// Signaler
type Signaler interface {
	Eval(bit int64, ctx Context) []Signal
}

type SignalFn func(bit int64, ctx Context) []Signal

func (f SignalFn) Eval(bit int64, ctx Context) []Signal {
	return f(bit, ctx)
}

// Modifier
type Modifier = func(Signaler) Signaler
