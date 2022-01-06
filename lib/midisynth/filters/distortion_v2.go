package filters

import "waver/lib/midisynth/waves/v2"

// Amplitude distortion filter
type DistortionFilter struct {
	Multiplier float64
}

func NewDistortionFilter(m float64) *DistortionFilter {
	return &DistortionFilter{
		Multiplier: m,
	}
}

func (df *DistortionFilter) Apply(w waves.Wave) waves.Wave {
	return &distortionImpl{
		wave: w,
		opts: df,
	}
}

type distortionImpl struct {
	wave waves.Wave
	opts *DistortionFilter
}

func (i *distortionImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	val := i.wave.Value(t, ctx) * i.opts.Multiplier
	if val > 1.0 {
		val = 1.0
	} else if val < -1.0 {
		val = -1.0
	}
	return val
}
