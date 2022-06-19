package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type SwingPan struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"freq,frequency"`
	Amplitude float64    `option:"amp,amplitude"`
}

func (SwingPan) New() Filter {
	return &SwingPan{
		Carrier:   waves.Sine,
		Frequency: 0.5,
		Amplitude: 1.0,
	}
}

func (p *SwingPan) Apply(input waves.Wave) waves.Wave {
	return &movingPanImpl{
		input:      input,
		opts:       p,
		carrierCtx: waves.NewNoteCtx(p.Frequency, p.Amplitude, -1.0, 0.0),
	}
}

type movingPanImpl struct {
	input      waves.Wave
	opts       *SwingPan
	carrierCtx *waves.NoteCtx
}

func (i *movingPanImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	pan := i.opts.Carrier.Value(ctx.AbsTime, i.carrierCtx) * i.carrierCtx.Amp

	r := (pan + 1.0) / 2.0
	l := 1.0 - r
	v := i.input.Value(t, ctx)
	if ctx.Channel == 0 {
		return v * l
	}
	return v * r
}
