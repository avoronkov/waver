package main

import (
	"gitlab.com/avoronkov/waver/lib/seq"
)

type Modifier = func(seq.Signaler) seq.Signaler

func Every(n int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			if bit%n == 0 {
				return fn.Eval(bit, ctx)
			}
			return nil
		}
		return seq.SignalFn(f)
	}
}

func Shift(n int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			return fn.Eval(bit-n, ctx)
		}
		return seq.SignalFn(f)
	}
}

func Sig(signal string) seq.SignalFn {
	return func(bit int64, ctx seq.Context) []string {
		return []string{signal}
	}
}

func OnBits(loop int64, bits ...int64) Modifier {
	return func(fn seq.Signaler) seq.Signaler {
		f := func(bit int64, ctx seq.Context) []string {
			mod := bit % loop
			for _, b := range bits {
				if mod == b {
					return fn.Eval(bit, ctx)
				}
			}
			return nil
		}
		return seq.SignalFn(f)
	}
}

func Chain(fn seq.Signaler, modifiers ...Modifier) seq.Signaler {
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

	_ = Chain(
		Note(NoteInstr(1), NoteOctave(3), NoteAmp(1)),
		Chord('A', 'C', 'E', 'G'),
		Every(7),
	)

	melody := Chain(
		Note(NoteInstr(2), NoteOctave(4), NoteAmp(1)),
		RandomNote('C', 'A', 'F', 'E', 'D'),
		Every(3),
	)

	seq.Run(
		kick,
		hat,
		snare,
		// harmony,
		melody,
	)
}
