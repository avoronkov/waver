package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"waver/wav"
)

func Sample() {
	w := wav.CreateDefaultWav()
	_ = w

	samplesPerSecond := 44100
	samples := samplesPerSecond * 4

	hz := 440.0

	var waveDuration float64 = float64(samplesPerSecond) / hz

	amp := 32000.0
	for i := 0; i < samples; i++ {
		x := 2.0 * math.Pi * float64(i) / waveDuration
		l := amp * math.Sin(x)
		w.Data.AddSample(int16(l)) // left
		// r := amp * math.Cos(x)
		r := l
		w.Data.AddSample(int16(r)) // right }
	}
	for i := 0; i < samples; i++ {
		x := 2.0 * math.Pi * float64(i) / waveDuration
		l := amp * math.Cos(x)
		w.Data.AddSample(int16(l)) // left
		// r := amp * math.Sin(x)
		r := -l
		w.Data.AddSample(int16(r)) // right }
	}

	f, err := os.OpenFile("sample.wav", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := w.Write(f); err != nil {
		log.Fatal(err)
	}
	fmt.Println("sample.wav generated.")
}

func main() {
	Sample()
}
