package main

import "flag"

var (
	port  int
	edo19 bool

	configPath string
)

func init() {
	flag.IntVar(&port, "port", 49161, "UPD port")
	flag.IntVar(&port, "p", 49161, "UPD port (shorthand)")

	flag.BoolVar(&edo19, "edo19", false, "Use EDO19 scale")
	flag.BoolVar(&edo19, "19", false, "Use EDO19 scale (shorthand)")

	flag.StringVar(&configPath, "config-path", "./etc/config.yaml", "Midisynth configuration")
	flag.StringVar(&configPath, "c", "./etc/config.json", "Midisynth configuration (shorthand)")
}
