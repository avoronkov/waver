package main

import "flag"

var (
	port  int
	edo19 bool
)

func init() {
	flag.IntVar(&port, "port", 49161, "UPD port")
	flag.IntVar(&port, "p", 49161, "UPD port (shorthand)")

	flag.BoolVar(&edo19, "edo19", false, "Use EDO19 scale")
	flag.BoolVar(&edo19, "19", false, "Use EDO19 scale (shorthand)")
}
