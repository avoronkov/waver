package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type TimeShift struct {
	period float64
	spread float64
}

func NewTimeShift(hz float64, spread float64) Filter {
	return &TimeShift{
		period: 1.0 / hz,
		spread: spread,
	}
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

func (i *timeShiftImpl) Value(tm float64, ctx *waves.NoteCtx) float64 {
	return i.input.Value(tm+(i.opts.spread*math.Sin(2.0*math.Pi*tm/i.opts.period)), ctx)
}
