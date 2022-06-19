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
	return MakeFilterImpl(ef, input, expImplFn)
}

func expImplFn(fx *Exponent, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	v := input.Value(t, ctx)
	if v < 0.0 {
		return -math.Pow(-v, fx.Value)
	}
	return math.Pow(v, fx.Value)
}
