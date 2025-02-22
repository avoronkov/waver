package parser

import (
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/types"
)

type Seq interface {
	Add(types.Signaler)
	Commit() error
	Assign(name string, value types.ValueFn)
	SetStopBit(bit int64)
}

type TempoSetter interface {
	SetTempo(int)
}

type ScaleSetter func(notes.Scale)

type InstrumentSet interface {
	AddInstrument(n string, in *instruments.Instrument)
	HasInstrument(n string) bool
}

type UserFunction struct {
	Name string
	Arg  string
	Fn   types.ValueFn
}
