package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/avoronkov/waver/etc"
	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/midisynth/dumper"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/midi"
	"github.com/avoronkov/waver/lib/midisynth/udp"
	"github.com/avoronkov/waver/lib/midisynth/unisynth"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/pragma"
	"github.com/avoronkov/waver/lib/project"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/parser"
)

func main() {
	flag.Parse()

	if newProject != "" {
		if err := project.New(newProject, etc.DefaultConfig); err != nil {
			log.Fatal(err)
		}
		return
	}

	var scale notes.Scale
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		scale = notes.NewEdo19()
	} else {
		log.Printf("Using Standard 12 tone scale.")
		scale = notes.NewStandard()
	}

	// TODO fix this a little
	common.Scale = scale

	udpInput := udp.New(udpPort, scale)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithSignalInput(udpInput),
	}

	if midiPort > 0 {
		midiInput := midi.NewInput(midiPort, scale)
		opts = append(opts, midisynth.WithSignalInput(midiInput))
	}

	var sequencer *seq.Sequencer
	if fileInput != "" {
		sequencer = seq.NewSequencer(
			seq.WithTempo(tempo),
			seq.WithStart(startBit),
			seq.WithShowingBits(showBits),
		)
		ps := parser.New(sequencer, scale, parser.WithFileInput(fileInput))
		check("Parser start", ps.Start(true))
		opts = append(opts, midisynth.WithSignalInput(sequencer))
	}

	if dump != "" {
		dumpOutput, err := dumper.NewJson(dump)
		check("Json dumper creation", err)
		opts = append(opts, midisynth.WithSignalOutput(dumpOutput))
	}

	// Instruments
	instSet := instruments.NewSet()
	cfg := config.New(getConfigPath(), instSet)
	check("MidiSynth initialization", cfg.InitMidiSynth())
	check("Config StartUpdateLoop", cfg.StartUpdateLoop())
	// .

	// Audio output
	audioOpts := []func(*unisynth.Output){
		unisynth.WithInstruments(instSet),
		unisynth.WithScale(scale),
		unisynth.WithTempo(tempo),
	}
	audioOutput, err := unisynth.New(audioOpts...)
	check("Syntheziser output", err)
	opts = append(opts, midisynth.WithSignalOutput(audioOutput))
	// .

	// Pragma parser
	if fileInput != "" {
		pragmaParser := pragma.New(
			fileInput,
			pragma.WithTempoSetter(sequencer),
			pragma.WithTempoSetter(audioOutput),
		)
		check("Pragma parser start", pragmaParser.Start(true))
	}
	// .

	m, err := midisynth.NewMidiSynth(opts...)
	check("Midisynth creation", err)

	check("Start", m.Start())
	check("Stop", m.Close())
	log.Printf("OK!")
}

func check(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func getConfigPath() string {
	log.Printf("[debug] fileInput = %v", fileInput)
	if fileInput != "" {
		confPath := fmt.Sprintf("%v.yml", strings.TrimSuffix(fileInput, filepath.Ext(fileInput)))
		log.Printf("[debug] confPath = %v", confPath)
		if _, err := os.Stat(confPath); err == nil {
			return confPath
		}
	}
	return configPath
}
