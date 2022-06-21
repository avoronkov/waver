package common

import (
	"github.com/avoronkov/waver/lib/seq/types"
)

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
	return &LazyEvaluatedList{
		values: []types.ValueFn(l),
		bit:    bit,
		ctx:    ctx,
	}
}

// Evaluated list of values

type EvaluatedList interface {
	IsValue()
	Len() int
	Get(i int) types.Value
}

var _ EvaluatedList = (*LazyEvaluatedList)(nil)

type LazyEvaluatedList struct {
	values []types.ValueFn
	bit    int64
	ctx    types.Context
}

func (l *LazyEvaluatedList) IsValue() {}

func (l *LazyEvaluatedList) Len() int {
	return len(l.values)
}

func (l *LazyEvaluatedList) Get(i int) types.Value {
	return l.values[i].Val(l.bit, l.ctx)
}

var _ EvaluatedList = (*GreedyEvaluatedList)(nil)

type GreedyEvaluatedList struct {
	values []types.Value
}

func (l *GreedyEvaluatedList) IsValue() {}

func (l *GreedyEvaluatedList) Len() int {
	return len(l.values)
}

func (l *GreedyEvaluatedList) Get(i int) types.Value {
	return l.values[i]
}
