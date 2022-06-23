package parser

import (
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/seq/types"
)

type Seq interface {
	Add(types.Signaler)
	Commit() error
	Assign(name string, value types.ValueFn)
}

type TempoSetter interface {
	SetTempo(int)
}

type InstrumentSet interface {
	AddInstrument(n string, in *instruments.Instrument)
}

type UserFunction struct {
	name string
	arg  string
	fn   types.ValueFn
}
