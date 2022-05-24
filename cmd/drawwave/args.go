package main

import "flag"

var (
	input string
)

func init() {
	flag.StringVar(&input, "input", "", "input wav file")
	flag.StringVar(&input, "i", "", "input wav file (shorthand)")
}
