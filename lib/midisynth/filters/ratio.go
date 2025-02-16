package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Ratio struct {
	Value float64 `option:"value"`
}

func (Ratio) New() Filter {
	return &Ratio{
		Value: 1.0,
	}
}

func (Ratio) Desc() string {
	return `Effect of increasion / decreasing speed of the input signal.`
}

func ratioImpl(fx *Ratio, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	return input.Value(t*fx.Value, ctx)
}

func (r *Ratio) Apply(wave waves.Wave) waves.Wave {
	return MakeFilterImpl(r, wave, ratioImpl)
}
