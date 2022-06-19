package filters

import (
	"math"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Flanger struct {
	shifter waves.Wave

	Frequency float64 `option:"freq,frequency"`
	MaxShift  float64 `option:"maxShift"`
	AbsTime   bool    `option:"abs"`

	shifterCtx *waves.NoteCtx
}

func (Flanger) New() Filter {
	return &Flanger{
		shifter: waves.WaveFn(func(t float64, ctx *waves.NoteCtx) float64 {
			x := 2.0 * math.Pi * t / ctx.Period
			return math.Cos(x)
		}),
		MaxShift:  0.02,
		Frequency: 4.0,
		AbsTime:   false,
	}
}

func (f *Flanger) Apply(wave waves.Wave) waves.Wave {
	f.shifterCtx = waves.NewNoteCtx(f.Frequency, 1.0, math.Inf(1), 0.0)
	if f.AbsTime {
		return MakeFilterImpl(f, wave, flangerAbsImplValue)
	} else {
		return MakeFilterImpl(f, wave, flangerImplValue)
	}
}

func flangerImplValue(fx *Flanger, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	v := input.Value(t, ctx)
	shift := fx.shifter.Value(t, fx.shifterCtx) * fx.MaxShift
	t1 := t - 0.5*shift + 0.5*fx.MaxShift
	v1 := input.Value(t1, ctx)
	return 0.5*v + 0.5*v1
}

func flangerAbsImplValue(fx *Flanger, input waves.Wave, t float64, ctx *waves.NoteCtx) float64 {
	v := input.Value(t, ctx)
	shift := fx.shifter.Value(ctx.AbsTime, fx.shifterCtx) * fx.MaxShift
	t1 := t - 0.5*shift + 0.5*fx.MaxShift
	v1 := input.Value(t1, ctx)
	return 0.5*v + 0.5*v1
}
