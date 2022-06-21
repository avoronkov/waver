package common

import "github.com/avoronkov/waver/lib/seq/types"

// Number
type Num int64

var _ types.Value = Num(0)

func (n Num) IsValue() {}

// Float
type Float float64

var _ types.Value = Float(0.0)

func (f Float) IsValue() {}

// Str
type Str string

var _ types.Value = Str("")

func (Str) IsValue() {}

// List of values
type List []types.ValueFn

var _ types.ValueFn = List(nil)

func (l List) Val(bit int64, ctx types.Context) types.Value {
	return EvaluatedList{
		values: []types.ValueFn(l),
		bit:    bit,
		ctx:    ctx,
	}
}

// Evaluated list of values
type EvaluatedList struct {
	values []types.ValueFn
	bit    int64
	ctx    types.Context
}

func (l EvaluatedList) IsValue() {}

func (l EvaluatedList) Len() int {
	return len(l.values)
}

func (l EvaluatedList) Get(i int) types.Value {
	return l.values[i].Val(l.bit, l.ctx)
}
