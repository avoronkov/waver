package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type SwingExp struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"frequency,freq"`
	Amplitude float64    `option:"amplitude,amp"`
	Inverse   bool       `option:"inverse"`

	carrierCtx *waves.NoteCtx
}

var cos = waves.WaveFn(func(t float64, ctx *waves.NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	return math.Cos(x)
})

func (SwingExp) New() Filter {
	return &SwingExp{
		Carrier:   cos,
		Frequency: 0.25,
		Amplitude: 1.0,
	}
}

func (SwingExp) Desc() string {
	return `Swinging Exponent is an effect of constantly changing exponent effect using periodic function (cosine).`
}

func (s *SwingExp) Apply(input waves.Wave) waves.Wave {
	s.carrierCtx = waves.NewNoteCtx(s.Frequency, s.Amplitude, -1.0, 0.0)
	return MakeFilterImpl(s, input, swingExpImpl)
}

func swingExpImpl(fx *SwingExp, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	c := fx.Carrier.Value(ctx.AbsTime, fx.carrierCtx)
	e := (1.0-c)*fx.carrierCtx.Amp/2.0 + 1.0
	if fx.Inverse {
		e = 1.0 / e
	}
	v := input.Value(t, ctx)
	if v <= 0 {
		return -math.Pow(-v, e)
	}
	return math.Pow(v, e)
}
