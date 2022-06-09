package filters

import (
	"fmt"
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Exponent struct {
	e float64
}

func NewExponent(value float64) Filter {
	return &Exponent{
		e: value,
	}
}

func (Exponent) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	val := 1.0
	if x, ok := options.(float64); ok {
		val = x
	} else {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "value":
				val = float64Of(value)
			default:
				return nil, fmt.Errorf("Unknown Exponent parameter: %v", param)
			}
		}
	}
	return NewExponent(val), nil
}

func (ef *Exponent) Apply(input waves.Wave) waves.Wave {
	return &expImpl{
		input: input,
		opts:  ef,
	}
}

type expImpl struct {
	input waves.Wave
	opts  *Exponent
}

func (i *expImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.input.Value(t, ctx)
	if v < 0.0 {
		return -math.Pow(-v, i.opts.e)
	}
	return math.Pow(v, i.opts.e)
}
