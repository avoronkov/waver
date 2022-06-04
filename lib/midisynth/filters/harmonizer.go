package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Harmonizer struct {
	harmonics []float64
}

func NewHarmonizer(opts ...func(*Harmonizer)) Filter {
	h := &Harmonizer{
		harmonics: []float64{
			0: 0.0,
			1: 1.0,
			2: 0.5,
			3: 0.25,
			4: 0.125,
			5: 0.0625,
		},
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Harmonizer) Apply(w waves.Wave) waves.Wave {
	return &harmonizerImpl{
		wave: w,
		opts: h,
	}
}

type harmonizerImpl struct {
	wave waves.Wave
	opts *Harmonizer
}

func (i *harmonizerImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	total := 0.0
	v := 0.0
	for n, s := range i.opts.harmonics[1:] {
		v += s * i.wave.Value(float64(n)*t, ctx)
		total += s
	}
	return v / total
}

// Options
func Harmonic(n int, strength float64) func(*Harmonizer) {
	return func(h *Harmonizer) {
		h.harmonics[n] = strength
	}
}
