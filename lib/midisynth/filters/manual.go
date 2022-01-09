// +build wip

package filters

import "gitlab.com/avoronkov/waver/lib/midisynth/waves"

type ManualControl struct {
	Release   float64
	releasing bool
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
	wave waves.Wave
	opts *ManualControl
}

func (i *manualControlImpl) Value(t float64, ctx *waves.NoteCtx) float64 {

}
