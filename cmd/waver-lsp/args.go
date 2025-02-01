package main

import "flag"

var (
	debugFile bool
)

func init() {
	flag.BoolVar(&debugFile, "debug-file", false, "Write log into debug file")
}
