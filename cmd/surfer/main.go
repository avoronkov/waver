package main

import (
	"fmt"
	"log"
	"os"

	"github.com/avoronkov/waver/lib/midisynth/waves"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}
	file := os.Args[1]
	sample, err := waves.ParseSampleFile(file)
	if err != nil {
		log.Fatalf("ParseSampleFile failed: %v", err)
	}

	data := sample.Data()
	slices := SlicesFromSamples(data)
	for _, slice := range slices {
		fmt.Printf("(%v) %v\n", len(slice), slice)
	}
	fmt.Println(len(slices))
}
