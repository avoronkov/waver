package filters

import "github.com/avoronkov/waver/lib/midisynth/waves"

type Filter interface {
	Apply(w waves.Wave) waves.Wave
}

type FilterAdsr interface {
	Filter

	IsAdsr()
}

type FilterManualControl interface {
	Filter

	IsManualControl()
}

func MakeFilterImpl[F any](
	filter F,
	input waves.Wave,
	impl func(filter F, input waves.Wave, t float64, ctx *waves.NoteCtx) float64,
) waves.Wave {
	waveFn := func(t float64, ctx *waves.NoteCtx) float64 {
		return impl(filter, input, t, ctx)
	}
	return waves.WaveFn(waveFn)
}
