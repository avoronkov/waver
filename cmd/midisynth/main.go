package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config/v2"
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

	// Experimantal section

	kick, err := waves.ReadSample("./samples/4-kick.wav")
	check(err)

	kickInst := instruments.NewInstrument(kick)
	m.AddInstrument(20, kickInst) // 'k'

	hat, err := waves.ReadSample("./samples/4-hat.wav")
	check(err)
	hatInst := instruments.NewInstrument(
		hat,
		filters.NewDelayFilter(
			filters.DelayTimes(3),
			filters.DelayInterval(0.125),
			filters.DelayFadeOut(0.5),
		),
	)
	m.AddInstrument(17, hatInst) // 'h'

	snare, err := waves.ReadSample("./samples/4-snare.wav")
	check(err)
	snareInst := instruments.NewInstrument(snare)
	m.AddInstrument(28, snareInst)

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
