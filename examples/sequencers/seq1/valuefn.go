package main

import (
	"math/rand"

	"github.com/avoronkov/waver/lib/seq"
)

type ValueFn interface {
	Val(bit int64, ctx seq.Context) Value
}

type ValueFunc func(bit int64, ctx seq.Context) Value

func (f ValueFunc) Val(bit int64, ctx seq.Context) Value {
	return f(bit, ctx)
}

func Const(n int64) ValueFn {
	return ValueFunc(func(int64, seq.Context) Value {
		return Num(n)
	})
}

func Lst(values ...int64) ValueFn {
	return ValueFunc(func(int64, seq.Context) Value {
		var res List
		for _, v := range values {
			res = append(res, Num(v))
		}
		return res
	})
}

func Var(name string) ValueFn {
	return ValueFunc(func(n int64, ctx seq.Context) Value {
		return ctx[name].(Value)
	})
}

//

func Random(vfns ...ValueFn) ValueFn {
	return ValueFunc(func(bit int64, ctx seq.Context) Value {
		i := rand.Intn(len(vfns))
		return vfns[i].Val(bit, ctx)
	})
}

type sequence struct {
	fns []ValueFn
	idx int
}

func (s *sequence) Val(bit int64, ctx seq.Context) (res Value) {
	res = s.fns[s.idx].Val(bit, ctx)
	s.idx = (s.idx + 1) % len(s.fns)
	return
}

func Sequence(vfns ...ValueFn) ValueFn {
	return &sequence{
		fns: vfns,
	}
}
