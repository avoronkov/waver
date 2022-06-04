package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type MovingExponent struct {
	initialValue float64
	speed        float64
	inverse      bool
}

func NewMovingExponent(opts ...func(*MovingExponent)) Filter {
	f := &MovingExponent{
		initialValue: 1.0,
	}
	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (e *MovingExponent) Apply(input waves.Wave) waves.Wave {
	return &movingExpImpl{
		input: input,
		opts:  e,
	}
}

type movingExpImpl struct {
	input waves.Wave
	opts  *MovingExponent
}

func (i *movingExpImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.input.Value(t, ctx)
	e := i.opts.initialValue + i.opts.speed*t
	if i.opts.inverse {
		e = 1.0 / e
	}
	if v < 0.0 {
		return -math.Pow(-v, e)
	}
	return math.Pow(v, e)
}

// Options
func MovingExponentInitialValue(v float64) func(f *MovingExponent) {
	return func(f *MovingExponent) {
		f.initialValue = v
	}
}

func MovingExponentSpeed(s float64) func(f *MovingExponent) {
	return func(f *MovingExponent) {
		f.speed = s
	}
}

func MovingExponentInverse(i bool) func(f *MovingExponent) {
	return func(f *MovingExponent) {
		f.inverse = i
	}
}
