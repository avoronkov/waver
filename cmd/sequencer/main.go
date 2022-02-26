package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/avoronkov/waver/lib/seq"
	"gitlab.com/avoronkov/waver/lib/seq/parser"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "Input file not specified.\n")
		flag.Usage()
		os.Exit(2)
	}

	input := flag.Arg(0)

	sequencer := seq.NewSequencer()
	ps := parser.New(input, sequencer)

	if err := ps.Start(); err != nil {
		log.Fatal(err)
	}
}
