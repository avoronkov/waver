package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Exponent struct {
	Value float64
}

func (Exponent) New() Filter {
	return &Exponent{
		Value: 1.0,
	}
}

func (ef *Exponent) Apply(input waves.Wave) waves.Wave {
	return &expImpl{
		input: input,
		opts:  ef,
	}
}

type expImpl struct {
	input waves.Wave
	opts  *Exponent
}

func (i *expImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.input.Value(t, ctx)
	if v < 0.0 {
		return -math.Pow(-v, i.opts.Value)
	}
	return math.Pow(v, i.opts.Value)
}
