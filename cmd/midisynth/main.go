package main

import (
	"log"
	"waver/lib/midisynth"
)

func main() {
	log.Printf("Strting midi syntesizer...")
	m, err := midisynth.NewMidiSynth(44100, 2, 2)
	if err != nil {
		log.Fatal(err)
	}
	m.Start()
	if err := m.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("OK!")
}
