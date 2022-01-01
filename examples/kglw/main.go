package main

import (
	"fmt"
	"log"
	"os"
	"waver/lib/adsr"
	n "waver/notes/kglw"
	"waver/wav"
)

func main() {
	w := wav.CreateDefaultWav()

	signal := &adsr.Signal{
		AttackLen:   1000,
		AttackLevel: 16000,
		DecayLen:    10000,
		DecayLevel:  8000,
		SusteinLen:  100,
		ReleaseLen:  3000,
	}

	for _, note := range []float64{n.I, n.II, n.III, n.IV, n.V, n.VI, n.VII, n.VIII, n.VII, n.VI, n.V, n.IV, n.III, n.II, n.I} {
		signal.PutSignal(note/2.0, w.Data)
	}

	name := "kglw.wav"
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
