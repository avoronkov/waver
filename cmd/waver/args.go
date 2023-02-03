package main

import (
	"flag"
)

var (
	udpPort  int
	midiPort int
	edo19    bool

	configPath string
	useConfig  bool

	dump string

	fileInput string

	newProject string

	tempo    int
	startBit int64

	showBits int64

	dumpWav bool

	shareFile string
)

func init() {
	flag.IntVar(&udpPort, "udp-port", 0, "UPD port (e.g. 49161)")
	flag.IntVar(&udpPort, "p", 0, "UPD port (shorthand, e.g. 49161)")

	flag.IntVar(&midiPort, "midi-port", 0, "MIDI port")
	flag.IntVar(&midiPort, "m", 0, "MIDI port (shorthand)")

	flag.BoolVar(&edo19, "edo19", false, "Use EDO19 scale")
	flag.BoolVar(&edo19, "19", false, "Use EDO19 scale (shorthand)")

	defaultConfig := "./etc/config.yml"
	flag.StringVar(&configPath, "config-path", defaultConfig, "Midisynth configuration")
	flag.StringVar(&configPath, "c", defaultConfig, "Midisynth configuration (shorthand)")

	flag.BoolVar(&useConfig, "config", false, "use instruments yaml config")

	flag.StringVar(&dump, "dump", "", "dump signals into file")
	flag.StringVar(&dump, "d", "", "dump signals into file (shorthand)")

	flag.StringVar(&fileInput, "input-file", "", "input sequencer file")
	flag.StringVar(&fileInput, "i", "", "input sequencer file (shorthand)")

	flag.StringVar(&newProject, "new", "", "initialize new project")

	flag.IntVar(&tempo, "tempo", 120, "set tempo")
	flag.IntVar(&tempo, "t", 120, "set tempo")

	flag.Int64Var(&startBit, "start", 0, "starting bit")

	flag.Int64Var(&showBits, "show-bits", 0, "Log bit number every n bits")

	flag.BoolVar(&dumpWav, "dump-wav", false, "dump audio output into wav file")

	flag.StringVar(&shareFile, "share", "", "Create sharable link with file content")
}
