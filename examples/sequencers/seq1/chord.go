package main

import "github.com/avoronkov/waver/lib/seq"

func Chord(notes ...rune) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) (res []string) {
			old := ctx["note"]
			for _, n := range notes {
				ctx["note"] = n
				res = append(res, fn.Eval(bit, ctx)...)
			}
			ctx["note"] = old
			return res
		}
		return seq.SignalFn(f)
	}
}
