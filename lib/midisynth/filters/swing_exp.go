package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type SwingExp struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"freq,frequency"`
	Amplitude float64    `option:"amp,amplitude"`
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

func (s *SwingExp) Apply(input waves.Wave) waves.Wave {
	return &swingExpImpl{
		wave:       input,
		opts:       s,
		carrierCtx: waves.NewNoteCtx(s.Frequency, s.Amplitude, -1.0, 0.0),
	}
}

type swingExpImpl struct {
	wave       waves.Wave
	opts       *SwingExp
	carrierCtx *waves.NoteCtx
}

func (i *swingExpImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	c := i.opts.Carrier.Value(ctx.AbsTime, i.carrierCtx)
	e := (1.0-c)*i.carrierCtx.Amp/2.0 + 1.0
	v := i.wave.Value(t, ctx)
	if v <= 0 {
		return -math.Pow(-v, e)
	}
	return math.Pow(v, e)
}
