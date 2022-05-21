package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

// Number
type Num int64

var _ types.Value = Num(0)

func (n Num) IsValue() {}

// List of values
type List []types.Value

var _ types.Value = List(nil)

func (l List) IsValue() {}

// Float
type Float float64

var _ types.Value = Float(0.0)

func (f Float) IsValue() {}
