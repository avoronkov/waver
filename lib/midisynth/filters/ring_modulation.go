package filters

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

// Ring (amplitude) modulation.
type Ring struct {
	Carrier    waves.Wave
	CarrierCtx *waves.NoteCtx
}

func NewRing(carrier waves.Wave, freq, amp float64) Filter {
	return &Ring{
		Carrier:    carrier,
		CarrierCtx: waves.NewNoteCtx(freq, amp, -1.0),
	}
}

func (Ring) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var carrier waves.Wave = waves.Sine
	var freq float64
	amp := 1.0

	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "wave":
				return nil, fmt.Errorf("Parameter 'wave' is not supported yet")
			case "freq", "frequency":
				freq = float64Of(value)
			case "amp", "amplitude":
				amp = float64Of(value)
			default:
				return nil, fmt.Errorf("Unknown AM parameter: %v", param)
			}
		}
	}
	return NewRing(carrier, freq, amp), nil
}

func (rf *Ring) Apply(input waves.Wave) waves.Wave {
	return &ringImpl{
		input: input,
		opts:  rf,
	}
}

type ringImpl struct {
	input waves.Wave
	opts  *Ring
}

func (i *ringImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	mp := (2.0 - i.opts.CarrierCtx.Amp + i.opts.CarrierCtx.Amp*i.opts.Carrier.Value(t, i.opts.CarrierCtx)) / 2.0
	return i.input.Value(t, ctx) * mp
}
