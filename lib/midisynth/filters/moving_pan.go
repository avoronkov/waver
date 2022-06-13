package filters

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type MovingPan struct {
	Carrier    waves.Wave
	CarrierCtx *waves.NoteCtx
}

func NewMovingPan(opts ...func(*MovingPan)) Filter {
	p := &MovingPan{
		Carrier:    waves.Sine,
		CarrierCtx: waves.NewNoteCtx(0.5, 1.0, -1.0, 0.0),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (MovingPan) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o []func(*MovingPan)
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "freq", "frequency":
				v := float64Of(value)
				o = append(o, MovingPanFreq(v))
			case "amp", "amplitude":
				v := float64Of(value)
				o = append(o, MovingPanAmp(v))
			default:
				return nil, fmt.Errorf("Unknown MovingPan parameter: %v", param)
			}
		}
	}
	return NewMovingPan(o...), nil
}

func (p *MovingPan) Apply(input waves.Wave) waves.Wave {
	return &movingPanImpl{
		input: input,
		opts:  p,
	}
}

type movingPanImpl struct {
	input waves.Wave
	opts  *MovingPan
}

func (i *movingPanImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	pan := i.opts.Carrier.Value(ctx.AbsTime, i.opts.CarrierCtx) * i.opts.CarrierCtx.Amp

	r := (pan + 1.0) / 2.0
	l := 1.0 - r
	v := i.input.Value(t, ctx)
	if ctx.Channel == 0 {
		return v * l
	}
	return v * r
}

// Options
func MovingPanFreq(freq float64) func(*MovingPan) {
	return func(p *MovingPan) {
		p.CarrierCtx.SetFrequency(freq)
	}
}

func MovingPanAmp(amp float64) func(*MovingPan) {
	return func(p *MovingPan) {
		p.CarrierCtx.Amp = amp
	}
}
