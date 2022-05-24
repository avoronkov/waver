package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Flanger struct {
	shifter    waves.Wave
	shifterCtx *waves.NoteCtx

	freq     float64
	maxShift float64
}

func NewFlanger(opts ...func(*Flanger)) Filter {
	f := &Flanger{
		shifter: waves.WaveFn(func(t float64, ctx *waves.NoteCtx) float64 {
			x := 2.0 * math.Pi * t / ctx.Period
			return math.Cos(x)
		}),
		maxShift: 0.02,
	}
	for _, opt := range opts {
		opt(f)
	}
	f.shifterCtx = waves.NewNoteCtx(f.freq, 1.0, math.Inf(1))
	return f
}

func (f *Flanger) Apply(wave waves.Wave) waves.Wave {
	return &flangerImpl{
		wave: wave,
		opts: f,
	}
}

type flangerImpl struct {
	wave waves.Wave
	opts *Flanger
}

func (i *flangerImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.wave.Value(t, ctx)
	shift := i.opts.shifter.Value(t, i.opts.shifterCtx) * i.opts.maxShift
	t1 := t - 0.5*shift + 0.5*i.opts.maxShift
	v1 := i.wave.Value(t1, ctx)
	return 0.5*v + 0.5*v1
}

// Options
func FlangerFreq(freq float64) func(f *Flanger) {
	return func(f *Flanger) {
		f.freq = freq
	}
}

// Maximum value of shift (in seconds)
func FlangerShift(shift float64) func(f *Flanger) {
	return func(f *Flanger) {
		f.maxShift = shift
	}
}
