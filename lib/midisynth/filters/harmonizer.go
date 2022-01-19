package filters

import "gitlab.com/avoronkov/waver/lib/midisynth/waves"

type Harmonizer struct {
	harmonics map[int]float64
}

func NewHarmonizer(opts ...func(*Harmonizer)) Filter {
	h := &Harmonizer{
		harmonics: map[int]float64{
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
	for _, s := range i.opts.harmonics {
		total += s
	}
	v := 0.0
	for n, s := range i.opts.harmonics {
		v += s * i.wave.Value(t, waves.NewNoteCtx(ctx.Freq*float64(n), ctx.Amp, ctx.Dur))
	}
	return v / total
}

// Options
func Harmonic(n int, strength float64) func(*Harmonizer) {
	return func(h *Harmonizer) {
		h.harmonics[n] = strength
	}
}
