package common

import (
	"fmt"
	"math/rand"

	"github.com/avoronkov/waver/lib/seq/types"
)

var (
	rndSeeded = false
	rnd       = rand.New(rand.NewSource(1))
)

func Srand(seed int64) {
	if rndSeeded {
		// Seed random only only
		return
	}
	rndSeeded = true
	rnd = rand.New(rand.NewSource(seed))
}

func Random(values types.ValueFn) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		vals := values.Val(bit, ctx)
		list, ok := vals.(EvaluatedList)
		if !ok {
			panic(fmt.Errorf("rand expects list, found: %v", vals))
		}
		i := rnd.Intn(list.Len())
		return list.Get(i)
	}
	return types.ValueFunc(f)
}
