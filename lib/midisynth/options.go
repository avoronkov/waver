package midisynth

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func WithUdpPort(port int) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.udpPort = port
	}
}

func WithMidiPort(port int) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.midiPort = port
	}
}

func WithWavSettings(settings *wav.Settings) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.settings = settings
	}
}

func WithScale(scale notes.Scale) func(m *MidiSynth) {
	return func(m *MidiSynth) {
		m.scale = scale
		if edoScale, ok := scale.(notes.EdoScale); ok {
			m.edo = edoScale.Edo()
		} else {
			m.edo = -1
		}
	}
}
