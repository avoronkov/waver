package common

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/types"
)

func Concat(values types.ValueFn) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		vals := values.Val(bit, ctx)
		list, ok := vals.(EvaluatedList)
		if !ok {
			panic(fmt.Errorf("Concat expects list, found: %v", vals))
		}
		res := &GreedyEvaluatedList{}
		l := list.Len()
		for i := 0; i < l; i++ {
			item := list.Get(i)
			itemList, ok := item.(EvaluatedList)
			if !ok {
				panic(fmt.Errorf("Concat: element of list is not a list: %v (%T)", item, item))
			}
			ilen := itemList.Len()
			for j := 0; j < ilen; j++ {
				x := itemList.Get(j)
				res.values = append(res.values, x)
			}
		}
		return res
	}
	return types.ValueFunc(f)
}
