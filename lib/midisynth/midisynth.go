package midisynth

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
)

type MidiSynth struct {
	osSignals chan os.Signal
	ch        chan *signals.Signal
	inputs    []signals.Input
	outputs   []signals.Output
}

func NewMidiSynth(opts ...func(*MidiSynth)) (*MidiSynth, error) {
	m := &MidiSynth{
		osSignals: make(chan os.Signal),
		ch:        make(chan *signals.Signal, 128),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m, nil
}

func (m *MidiSynth) Start() error {
	if len(m.inputs) == 0 {
		return fmt.Errorf("No inputs specified")
	}
	if len(m.outputs) == 0 {
		return fmt.Errorf("No outputs specified")
	}

	for _, input := range m.inputs {
		if err := input.Start(m.ch); err != nil {
			return err
		}
	}

	start := time.Now()
	sec := float64(time.Second)

L:
	for {
		select {
		case sig := <-m.ch:
			if sig == nil {
				continue L
			}
			tm := float64(time.Since(start)) / sec
			for _, output := range m.outputs {
				go output.ProcessAsync(tm, sig)
			}
		case <-m.osSignals:
			log.Printf("Interupting...")
			break L
		}
	}
	return nil
}

func (m *MidiSynth) Close() error {
	for _, input := range m.inputs {
		if err := input.Close(); err != nil {
			log.Printf("Input Close() failed: %v", err)
		}
	}
	for _, output := range m.outputs {
		if err := output.Close(); err != nil {
			log.Printf("Output Close() failed: %v", err)
		}
	}
	return nil
}
