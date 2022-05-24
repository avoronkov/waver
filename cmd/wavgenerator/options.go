package main

import (
	"fmt"

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

var _ = WithWavSettings

func WithScale(scale notes.Scale) func(*WavGenerator) {
	return func(g *WavGenerator) {
		g.scale = scale
	}
}

func WithChannels(n int) func(*WavGenerator) {
	if !(n == 1 || n == 2) {
		panic(fmt.Errorf("Unsupported number of channels: %v", n))
	}
	return func(g *WavGenerator) {
		g.channels = n
	}
}
