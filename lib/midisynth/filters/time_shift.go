package filters

import (
	"fmt"
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type TimeShift struct {
	shifter    waves.Wave
	shifterCtx *waves.NoteCtx

	freq float64
	amp  float64
}

func NewTimeShift(opts ...func(*TimeShift)) Filter {
	f := &TimeShift{
		shifter: waves.Sine,
	}

	for _, o := range opts {
		o(f)
	}

	f.shifterCtx = waves.NewNoteCtx(f.freq, f.amp, math.Inf(1), 0.0)

	return f
}

func (TimeShift) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o []func(*TimeShift)
	o = append(o, TimeShiftCarrierWave(waves.Sine))
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "wave":
				return nil, fmt.Errorf("Parameter 'wave' is not supported yet")
			case "freq", "frequency":
				o = append(o, TimeShiftFrequency(float64Of(value)))
			case "amp", "amplitude":
				o = append(o, TimeShiftAmplitude(float64Of(value)))
			default:
				return nil, fmt.Errorf("Unknown Time Shift parameter: %v", param)
			}
		}
	}
	return NewTimeShift(o...), nil
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
	newCtx := waves.NewNoteCtx(newFreq, ctx.Amp, ctx.Dur, ctx.AbsTime)
	return i.input.Value(t, newCtx)
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
