package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

type Num int64

var _ types.Value = Num(0)

func (n Num) IsValue() {}

type List []types.Value

var _ types.Value = List(nil)

func (l List) IsValue() {}
