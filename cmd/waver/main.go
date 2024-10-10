package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/avoronkov/waver/etc"
	"github.com/avoronkov/waver/lib/midisynth"
	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/midisynth/dumper"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/midi"
	"github.com/avoronkov/waver/lib/midisynth/output/pulse"
	"github.com/avoronkov/waver/lib/midisynth/signals"
	"github.com/avoronkov/waver/lib/midisynth/udp"
	"github.com/avoronkov/waver/lib/midisynth/unisynth"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/project"
	"github.com/avoronkov/waver/lib/seq"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/share"
)

func main() {
	flag.Parse()

	if newProject != "" {
		if err := project.New(newProject, etc.DefaultConfig); err != nil {
			log.Fatal(err)
		}
		return
	}

	if shareFile != "" {
		link, err := share.MakeLinkFromFile(shareFile)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(link)
		return
	}

	if seed >= 0 {
		common.Srand(seed)
		log.Printf("Using seed: %v", seed)
	} else if randomSeed {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		seed := rng.Int63()
		common.Srand(seed)
		log.Printf("Using seed: %v", seed)
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

	channel := make(chan signals.Interface, 128)

	opts := []func(*midisynth.MidiSynth){
		midisynth.WithChannel(channel),
	}

	if udpPort > 0 {
		udpInput := udp.New(udpPort, scale)
		opts = append(opts, midisynth.WithSignalInput(udpInput))
	}

	if midiPort > 0 {
		midiInput := midi.NewInput(midiPort, scale)
		opts = append(opts, midisynth.WithSignalInput(midiInput))
	}

	// Instruments
	instSet := instruments.NewSet()

	wavSettings := wav.Default

	// Pulseaudio player
	player, err := pulse.New(wavSettings.SampleRate, wavSettings.ChannelNum, wavSettings.BitDepthInBytes)
	if err != nil {
		log.Fatalf("pulse.New failed: %v", err)
	}

	// Audio output
	audioOpts := []func(*unisynth.Output){
		unisynth.WithInstruments(instSet),
		unisynth.WithScale(scale),
		unisynth.WithTempo(tempo),
		unisynth.WithWavSettings(wavSettings),
		unisynth.WithPlayer(player),
		unisynth.WithDelay(20 * time.Millisecond),
	}

	if fileInput == "" && flag.NArg() > 0 {
		fileInput = flag.Arg(0)
	}

	if fileInput != "" && dumpWav {
		fileOutput := strings.TrimSuffix(fileInput, filepath.Ext(fileInput)) + ".wav"
		audioOpts = append(audioOpts,
			unisynth.WithWavFileDump(fileOutput),
			unisynth.WithWavSpaceLeft(wavSpaceLeft),
			unisynth.WithWavSpaceRight(wavSpaceRight),
		)
	}

	audioOutput, err := unisynth.New(audioOpts...)
	check("Syntheziser output", err)
	opts = append(opts, midisynth.WithSignalOutput(audioOutput))
	// .

	// File sequencer
	var sequencer *seq.Sequencer
	if fileInput != "" {
		sequencer = seq.NewSequencer(
			seq.WithTempo(tempo),
			seq.WithStart(startBit),
			seq.WithShowingBits(showBits),
			seq.WithChannel(channel),
		)
		ps := parser.New(
			parser.WithSeq(sequencer),
			parser.WithScale(scale),
			parser.WithFileInput(fileInput),
			parser.WithTempoSetter(sequencer),
			parser.WithInstrumentSet(instSet),
		)
		check("Parser start", ps.Start(true))
		opts = append(opts, midisynth.WithSignalInput(sequencer))
	}

	if dump != "" {
		dumpOutput, err := dumper.NewJson(dump)
		check("Json dumper creation", err)
		opts = append(opts, midisynth.WithSignalOutput(dumpOutput))
	}

	if useConfig {
		cfg := config.New(getConfigPath(), instSet)
		check("MidiSynth initialization", cfg.InitMidiSynth())
		check("Config StartUpdateLoop", cfg.StartUpdateLoop())
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
