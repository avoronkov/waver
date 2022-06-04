package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Exponent struct {
	e float64
}

func NewExponent(value float64) Filter {
	return &Exponent{
		e: value,
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
		return -math.Pow(-v, i.opts.e)
	}
	return math.Pow(v, i.opts.e)
}
