package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type ManualControl struct {
	Release float64
}

func NewManualControlFilter(release float64) Filter {
	return &ManualControl{
		Release: release,
	}
}

func (mc *ManualControl) Apply(w waves.Wave) waves.Wave {
	return &manualControlImpl{
		wave: w,
		opts: mc,
	}
}

type manualControlImpl struct {
	wave             waves.Wave
	opts             *ManualControl
	releasing        bool
	releaseStartTime *float64
}

func (i *manualControlImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	if i.releasing {
		return math.NaN()
	}
	return i.wave.Value(t, ctx)
}

func (i *manualControlImpl) Release() {
	i.releasing = true
}
