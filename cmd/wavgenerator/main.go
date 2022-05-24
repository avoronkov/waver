package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func main() {
	flag.Parse()

	instSet := instruments.NewSet()
	cfg := config.New(configPath, instSet)
	if err := cfg.InitMidiSynth(); err != nil {
		log.Fatal(err)
	}

	opts := []func(*WavGenerator){
		WithInput(input),
		WithOutput(output),
		WithInstruments(instSet),
	}
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		opts = append(opts, WithScale(notes.NewEdo19()))
	} else {
		log.Printf("Using Standard 12 tone scale.")
		opts = append(opts, WithScale(notes.NewStandard()))
	}
	if mono {
		log.Printf("Using mono output.")
		opts = append(opts, WithChannels(1))
	}
	gen, err := NewWavGenerator(opts...)
	if err != nil {
		log.Fatal(err)
	}
	if err := gen.Run(); err != nil {
		log.Fatal(err)
	}
}
