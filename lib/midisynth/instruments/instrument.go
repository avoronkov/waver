package instruments

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Instrument struct {
	initialWave waves.Wave

	fx []filters.Filter

	resultingWave waves.Wave

	manualControl filters.Filter
	adsr          filters.Filter
}

var defaultAdsr = filters.NewAdsrFilter()
var defaultManual = filters.NewManualControlFilter(0.125)

func NewInstrument(wave waves.Wave, fx ...filters.Filter) *Instrument {
	in := &Instrument{
		initialWave:   wave,
		fx:            fx,
		manualControl: defaultManual,
		adsr:          defaultAdsr,
	}

	w := in.initialWave
	for _, f := range fx {
		w = f.Apply(w)

		if _, ok := f.(filters.FilterAdsr); ok {
			in.adsr = f
			continue
		}
		if _, ok := f.(filters.FilterManualControl); ok {
			in.manualControl = f
			continue
		}
		in.fx = append(in.fx, f)
	}
	in.resultingWave = in.adsr.Apply(w)

	return in
}

func (i *Instrument) Wave() waves.Wave {
	return i.resultingWave
}

func (i *Instrument) WaveControlled() waves.WaveControlled {
	w := i.initialWave
	for _, f := range i.fx {
		w = f.Apply(w)
	}
	return i.manualControl.Apply(w).(waves.WaveControlled)
}
