package main

import "flag"

var (
	debugFile bool

	reference bool
)

func init() {
	flag.BoolVar(&debugFile, "debug-file", false, "Write log into debug file")
	flag.BoolVar(&reference, "ref", false, "Print language reference")
}
