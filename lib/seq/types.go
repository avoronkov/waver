package seq

type Signal = string

type Context = map[string]interface{}

type Signaler interface {
	Eval(bit int64, ctx Context) []Signal
}

type SignalFn func(bit int64, ctx Context) []Signal

func (f SignalFn) Eval(bit int64, ctx Context) []Signal {
	return f(bit, ctx)
}
