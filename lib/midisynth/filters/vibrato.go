package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type VibratoFilter struct {
	shifter    waves.Wave
	shifterCtx *waves.NoteCtx

	freq float64
	amp  float64
}

func NewVibrato(opts ...func(*VibratoFilter)) Filter {
	f := &VibratoFilter{
		shifter: waves.Sine,
	}

	for _, o := range opts {
		o(f)
	}

	f.shifterCtx = waves.NewNoteCtx(f.freq, f.amp, math.Inf(1))

	return f
}

func (v *VibratoFilter) Apply(input waves.Wave) waves.Wave {
	return &vibratoImpl{
		input: input,
		opts:  v,
	}
}

type vibratoImpl struct {
	input waves.Wave
	opts  *VibratoFilter
}

func (i *vibratoImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	shift := i.opts.shifter.Value(t, i.opts.shifterCtx) * i.opts.shifterCtx.Amp
	return i.input.Value(t+shift, ctx)
}

// Options

func VibratoFrequency(v float64) func(f *VibratoFilter) {
	return func(f *VibratoFilter) {
		f.freq = v
	}
}

func VibratoAmplitude(v float64) func(f *VibratoFilter) {
	return func(f *VibratoFilter) {
		f.amp = v
	}
}

func VibratoCarrierWave(w waves.Wave) func(f *VibratoFilter) {
	return func(f *VibratoFilter) {
		f.shifter = w
	}
}
