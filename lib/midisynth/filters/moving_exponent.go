package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type MovingExponent struct {
	InitialValue float64 `option:"initialValue"`
	Speed        float64 `option:"speed"`
	Inverse      bool    `option:"inverse"`
}

func (MovingExponent) New() Filter {
	return &MovingExponent{
		InitialValue: 1.0,
		Speed:        0.25,
	}
}

func (e *MovingExponent) Apply(input waves.Wave) waves.Wave {
	return MakeFilterImpl(e, input, movingExpValue)
}

func movingExpValue(fx *MovingExponent, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	v := input.Value(t, ctx)
	e := fx.InitialValue + fx.Speed*t
	if fx.Inverse {
		e = 1.0 / e
	}
	if v < 0.0 {
		return -math.Pow(-v, e)
	}
	return math.Pow(v, e)
}
