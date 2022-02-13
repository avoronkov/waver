package main

import "flag"

var (
	input      string
	output     string
	configPath string
	edo19      bool
)

func init() {
	flag.StringVar(&input, "input", "", "input file")
	flag.StringVar(&input, "i", "", "input file (shorthand)")

	flag.StringVar(&output, "output", "", "output file")
	flag.StringVar(&output, "o", "", "output file (shorthand)")

	flag.StringVar(&configPath, "config", "", "instruments config file")
	flag.StringVar(&configPath, "c", "", "instruments config file (shorthand)")

	flag.BoolVar(&edo19, "edo19", false, "Use EDO19 scale")
	flag.BoolVar(&edo19, "19", false, "Use EDO19 scale (shorthand)")
}
