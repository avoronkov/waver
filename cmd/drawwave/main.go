package main

import (
	"io"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/player"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

func main() {
	// Instrument
	in := instruments.NewInstrument(
		&waves.SineSine{},
		filters.NewAdsrFilter(),
	)
	// .

	play := player.New(wav.Default)
	hz := 440.0
	amp := 1.0
	dur := 0.1
	reader, _ := play.PlayContext2(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	draw := &Drawer{}
	if err := draw.Draw(data, "wave.html"); err != nil {
		log.Fatal(err)
	}
}
