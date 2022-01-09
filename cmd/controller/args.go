package main

import "flag"

var (
	p int

	configPath string
)

func init() {
	flag.IntVar(&p, "p", 24, "client port number")

	flag.StringVar(&configPath, "config-path", "./etc/config.json", "Midisynth configuration")
	flag.StringVar(&configPath, "c", "./etc/config.json", "Midisynth configuration (shorthand)")
}
