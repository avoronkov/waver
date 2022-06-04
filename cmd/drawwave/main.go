package main

import (
	"flag"
	"io"
	"log"

	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/player"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

func main() {
	flag.Parse()
	if input == "" {
		log.Fatal("Input file name not specified")
	}
	snare, err := waves.ReadSample(input)
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
	reader, _ := play.PlayContext(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	draw := &Drawer{}
	if err := draw.Draw(data, input+".html"); err != nil {
		log.Fatal(err)
	}
}
