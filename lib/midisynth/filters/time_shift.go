package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type TimeShift struct {
	Carrier   waves.Wave `option:"carrier"`
	Frequency float64    `option:"frequency,freq"`
	Amplitude float64    `option:"amplitude,amp"`
}

func (TimeShift) New() Filter {
	return &TimeShift{
		Carrier:   waves.Sine,
		Frequency: 4.0,
		Amplitude: 0.05,
	}
}

func (TimeShift) Desc() string {
	return `Experimental time shift effect.`
}

func (ts *TimeShift) Apply(w waves.Wave) waves.Wave {
	return &timeShiftImpl{
		input:      w,
		opts:       ts,
		shifterCtx: waves.NewNoteCtx(ts.Frequency, ts.Amplitude, math.Inf(1), 0.0),
	}
}

type timeShiftImpl struct {
	input      waves.Wave
	opts       *TimeShift
	shifterCtx *waves.NoteCtx
}

func (i *timeShiftImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	newFreq := (i.opts.Carrier.Value(t, i.shifterCtx)*i.shifterCtx.Amp + 1.0) * ctx.Freq
	newCtx := waves.NewNoteCtx(newFreq, ctx.Amp, ctx.Dur, ctx.AbsTime)
	return i.input.Value(t, newCtx)
}
