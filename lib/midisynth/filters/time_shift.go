package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type TimeShift struct {
	shifter    waves.Wave
	shifterCtx *waves.NoteCtx

	freq float64
	amp  float64
}

func NewTimeShift(opts ...func(*TimeShift)) Filter {
	f := &TimeShift{
		shifter: &waves.Sine{},
	}

	for _, o := range opts {
		o(f)
	}

	f.shifterCtx = waves.NewNoteCtx(f.freq, f.amp, math.Inf(1))

	return f
}

func (ts *TimeShift) Apply(w waves.Wave) waves.Wave {
	return &timeShiftImpl{
		input: w,
		opts:  ts,
	}
}

type timeShiftImpl struct {
	input waves.Wave
	opts  *TimeShift
}

func (i *timeShiftImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	newFreq := (i.opts.shifter.Value(t, i.opts.shifterCtx)*i.opts.shifterCtx.Amp + 1.0) * ctx.Freq
	newCtx := waves.NewNoteCtx(newFreq, ctx.Amp, ctx.Dur)
	res := i.input.Value(t, newCtx)
	return res
}

// Options

func TimeShiftFrequency(v float64) func(f *TimeShift) {
	return func(f *TimeShift) {
		f.freq = v
	}
}

func TimeShiftAmplitude(v float64) func(f *TimeShift) {
	return func(f *TimeShift) {
		f.amp = v
	}
}

func TimeShiftCarrierWave(w waves.Wave) func(f *TimeShift) {
	return func(f *TimeShift) {
		f.shifter = w
	}
}
