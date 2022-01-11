package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
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

	// Experimental section

	m.AddInstrument(8, instruments.NewInstrument(
		&waves.Sine{},
		filters.NewVibrato(&waves.Sine{}, 10.0, 0.05),
		filters.NewAdsrFilter(),
	))

	m.AddInstrument(9, instruments.NewInstrument(
		&waves.Triangle{},
		filters.NewRing(&waves.Sine{}, 4.0, 1.0),
		filters.NewAdsrFilter(),
	))

	m.AddInstrument(7, instruments.NewInstrument(
		&waves.Sine{},
		filters.NewTimeShift(10.0, 0.01),
		filters.NewAdsrFilter(),
	))

	// .

	m.Start()
	if err := m.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("OK!")
}
