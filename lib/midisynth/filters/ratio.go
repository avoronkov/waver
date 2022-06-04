package filters

import "github.com/avoronkov/waver/lib/midisynth/waves"

type Ratio struct {
	value float64
}

func NewRatio(value float64) Filter {
	return &Ratio{
		value: value,
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
	return i.input.Value(t*i.opts.value, ctx)
}
