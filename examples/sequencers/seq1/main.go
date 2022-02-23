package main

import (
	"gitlab.com/avoronkov/waver/lib/seq"
)

type Modifier = func(seq.SignalFn) seq.SignalFn

func Every(n int64) func(seq.SignalFn) seq.SignalFn {
	return func(fn seq.SignalFn) seq.SignalFn {
		return func(bit int64) []string {
			if bit%n == 0 {
				return fn(bit)
			}
			return nil
		}
	}
}

func Shift(n int64) func(seq.SignalFn) seq.SignalFn {
	return func(fn seq.SignalFn) seq.SignalFn {
		return func(bit int64) []string {
			return fn(bit - n)
		}
	}
}

func Sig(signal string) seq.SignalFn {
	return func(bit int64) []string {
		return []string{signal}
	}
}

func OnBits(loop int64, bits ...int64) Modifier {
	return func(fn seq.SignalFn) seq.SignalFn {
		return func(bit int64) []string {
			mod := bit % loop
			for _, b := range bits {
				if mod == b {
					return fn(bit)
				}
			}
			return nil
		}
	}
}

func Chain(fn seq.SignalFn, modifiers ...func(seq.SignalFn) seq.SignalFn) seq.SignalFn {
	res := fn
	for _, md := range modifiers {
		res = md(res)
	}
	return res
}

func main() {
	kick := Chain(Sig("z4k2"), Every(8))
	hat := Chain(Sig("z4h2"), Every(8), Shift(4))
	snare := Chain(Sig("z4s2"), OnBits(16, 3, 6, 9, 12), Shift(-1))

	seq.Run(
		kick,
		hat,
		snare,
	)
}
