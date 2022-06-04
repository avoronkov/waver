package main

import (
	"github.com/avoronkov/waver/lib/seq"
)

func main() {
	kick := Chain(Sig("z4k2"), Every(8))
	hat := Chain(Sig("z4h2"), Every(8), Shift(4))
	_ = hat
	snare := Chain(Sig("z4s2"), OnBits(16, 3, 6, 9, 12), Shift(-1))
	_ = snare

	harmony := Chain(
		Note(Const(1), Lst(int64(A4), int64(C4), int64(E4), int64(G4)), NoteAmp(Const(1))),
		Every(7),
	)

	melody := Chain(
		Note(
			Const(2),
			Random(Const(int64(C2)), Const(int64(A2)), Const(int64(F2)), Const(int64(E2))),
			NoteAmp(Const(1)),
		),
		Every(3),
	)
	_ = melody

	seq.Run(
		kick,
		hat,
		snare,
		harmony,
		melody,
	)
}
