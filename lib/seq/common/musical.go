package common

import "github.com/avoronkov/waver/lib/seq/types"

func ChordFn(keyNote types.ValueFn, steps ...int64) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		kn := keyNote.Val(bit, ctx)
		k := int64(kn.(Num))
		// TODO use GreedyEvaluated list
		res := &LazyEvaluatedList{
			bit: bit,
			ctx: ctx,
		}
		for _, step := range steps {
			res.values = append(res.values, Const(k+step))
		}
		return res
	}
	return types.ValueFunc(f)
}
