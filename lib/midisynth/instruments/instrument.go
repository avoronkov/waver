package instruments

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Instrument struct {
	initialWave waves.Wave

	fx       []filters.Filter
	manualFx []filters.Filter

	resultingWave waves.Wave

	manualControl filters.Filter
	adsr          filters.Filter
}

var _ Interface = (*Instrument)(nil)

var defaultAdsr = filters.NewAdsrFilter()
var defaultManual = filters.NewManualControlFilter(0.125)

func NewInstrument(wave waves.Wave, fx ...filters.Filter) *Instrument {
	in := &Instrument{
		initialWave:   wave,
		manualControl: defaultManual,
		adsr:          defaultAdsr,
	}

	/// TODO fix manual / auto filters problem

	adsrApplied := false
	manualApplied := false
	waitManual := false
	waitAdsr := false

	for _, f := range fx {
		if waitManual {
			waitManual = false
			if _, ok := f.(filters.FilterManualControl); ok {
				in.manualFx = append(in.manualFx, f)
			} else {
				in.manualFx = append(in.manualFx, defaultManual)

				in.fx = append(in.fx, f)
				in.manualFx = append(in.manualFx, f)
			}
			continue
		}

		if waitAdsr {
			waitAdsr = false
			if _, ok := f.(filters.FilterAdsr); ok {
				in.fx = append(in.fx, f)
			} else {
				in.fx = append(in.fx, defaultAdsr)

				in.fx = append(in.fx, f)
				in.manualFx = append(in.manualFx, f)
			}
			continue
		}

		if _, ok := f.(filters.FilterAdsr); ok {
			in.fx = append(in.fx, f)
			adsrApplied = true
			waitManual = true
			continue
		}

		if _, ok := f.(filters.FilterManualControl); ok {
			in.manualFx = append(in.manualFx, f)
			manualApplied = true
			waitAdsr = true
			continue
		}
		in.fx = append(in.fx, f)
		in.manualFx = append(in.manualFx, f)
	}

	if !adsrApplied {
		in.fx = append(in.fx, defaultAdsr)
	}
	if !manualApplied {
		in.manualFx = append(in.manualFx, defaultManual)
	}

	// evaluate resulting wave
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
	for _, f := range i.manualFx {
		w = f.Apply(w)
	}
	// TODO this will break when delay is the last
	return w.(waves.WaveControlled)
}
