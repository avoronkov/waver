package pragma

import (
	"github.com/avoronkov/waver/lib/midisynth/instruments"
)

type TempoSetter interface {
	SetTempo(int)
}

type InstrumentSet interface {
	AddInstrument(n int, in *instruments.Instrument)
	AddSampledInstrument(name string, in *instruments.Instrument)
}
