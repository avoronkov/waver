package filters

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Ratio struct {
	value float64
}

func NewRatio(value float64) Filter {
	return &Ratio{
		value: value,
	}
}

func (Ratio) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	ratio := 1.0
	if x, ok := options.(float64); ok {
		ratio = x
	} else if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "value":
				val := float64Of(value)
				ratio = val
			default:
				return nil, fmt.Errorf("Unknon Ratio parameter: %v", param)
			}
		}
	}
	return NewRatio(ratio), nil
}

func (r *Ratio) Apply(wave waves.Wave) waves.Wave {
	return &ratioImpl{
		input: wave,
		opts:  r,
	}
}

type ratioImpl struct {
	input waves.Wave
	opts  *Ratio
}

func (i *ratioImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	return i.input.Value(t*i.opts.value, ctx)
}
