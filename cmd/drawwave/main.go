package main

import (
	"io"
	"log"
	"waver/lib/midisynth/filters"
	instr "waver/lib/midisynth/instruments"
	"waver/lib/midisynth/player"
	"waver/lib/midisynth/wav"
	waves2 "waver/lib/midisynth/waves/v2"
)

func main() {
	in := instr.NewInstrument(&waves2.Saw{}, filters.NewAdsrFilter())
	play := player.New(wav.Default)
	hz := 440.0
	amp := 1.0
	dur := 0.25
	reader, _ := play.PlayContext(in.Wave(), waves2.NewNoteCtx(hz, amp, dur))

	data, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	draw := &Drawer{}
	if err := draw.Draw(data, "wave.html"); err != nil {
		log.Fatal(err)
	}
}
