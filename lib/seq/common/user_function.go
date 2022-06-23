package common

import (
	"github.com/avoronkov/waver/lib/seq/types"
)

func UserFunction(argName string, arg, body types.ValueFn) types.ValueFn {
	f := func(bit int64, ctx types.Context) types.Value {
		newContext := ctx.Copy()
		_ = newContext.Put(argName, arg)
		return body.Val(bit, newContext)
	}
	return types.ValueFunc(f)
}
