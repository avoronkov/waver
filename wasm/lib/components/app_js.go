//go:build js

package components

import (
	"log"

	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/signals"
	synth "github.com/avoronkov/waver/lib/midisynth/unisynth"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/parser"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (ap *App) OnNav(ctx app.Context) {
	log.Printf("[JS] initMain()...")

	tempo := 120
	var startBit int64 = 0

	channel := make(chan signals.Interface, 128)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithChannel(channel),
	}

	scale := notes.NewStandard()
	common.Scale = scale

	goSequencer := seq.NewSequencer(
		seq.WithTempo(tempo),
		seq.WithStart(startBit),
		seq.WithChannel(channel),
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
	goParser := parser.New(
		goSequencer,
		scale,
		parser.WithTempoSetter(goSequencer),
		parser.WithInstrumentSet(instSet),
	)

	opts = append(opts,
		midisynth.WithSignalInput(goSequencer),
		midisynth.WithLoggingFunction(doLog),
	)

	midiSynth, err := midisynth.NewMidiSynth(opts...)
	check(err)

	// Finally starting sequencer
	// j.check(m.Start())
	// j.check(m.Close())

	ap.mParser = goParser
	ap.mSequencer = goSequencer
	ap.mMidiSynth = midiSynth

	go func() {
		log.Printf("mMidiSynth.Start()...")
		check(ap.mMidiSynth.Start())
		log.Printf("mMidiSynth.Finished()")
	}()
}
