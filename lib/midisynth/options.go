package midisynth

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
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
