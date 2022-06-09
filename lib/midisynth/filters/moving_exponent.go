package filters

import (
	"fmt"
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

func (MovingExponent) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o []func(*MovingExponent)
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "initialValue":
				val := float64Of(value)
				o = append(o, MovingExponentInitialValue(val))
			case "speed":
				val := float64Of(value)
				o = append(o, MovingExponentSpeed(val))
			case "inverse":
				val := value.(bool)
				o = append(o, MovingExponentInverse(val))
			default:
				return nil, fmt.Errorf("Unknon Moving Exponent parameter: %v", param)
			}
		}
	}
	return NewMovingExponent(o...), nil
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
