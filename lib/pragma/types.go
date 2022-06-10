package pragma

import (
	"github.com/avoronkov/waver/lib/midisynth/instruments"
)

type TempoSetter interface {
	SetTempo(int)
}

type InstrumentSet interface {
	AddInstrument(n string, in *instruments.Instrument)
}
