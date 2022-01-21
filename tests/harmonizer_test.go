package main

import (
	"io"
	"testing"

	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/player"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
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
		reader, _ := play.PlayContext2(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

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
		reader, _ := play.PlayContext2(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

		_, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
	}
}
