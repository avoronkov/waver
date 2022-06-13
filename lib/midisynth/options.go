package midisynth

import (
	"github.com/avoronkov/waver/lib/midisynth/signals"
)

func WithSignalInput(input signals.Input) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.inputs = append(m.inputs, input)
	}
}

func WithSignalOutput(output signals.Output) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.outputs = append(m.outputs, output)
	}
}

func WithLoggingFunction(logf func(format string, v ...any)) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.logf = logf
	}
}

func WithChannel(ch chan signals.Interface) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.ch = ch
	}
}
