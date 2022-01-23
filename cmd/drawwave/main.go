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
	snare, err := waves.ReadSample("../../samples/4-snare.wav")
	if err != nil {
		panic(err)
	}

	// Instrument
	in := instruments.NewInstrument(
		snare,
		// filters.NewAdsrFilter(),
		filters.NewDelayFilter(
			filters.DelayInterval(0.125),
			filters.DelayFadeOut(0.5),
			filters.DelayTimes(2),
		),
	)
	// .

	play := player.New(wav.Default)
	hz := 55.0
	amp := 1.0
	dur := 1.0
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
