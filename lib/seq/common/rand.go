package common

import (
	"fmt"
	"math/rand"

	"github.com/avoronkov/waver/lib/seq/types"
)

func Random(values types.ValueFn) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		vals := values.Val(bit, ctx)
		list, ok := vals.(EvaluatedList)
		if !ok {
			panic(fmt.Errorf("rand expects list, found: %v", vals))
		}
		i := rand.Intn(list.Len())
		return list.Get(i)
	}
	return types.ValueFunc(f)
}
