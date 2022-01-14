package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config/v2"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func main() {
	flag.Parse()
	log.Printf("Starting midi syntesizer on port %v...", port)
	var scale notes.Scale
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		scale = notes.NewEdo19()
	} else {
		log.Printf("Using Standard 12 tone scale.")
		scale = notes.NewStandard()
	}
	m, err := midisynth.NewMidiSynth(wav.Default, scale, port)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &config.Config{}
	if err := cfg.InitMidiSynth(configPath, m); err != nil {
		log.Fatal(err)
	}

	// Experimantal section

	// .

	m.Start()
	if err := m.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("OK!")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
