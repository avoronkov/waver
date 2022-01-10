package instruments

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Instrument struct {
	initialWave waves.Wave

	fx []filters.Filter

	resultingWave waves.Wave
}

func NewInstrument(wave waves.Wave, fx ...filters.Filter) *Instrument {
	in := &Instrument{
		initialWave: wave,
		fx:          fx,
	}

	w := in.initialWave
	for _, f := range in.fx {
		w = f.Apply(w)
	}
	in.resultingWave = w

	return in
}

func (i *Instrument) Wave() waves.Wave {
	return i.resultingWave
}

func (i *Instrument) WaveControlled() waves.WaveControlled {
	w := i.initialWave
	for _, f := range i.fx[:len(i.fx)-1] {
		w = f.Apply(w)
	}
	manual := filters.NewManualControlFilter(0.0)
	return manual.Apply(w).(waves.WaveControlled)
}
