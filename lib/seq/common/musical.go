package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

func ChordFn(keyNote types.ValueFn, steps ...int64) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		kn := keyNote.Val(bit, ctx)
		k := int64(kn.(Num))
		res := List{}
		for _, step := range steps {
			res = append(res, Num(k+step))
		}
		return res
	}
	return types.ValueFunc(f)
}
