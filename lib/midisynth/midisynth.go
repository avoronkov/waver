package midisynth

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/avoronkov/waver/lib/midisynth/signals"
)

type MidiSynth struct {
	osSignals chan os.Signal
	ch        chan signals.Interface
	inputs    []signals.Input
	outputs   []signals.Output
	logf      func(format string, v ...any)
	closed    bool
}

func NewMidiSynth(opts ...func(*MidiSynth)) (*MidiSynth, error) {
	m := &MidiSynth{
		osSignals: make(chan os.Signal),
		logf:      log.Printf,
	}
	signal.Notify(m.osSignals, os.Interrupt)
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
		go func(in signals.Input) {
			defer func() {
				if r := recover(); r != nil {
					m.logf("Input recovered: %v\n%v", r, string(debug.Stack()))
				}
			}()
			if err := in.Run(m.ch); err != nil {
				m.logf("Input run failed: %v", err)
			}
		}(input)
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
			if _, ok := sig.(*signals.Stop); ok {
				log.Printf("Stopping...")
				m.Close()
				break L
			}
			tm := float64(time.Since(start)) / sec
			for _, output := range m.outputs {
				go output.ProcessAsync(tm, sig)
			}
		case <-m.osSignals:
			log.Printf("Interupting...")
			m.Close()
			break L
		}
	}
	return nil
}

func (m *MidiSynth) Close() error {
	// Handle possible double close
	if m.closed {
		return nil
	}
	m.closed = true
	for _, input := range m.inputs {
		if err := input.Close(); err != nil {
			log.Printf("Input Close() failed: %v", err)
		}
	}
	log.Printf("[MidiSynth] Closing %v outputs", len(m.outputs))
	for _, output := range m.outputs {
		if err := output.Close(); err != nil {
			log.Printf("Output Close() failed: %v", err)
		}
	}
	return nil
}
