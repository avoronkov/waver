package main

import "flag"

var (
	udpPort  int
	midiPort int
	edo19    bool

	configPath string

	dump string
)

func init() {
	flag.IntVar(&udpPort, "udp-port", 49161, "UPD port")
	flag.IntVar(&udpPort, "p", 49161, "UPD port (shorthand)")

	flag.IntVar(&midiPort, "midi-port", 0, "MIDI port")
	flag.IntVar(&midiPort, "m", 0, "MIDI port (shorthand)")

	flag.BoolVar(&edo19, "edo19", false, "Use EDO19 scale")
	flag.BoolVar(&edo19, "19", false, "Use EDO19 scale (shorthand)")

	defaultConfig := "./etc/config.yml"
	flag.StringVar(&configPath, "config-path", defaultConfig, "Midisynth configuration")
	flag.StringVar(&configPath, "c", defaultConfig, "Midisynth configuration (shorthand)")

	flag.StringVar(&dump, "dump", "", "dump signals into file")
	flag.StringVar(&dump, "d", "", "dump signals into file (shorthand)")
}
