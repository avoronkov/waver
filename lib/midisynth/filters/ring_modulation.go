package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

// Ring (amplitude) modulation.
type Ring struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"frequency,freq"`
	Amplitude float64    `option:"amplitude,amp"`
}

func (Ring) New() Filter {
	return &Ring{
		Carrier:   waves.Sine,
		Frequency: 4.0,
		Amplitude: 1.0,
	}
}

func (Ring) Desc() string {
	return `Amplitude modulation.`
}

func (rf *Ring) Apply(input waves.Wave) waves.Wave {
	return &ringImpl{
		input:      input,
		opts:       rf,
		carrierCtx: waves.NewNoteCtx(rf.Frequency, rf.Amplitude, -1.0, 0.0),
	}
}

type ringImpl struct {
	input      waves.Wave
	opts       *Ring
	carrierCtx *waves.NoteCtx
}

func (i *ringImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	mp := (2.0 - i.carrierCtx.Amp + i.carrierCtx.Amp*i.opts.Carrier.Value(t, i.carrierCtx)) / 2.0
	return i.input.Value(t, ctx) * mp
}
