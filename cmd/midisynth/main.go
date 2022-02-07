package main

import (
	"flag"
	"log"

	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/filters"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/midi"
	"gitlab.com/avoronkov/waver/lib/midisynth/udp"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
	"gitlab.com/avoronkov/waver/lib/notes"
)

func main() {
	flag.Parse()

	udpInput := udp.New(udpPort)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithSignalInput(udpInput),
	}
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		opts = append(opts, midisynth.WithScale(notes.NewEdo19()))
	} else {
		log.Printf("Using Standard 12 tone scale.")
		opts = append(opts, midisynth.WithScale(notes.NewStandard()))
	}

	m, err := midisynth.NewMidiSynth(opts...)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.New(configPath, m)
	if err := cfg.InitMidiSynth(); err != nil {
		log.Fatal(err)
	}

	if err := cfg.StartUpdateLoop(); err != nil {
		log.Fatal(err)
	}

	if midiPort > 0 {
		opts := []func(*midi.Proc){
			midi.WithDispatcher(cfg),
		}
		if edo19 {
			opts = append(opts, midi.Edo19())
		}

		midiProc := midi.NewProc(m, midiPort, opts...)
		midisynth.WithMidiProc(midiProc)(m)
	}

	// Experimantal section
	in := instruments.NewInstrument(
		waves.SineSine,
		filters.NewAdsrFilter(),
	)

	m.AddInstrument(10, in)

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
