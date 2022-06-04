package main

import (
	"io"
	"testing"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/player"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

func BenchmarkHarmonizerFlanger(b *testing.B) {
	in := instruments.NewInstrument(
		waves.Triangle,
		filters.NewHarmonizer(),
		filters.NewFlanger(filters.FlangerFreq(0.1)),
		filters.NewAdsrFilter(),
	)

	play := player.New(wav.Default)
	hz := 55.0
	amp := 1.0
	dur := 0.1
	for i := 0; i < b.N; i++ {
		reader, _ := play.PlayContext(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

		_, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFlanger(b *testing.B) {
	in := instruments.NewInstrument(
		waves.Triangle,
		filters.NewFlanger(filters.FlangerFreq(0.1)),
		filters.NewAdsrFilter(),
	)

	play := player.New(wav.Default)
	hz := 55.0
	amp := 1.0
	dur := 0.1
	for i := 0; i < b.N; i++ {
		reader, _ := play.PlayContext(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

		_, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
	}
}
