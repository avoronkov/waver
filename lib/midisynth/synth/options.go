package synth

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func WithWavSettings(settings *wav.Settings) func(*Output) {
	return func(o *Output) {
		o.settings = settings
	}
}

func WithScale(scale notes.Scale) func(*Output) {
	return func(o *Output) {
		o.scale = scale
	}
}

func WithInstruments(set InstrumentSet) func(*Output) {
	return func(o *Output) {
		o.instruments = set
	}
}

func WithTempo(n int) func(*Output) {
	return func(o *Output) {
		o.tempo = n
	}
}
