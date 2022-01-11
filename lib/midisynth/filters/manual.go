package filters

import (
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type ManualControl struct {
	Release float64
}

var _ FilterManualControl = (*ManualControl)(nil)

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

func (mc *ManualControl) IsManualControl() {}

type manualControlImpl struct {
	wave             waves.Wave
	opts             *ManualControl
	releasing        bool
	releaseStartTime *float64
}

func float64ptr(x float64) *float64 { return &x }

func (i *manualControlImpl) Value(t float64, ctx *waves.NoteCtx) float64 {
	mult := 1.0
	if i.releasing {
		if i.releaseStartTime == nil {
			i.releaseStartTime = float64ptr(t)
		}
		if t > *(i.releaseStartTime)+i.opts.Release {
			return math.NaN()
		}
		mult = 1.0 - (t-*i.releaseStartTime)/i.opts.Release
	}
	return i.wave.Value(t, ctx) * ctx.Amp * mult
}

func (i *manualControlImpl) Release() {
	i.releasing = true
}
