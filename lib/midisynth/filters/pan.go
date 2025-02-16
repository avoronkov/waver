package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type Pan struct {
	Left  float64 `option:"left,l"`
	Right float64 `option:"right,r"`
}

func (Pan) New() Filter {
	return &Pan{
		Left:  1.0,
		Right: 1.0,
	}
}

func (Pan) Desc() string {
	return `Panning is the distribution of an audio signal into a new stereo or multi-channel sound field determined by a pan control setting.`
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
	ch := i.opts.Left
	if ctx.Channel == 1 {
		ch = i.opts.Right
	}
	return i.wave.Value(tm, ctx) * ch
}
