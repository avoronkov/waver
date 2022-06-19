package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type VibratoFilter struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"freq,frequency"`
	Amplitude float64    `option:"amp,amplitude"`
}

func (VibratoFilter) New() Filter {
	return &VibratoFilter{
		Carrier:   waves.Sine,
		Frequency: 1.0,
		Amplitude: 0.5,
	}
}

func (v *VibratoFilter) Apply(input waves.Wave) waves.Wave {
	return &vibratoImpl{
		input:      input,
		opts:       v,
		shifterCtx: waves.NewNoteCtx(v.Frequency, v.Amplitude, math.Inf(1), 0.0),
	}
}

type vibratoImpl struct {
	input      waves.Wave
	opts       *VibratoFilter
	shifterCtx *waves.NoteCtx
}

func (i *vibratoImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	shift := i.opts.Carrier.Value(t, i.shifterCtx) * i.opts.Amplitude
	return i.input.Value(t+shift, ctx)
}
