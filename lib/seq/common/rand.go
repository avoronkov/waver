package common

import (
	"math/rand"

	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func Random(values ...types.ValueFn) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		i := rand.Intn(len(values))
		return values[i].Val(bit, ctx)
	}
	return types.ValueFunc(f)
}
