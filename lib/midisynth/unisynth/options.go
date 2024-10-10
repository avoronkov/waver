package unisynth

import (
	"time"

	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/notes"
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

func WithWavFileDump(file string) func(*Output) {
	return func(o *Output) {
		o.wavFilename = file
	}
}

func WithWavSpaceRight(value float64) func(*Output) {
	return func(o *Output) {
		o.wavSpaceRight = value
	}
}

func WithWavSpaceLeft(value float64) func(*Output) {
	return func(o *Output) {
		o.wavSpaceLeft = value
	}
}

func WithPlayer(p PlayCloser) func(*Output) {
	return func(o *Output) {
		o.player = p
	}
}

func WithDelay(d time.Duration) func(*Output) {
	return func(o *Output) {
		o.delay = d
	}
}
