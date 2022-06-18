package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Ratio struct {
	Value float64
}

func (Ratio) New() Filter {
	return &Ratio{
		Value: 1.0,
	}
}

func (r *Ratio) Apply(wave waves.Wave) waves.Wave {
	return &ratioImpl{
		input: wave,
		opts:  r,
	}
}

type ratioImpl struct {
	input waves.Wave
	opts  *Ratio
}

func (i *ratioImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	return i.input.Value(t*i.opts.Value, ctx)
}
