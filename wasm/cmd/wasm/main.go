package main

import (
	"fmt"
	"syscall/js"

	"gitlab.com/avoronkov/waver/etc"
	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/synth"
	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/parser"
)

var goParser *parser.Parser

func main() {
	tempo := 120
	var startBit int64 = 0

	opts := []func(*midisynth.MidiSynth){}

	scale := notes.NewStandard()
	common.Scale = scale

	sequencer := seq.NewSequencer(
		seq.WithTempo(tempo),
		seq.WithStart(startBit),
	)

	goParser = parser.New(sequencer, scale)
	// TODO method to send data to parser

	opts = append(opts, midisynth.WithSignalInput(sequencer))

	// Instruments
	instSet := instruments.NewSet()
	cfg := config.New("", instSet)
	check(cfg.UpdateData(etc.DefaultConfig))

	// Audio output
	audioOpts := []func(*synth.Output){
		synth.WithInstruments(instSet),
		synth.WithScale(scale),
		synth.WithTempo(tempo),
	}
	audioOutput, err := synth.New(audioOpts...)
	check(err)

	opts = append(opts, midisynth.WithSignalOutput(audioOutput))

	m, err := midisynth.NewMidiSynth(opts...)
	check(err)

	//
	c := make(chan struct{}, 0)
	fmt.Println("Go Web Assembly")

	js.Global().Set("goPlay", js.FuncOf(jsPlay))

	// Finally starting sequencer
	check(m.Start())
	check(m.Close())

	<-c
}
