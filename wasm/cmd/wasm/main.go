package main

import (
	"fmt"
	"runtime/debug"
	"syscall/js"

	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	synth "github.com/avoronkov/waver/lib/midisynth/unisynth"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/parser"
)

var (
	goParser    *parser.Parser
	goSequencer *seq.Sequencer
	goCfg       *config.Config
)

func doMain() {
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

	// Instruments
	instSet := instruments.NewSet()

	// Audio output
	audioOpts := []func(*synth.Output){
		synth.WithInstruments(instSet),
		synth.WithScale(scale),
		synth.WithTempo(tempo),
	}
	audioOutput, err := synth.New(audioOpts...)
	check(err)

	opts = append(opts, midisynth.WithSignalOutput(audioOutput))

	// Parser
	goParser = parser.New(
		goSequencer,
		scale,
		parser.WithTempoSetter(goSequencer),
		parser.WithTempoSetter(audioOutput),
		parser.WithInstrumentSet(instSet),
	)

	opts = append(opts,
		midisynth.WithSignalInput(goSequencer),
		midisynth.WithLoggingFunction(doLog),
	)

	m, err := midisynth.NewMidiSynth(opts...)
	check(err)

	//
	c := make(chan struct{}, 0)

	// Finally starting sequencer
	check(m.Start())
	check(m.Close())

	<-c
}

func main() {
	defer doRecover()
	doMain()
}

func doLog(format string, v ...any) {
	js.Global().Call("logMessage", fmt.Sprintf(format, v...))
}

func doRecover() {
	if r := recover(); r != nil {
		doLog("RECOVERED: %v\n%s", r, debug.Stack())
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
