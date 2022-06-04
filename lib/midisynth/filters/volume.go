package filters

import "github.com/avoronkov/waver/lib/midisynth/waves"

type Volume struct {
}

var _ Filter = (*Volume)(nil)

func (v *Volume) Apply(w waves.Wave) waves.Wave {
	return &volumeImpl{
		wave: w,
	}
}

type volumeImpl struct {
	wave waves.Wave
}

func (i *volumeImpl) Value(tm float64, ctx *waves.NoteCtx) float64 {
	return i.wave.Value(tm, ctx) * ctx.Amp
}
