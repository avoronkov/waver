package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/midi"
	"gitlab.com/avoronkov/waver/lib/midisynth/synth"
	"gitlab.com/avoronkov/waver/lib/midisynth/udp"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func main() {
	flag.Parse()

	udpInput := udp.New(udpPort)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithSignalInput(udpInput),
	}

	if midiPort > 0 {
		midiInput := midi.NewInput(midiPort)
		opts = append(opts, midisynth.WithSignalInput(midiInput))
	}

	// Instruments
	instSet := instruments.NewSet()
	cfg := config.New(configPath, instSet)
	if err := cfg.InitMidiSynth(); err != nil {
		log.Fatal(err)
	}
	if err := cfg.StartUpdateLoop(); err != nil {
		log.Fatal(err)
	}
	// .

	// Audio output
	audioOpts := []func(*synth.Output){
		synth.WithInstruments(instSet),
	}
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		audioOpts = append(audioOpts, synth.WithScale(notes.NewEdo19()))
	} else {
		log.Printf("Using Standard 12 tone scale.")
		audioOpts = append(audioOpts, synth.WithScale(notes.NewStandard()))
	}
	audioOutput, err := synth.New(audioOpts...)
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, midisynth.WithSignalOutput(audioOutput))
	// .

	m, err := midisynth.NewMidiSynth(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Experimantal section
	/*
		in := instruments.NewInstrument(
			waves.SineSine,
			filters.NewAdsrFilter(),
		)

		m.AddInstrument(10, in)

	*/
	// .

	if err := m.Start(); err != nil {
		log.Fatal("Start failed: ", err)
	}
	if err := m.Close(); err != nil {
		log.Fatal("Stop failed: ", err)
	}
	log.Printf("OK!")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
