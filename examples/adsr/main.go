package main

import (
	"fmt"
	"log"
	"os"
	n "waver/notes"
	"waver/wav"
)

func main() {
	w := wav.CreateDefaultWav()

	signal := &Signal{
		AttackLen:   1000,
		AttackLevel: 16000,
		DecayLen:    10000,
		DecayLevel:  8000,
		SusteinLen:  100,
		ReleaseLen:  13000,
	}

	for _, note := range []float64{n.C2, n.D2, n.E2, n.F2, n.G2, n.A2, n.B2, n.C3, n.B2, n.A2, n.G2, n.F2, n.E2, n.D2, n.C2} {
		signal.PutSignal(note, w.Data)
	}

	name := "signal.wav"
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := w.Write(f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v generated.", name)
}
