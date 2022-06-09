package filters

import (
	"fmt"
	"sort"
	"strconv"

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

func (Harmonizer) Create(options any) (fx Filter, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	opts := options.(map[string]any)
	keys := make([]string, 0, len(opts))
	for param := range opts {
		keys = append(keys, param)
	}
	sort.Strings(keys)
	var o []func(*Harmonizer)
	for _, param := range keys {
		n, err := strconv.Atoi(param)
		if err != nil {
			return nil, fmt.Errorf("Incorrect Harmonizer param: %v", param)
		}
		value := opts[param]
		v := float64Of(value)
		o = append(o, Harmonic(n, v))
	}
	return NewHarmonizer(o...), nil
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
