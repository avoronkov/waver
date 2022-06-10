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

var _ types.Value = List(nil)

func (l List) IsValue() {}

func (l List) Len() int {
	return len(l)
}

func (l List) Get(i int, bit int64, ctx types.Context) types.Value {
	return l[i].Val(bit, ctx)
}
