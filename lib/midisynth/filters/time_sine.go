package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type SinusTime struct{}

func NewSinusTime() Filter {
	return &SinusTime{}
}

func (s *SinusTime) Apply(wave waves.Wave) waves.Wave {
	return &sinusTimeImpl{
		wave: wave,
		opts: s,
	}
}

type sinusTimeImpl struct {
	wave waves.Wave
	opts *SinusTime
}

func (i *sinusTimeImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	newT := math.Pi * (math.Sin(2*math.Pi*t/ctx.Period) + 1.0)
	return i.wave.Value(newT, ctx)
}
