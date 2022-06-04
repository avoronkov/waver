package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/avoronkov/waver/lib/midisynth/dumper"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/notes"
	wavfmt "github.com/avoronkov/waver/wav"
)

type WavGenerator struct {
	input    string
	output   string
	instSet  *instruments.Set
	settings *wav.Settings
	tempo    int
	scale    notes.Scale

	channels int

	storage *SamplesStorage

	inputFile *os.File
	decoder   *json.Decoder
}

func NewWavGenerator(opts ...func(*WavGenerator)) (*WavGenerator, error) {
	g := &WavGenerator{
		tempo:    120,
		channels: 2,
		storage:  new(SamplesStorage),
	}

	for _, opt := range opts {
		opt(g)
	}
	if g.input == "" {
		return nil, fmt.Errorf("Input file not specified")
	}
	if g.output == "" {
		return nil, fmt.Errorf("Output file not specified")
	}
	if g.instSet == nil {
		return nil, fmt.Errorf("Instruments set not specified")
	}
	if g.settings == nil {
		g.settings = wav.Default
	}
	if g.scale == nil {
		g.scale = notes.NewStandard()
	}

	return g, nil
}

func (g *WavGenerator) Run() error {
	if err := g.openInput(); err != nil {
		return err
	}
	defer g.inputFile.Close()
	// Read samples one by one
	log.Printf("Processing input file...")
	for {
		var sig dumper.SignalJson
		err := g.decoder.Decode(&sig)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Cannot decode signal: %w", err)
		}
		if err := g.processSignal(&sig); err != nil {
			return fmt.Errorf("Cannot process signal '%v': %w", sig, err)
		}
	}
	// Normalize samples
	log.Printf("Normalizing samples...")
	g.storage.Normalize()
	// Convert to []int16
	log.Print("Performing descretization...")
	int16Samples := g.storage.ToInt16List()
	// Write samples to file
	log.Printf("Writing data into output file...")
	if err := g.writeSamplesIntoWavFile(int16Samples); err != nil {
		return err
	}
	log.Printf("Output file is generated: %v", g.output)
	return nil
}

func (g *WavGenerator) openInput() error {
	f, err := os.Open(g.input)
	if err != nil {
		return err
	}
	g.inputFile = f
	g.decoder = json.NewDecoder(g.inputFile)
	return nil
}

func (g *WavGenerator) processSignal(s *dumper.SignalJson) error {
	// 1. Find starting frame of the sample
	frame := int(float64(g.settings.SampleRate) * s.T)
	// 2. Generate wave for specified instrument/note
	wave, err := g.getWaveForSignal(s)
	if err != nil {
		return err
	}
	// 3. TODO Calculate the frequency
	freq := s.Note.Freq

	// 4. Calculate the duration
	if s.DurationBits == 0 {
		panic("s.DurationBits == 0")
	}
	dur := 15.0 * float64(s.DurationBits) / float64(g.tempo)
	// 5. Create note contexts
	noteCtx := waves.NewNoteCtx(freq, s.Amp, dur)
	// 6. While wave is not over store values into the storage

	tm := 0.0
	dt := 1.0 / float64(g.settings.SampleRate)
	for {
		value := wave.Value(tm, noteCtx)
		if math.IsNaN(value) {
			break
		}
		g.storage.AddSample(frame, value)

		tm += dt
		frame++
	}
	return nil
}

func (g *WavGenerator) getWaveForSignal(sig *dumper.SignalJson) (waves.Wave, error) {
	if sig.Manual || sig.Stop {
		return nil, fmt.Errorf("Manual controlled waves are not supported yet")
	}
	if sig.Sample != "" {
		wave, ok := g.instSet.Sample(sig.Sample)
		if !ok {
			return nil, fmt.Errorf("Sample not found: %v", sig.Sample)
		}
		return wave, nil
	}

	// Regular wave
	wave, ok := g.instSet.Wave(sig.Instrument)
	if !ok {
		return nil, fmt.Errorf("Instrument not found: %v", sig.Instrument)
	}
	return wave, nil
}

func (g *WavGenerator) writeSamplesIntoWavFile(samples []int16) error {
	w := wavfmt.CreateDefaultWav()
	w.Fmt.NumberOfChannels = uint16(g.channels)
	for _, sample := range samples {
		w.Data.AddSample(sample) // left
		if g.channels == 2 {
			w.Data.AddSample(sample) // right
		}
	}

	f, err := os.OpenFile(g.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := w.Write(f); err != nil {
		return err
	}
	return nil
}
