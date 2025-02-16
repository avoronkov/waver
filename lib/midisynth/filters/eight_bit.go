package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type EightBit struct {
	Bits int `option:"bits"`
}

func (EightBit) New() Filter {
	return &EightBit{
		Bits: 8,
	}
}

func (EightBit) Desc() string {
	return `8-bit sound.`
}

func (ef *EightBit) Apply(input waves.Wave) waves.Wave {
	return &eightBitImpl{
		input: input,
		opts:  ef,
	}
}

type eightBitImpl struct {
	input waves.Wave
	opts  *EightBit
}

func (i *eightBitImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.input.Value(t, ctx)
	n := math.Pow(2.0, float64(i.opts.Bits/2))
	return math.Ceil(v*n) / n
}
