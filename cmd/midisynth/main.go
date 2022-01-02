package main

import (
	"log"
	"waver/lib/midisynth"
	"waver/lib/notes"
)

func main() {
	port := 49161
	log.Printf("Starting midi syntesizer on port %v...", port)
	m, err := midisynth.NewMidiSynth(44100, 2, 2, notes.NewStandard(), port)
	if err != nil {
		log.Fatal(err)
	}
	m.Start()
	if err := m.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("OK!")
}
