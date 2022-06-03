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

var (
	goParser    *parser.Parser
	goSequencer *seq.Sequencer
	goCfg       *config.Config
)

func main() {
	// Export JS functions
	fmt.Println("Go Web Assembly: goPlay, goGetDefaultCode")

	js.Global().Set("goPlay", js.FuncOf(jsPlay))
	js.Global().Set("goGetDefaultCode", js.FuncOf(jsGetDefaultCode))
	js.Global().Set("goPause", js.FuncOf(jsPause))
	js.Global().Set("goUpdateInstruments", js.FuncOf(jsUpdateInstruments))
	js.Global().Set("goGetDefaultInstruments", js.FuncOf(jsGetDefaultInstruments))

	js.Global().Call("initPage")

	// Init waver
	tempo := 120
	var startBit int64 = 0

	opts := []func(*midisynth.MidiSynth){}

	scale := notes.NewStandard()
	common.Scale = scale

	goSequencer = seq.NewSequencer(
		seq.WithTempo(tempo),
		seq.WithStart(startBit),
	)

	goParser = parser.New(goSequencer, scale)
	// TODO method to send data to parser

	opts = append(opts, midisynth.WithSignalInput(goSequencer))

	// Instruments
	instSet := instruments.NewSet()
	goCfg = config.New("", instSet)
	check(goCfg.UpdateData(etc.DefaultConfig))

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

	// Finally starting sequencer
	check(m.Start())
	check(m.Close())

	<-c
}
