package filters

import (
	"fmt"
	"log"
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type SwingExp struct {
	Carrier    waves.Wave
	CarrierCtx *waves.NoteCtx
}

var cos = waves.WaveFn(func(t float64, ctx *waves.NoteCtx) float64 {
	x := 2.0 * math.Pi * t / ctx.Period
	return math.Cos(x)
})

func NewSwingExp(opts ...func(*SwingExp)) Filter {
	f := &SwingExp{
		Carrier:    cos,
		CarrierCtx: waves.NewNoteCtx(0.2, 1.0, -1.0, 0.0),
	}
	for _, o := range opts {
		o(f)
	}
	return f
}

func (SwingExp) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	var o []func(*SwingExp)
	if options != nil {
		opts := options.(map[string]any)
		for param, value := range opts {
			switch param {
			case "amp", "amplitude":
				v := float64Of(value)
				o = append(o, SwingExpAmp(v))
			case "freq", "frequency":
				v := float64Of(value)
				o = append(o, SwingExpFreq(v))
			default:
				return nil, fmt.Errorf("Unknown SwingExp parameter: %v", param)
			}
		}
	}
	return NewSwingExp(o...), nil
}

func (s *SwingExp) Apply(input waves.Wave) waves.Wave {
	return &swingExpImpl{
		wave: input,
		opts: s,
	}
}

type swingExpImpl struct {
	wave waves.Wave
	opts *SwingExp
	n    int64
}

func (i *swingExpImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	c := i.opts.Carrier.Value(ctx.AbsTime, i.opts.CarrierCtx)
	e := (1.0-c)*i.opts.CarrierCtx.Amp/2.0 + 1.0
	// e := 1.0
	if i.n == 0 {
		log.Printf("t: %v, c: %v e: %v", ctx.AbsTime, c, e)
	}
	i.n++
	if i.n == 44100 {
		i.n = 0
	}
	v := i.wave.Value(t, ctx)
	if v <= 0 {
		return -math.Pow(-v, e)
	}
	return math.Pow(v, e)
}

// Options
func SwingExpAmp(amp float64) func(*SwingExp) {
	return func(f *SwingExp) {
		f.CarrierCtx.Amp = amp
	}
}

func SwingExpFreq(freq float64) func(*SwingExp) {
	return func(f *SwingExp) {
		f.CarrierCtx.SetFrequency(freq)
	}
}
