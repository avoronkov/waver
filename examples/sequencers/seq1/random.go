package main

import (
	"math/rand"

	"gitlab.com/avoronkov/waver/lib/seq"
)

func RandomNote(notes ...rune) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			i := rand.Intn(len(notes))
			old := ctx["note"]
			ctx["note"] = notes[i]
			res := fn.Eval(bit, ctx)
			ctx["note"] = old
			return res
		}
		return seq.SignalFn(f)
	}
}
