package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/avoronkov/waver/etc"
	"gitlab.com/avoronkov/waver/lib/midisynth"
	"gitlab.com/avoronkov/waver/lib/midisynth/config"
	"gitlab.com/avoronkov/waver/lib/midisynth/dumper"
	"gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/midi"
	"gitlab.com/avoronkov/waver/lib/midisynth/synth"
	"gitlab.com/avoronkov/waver/lib/midisynth/udp"
	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/project"
	"gitlab.com/avoronkov/waver/lib/seq"
	"gitlab.com/avoronkov/waver/lib/seq/parser"
)

func main() {
	flag.Parse()

	if newProject != "" {
		if err := project.New(newProject, etc.DefaultConfig); err != nil {
			log.Fatal(err)
		}
		return
	}

	udpInput := udp.New(udpPort)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithSignalInput(udpInput),
	}

	if midiPort > 0 {
		midiInput := midi.NewInput(midiPort)
		opts = append(opts, midisynth.WithSignalInput(midiInput))
	}

	if fileInput != "" {
		sequencer := seq.NewSequencer()
		ps := parser.New(fileInput, sequencer)
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
	audioOpts := []func(*synth.Output){
		synth.WithInstruments(instSet),
	}
	if edo19 {
		log.Printf("Using EDO-19 scale.")
		audioOpts = append(audioOpts, synth.WithScale(notes.NewEdo19()))
	} else {
		log.Printf("Using Standard 12 tone scale.")
		audioOpts = append(audioOpts, synth.WithScale(notes.NewStandard()))
	}
	audioOutput, err := synth.New(audioOpts...)
	check("Syntheziser output", err)
	opts = append(opts, midisynth.WithSignalOutput(audioOutput))
	// .

	m, err := midisynth.NewMidiSynth(opts...)
	check("Midisynth creation", err)

	// Experimantal section
	/*
		in := instruments.NewInstrument(
			waves.SineSine,
			filters.NewAdsrFilter(),
		)

		m.AddInstrument(10, in)

	*/
	// .

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
