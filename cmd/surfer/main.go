package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/avoronkov/waver/lib/surfer"
)

func main() {
	flag.Parse()
	if wavFile == "" {
		log.Fatal("Wav file is not specified")
	}
	if forthFile == "" {
		log.Fatal("Forth file is not specified")
	}
	if outFile == "" {
		log.Fatal("Output file is not specified")
	}

	in := surfer.NewInterpreter()
	err := in.Run(forthFile, wavFile, outFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, "OK")
}
