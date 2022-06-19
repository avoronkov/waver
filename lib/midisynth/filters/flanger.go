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
	shifterCtx := waves.NewNoteCtx(f.Frequency, 1.0, math.Inf(1), 0.0)
	if f.AbsTime {
		return &flangerAbsImpl{
			wave:       wave,
			opts:       f,
			shifterCtx: shifterCtx,
		}
	} else {
		return &flangerImpl{
			wave:       wave,
			opts:       f,
			shifterCtx: shifterCtx,
		}
	}
}

type flangerImpl struct {
	wave       waves.Wave
	opts       *Flanger
	shifterCtx *waves.NoteCtx
}

func (i *flangerImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.wave.Value(t, ctx)
	shift := i.opts.shifter.Value(t, i.shifterCtx) * i.opts.MaxShift
	t1 := t - 0.5*shift + 0.5*i.opts.MaxShift
	v1 := i.wave.Value(t1, ctx)
	return 0.5*v + 0.5*v1
}

type flangerAbsImpl struct {
	wave       waves.Wave
	opts       *Flanger
	shifterCtx *waves.NoteCtx
}

func (i *flangerAbsImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	v := i.wave.Value(t, ctx)
	shift := i.opts.shifter.Value(ctx.AbsTime, i.shifterCtx) * i.opts.MaxShift
	t1 := t - 0.5*shift + 0.5*i.opts.MaxShift
	v1 := i.wave.Value(t1, ctx)
	return 0.5*v + 0.5*v1
}
