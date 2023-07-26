package main

import "flag"

var (
	wavFile   string
	forthFile string
	outFile   string
)

func init() {
	flag.StringVar(&wavFile, "w", "", "wav file to process")
	flag.StringVar(&forthFile, "f", "", "forth instructions")
	flag.StringVar(&outFile, "o", "", "output file")
}
