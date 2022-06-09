package filters

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

// Amplitude distortion filter
type DistortionFilter struct {
	Multiplier float64
}

func NewDistortionFilter(m float64) *DistortionFilter {
	return &DistortionFilter{
		Multiplier: m,
	}
}

func (f DistortionFilter) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	value := 1.0
	if options != nil {
		opts := options.(map[string]any)
		for param, v := range opts {
			switch param {
			case "value":
				value = float64Of(v)
			default:
				return nil, fmt.Errorf("Unknown Distortion parameter: %v", param)
			}
		}
	}
	return NewDistortionFilter(value), nil
}

func (df *DistortionFilter) Apply(w waves.Wave) waves.Wave {
	return &distortionImpl{
		wave: w,
		opts: df,
	}
}

type distortionImpl struct {
	wave waves.Wave
	opts *DistortionFilter
}

func (i *distortionImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	val := i.wave.Value(t, ctx) * i.opts.Multiplier
	if val > 1.0 {
		val = 1.0
	} else if val < -1.0 {
		val = -1.0
	}
	return val
}
