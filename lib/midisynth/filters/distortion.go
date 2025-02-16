package filters

import (
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

// Amplitude distortion filter
type Distortion struct {
	Value float64 `option:"value"`
}

func (Distortion) New() Filter {
	return &Distortion{
		Value: 1.0,
	}
}

func (Distortion) Desc() string {
	return `Distortion is a form of audio signal processing used to alter the sound of amplified electric musical instruments, usually by increasing their gain, producing a "fuzzy", "growling", or "gritty" tone.`
}

func (df *Distortion) Apply(w waves.Wave) waves.Wave {
	return &distortionImpl{
		wave: w,
		opts: df,
	}
}

type distortionImpl struct {
	wave waves.Wave
	opts *Distortion
}

func (i *distortionImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	val := i.wave.Value(t, ctx) * i.opts.Value
	if val > 1.0 {
		val = 1.0
	} else if val < -1.0 {
		val = -1.0
	}
	return val
}
