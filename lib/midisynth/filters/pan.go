package filters

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Pan struct {
	channels []float64
}

func NewPan(opts ...func(*Pan)) Filter {
	p := &Pan{
		channels: []float64{1.0, 1.0},
	}

	for _, o := range opts {
		o(p)
	}

	return p
}

func (Pan) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o []func(*Pan)
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "l", "left":
				v := float64Of(value)
				o = append(o, PanLeft(v))
			case "r", "right":
				v := float64Of(value)
				o = append(o, PanRight(v))
			default:
				return nil, fmt.Errorf("Unknown Pan parameter: %v", param)
			}
		}
	}
	return NewPan(o...), nil
}

func (p *Pan) Apply(w waves.Wave) waves.Wave {
	return &panImpl{
		wave: w,
		opts: p,
	}
}

type panImpl struct {
	wave waves.Wave
	opts *Pan
}

func (i *panImpl) Value(tm float64, ctx *waves.NoteCtx) (res float64) {
	return i.wave.Value(tm, ctx) * i.opts.channels[ctx.Channel]
}

// Options

func PanLeft(v float64) func(*Pan) {
	return func(p *Pan) {
		p.channels[0] = v
	}
}

func PanRight(v float64) func(*Pan) {
	return func(p *Pan) {
		p.channels[1] = v
	}
}
