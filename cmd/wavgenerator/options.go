package main

import (
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func WithInput(file string) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.input = file
	}
}

func WithOutput(file string) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.output = file
	}
}

func WithInstruments(set *instruments.Set) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.instSet = set
	}
}

func WithWavSettings(settings *wav.Settings) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.settings = settings
	}
}

func WithScale(scale notes.Scale) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.scale = scale
	}
}
