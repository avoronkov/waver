package midisynth

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	instr "gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/midi"
	"gitlab.com/avoronkov/waver/lib/midisynth/player"
	"gitlab.com/avoronkov/waver/lib/midisynth/signals"
	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
	"gitlab.com/avoronkov/waver/lib/notes"
)

type MidiSynth struct {
	settings *wav.Settings

	play *player.Player

	context *oto.Context

	p oto.Player

	scale notes.Scale

	// Midi port for controller client
	midiProc *midi.Proc

	tempo int

	instruments map[int]*instr.Instrument
	samples     map[string]*instr.Instrument

	midiChan  chan string
	osSignals chan os.Signal
	ch        chan *signals.Signal
	inputs    []signals.Input

	edo int
}

func NewMidiSynth(opts ...func(*MidiSynth)) (*MidiSynth, error) {
	m := &MidiSynth{
		settings: wav.Default,
		// play:        player.New(settings),
		// context:     c,
		// scale:       scale,
		tempo:       120,
		instruments: make(map[int]*instr.Instrument),
		samples:     make(map[string]*instr.Instrument),
		midiChan:    make(chan string),
		osSignals:   make(chan os.Signal),
		ch:          make(chan *signals.Signal),
		edo:         12,
	}
	for _, opt := range opts {
		opt(m)
	}

	// Init scale
	if m.scale == nil {
		m.scale = notes.NewStandard()
	}

	// Init oto.Context
	c, ready, err := oto.NewContext(m.settings.SampleRate, m.settings.ChannelNum, m.settings.BitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready

	// Init Player
	m.play = player.New(m.settings)

	m.context = c
	return m, nil
}

func (m *MidiSynth) AddInstrument(n int, in *instr.Instrument) {
	m.instruments[n] = in
}

func (m *MidiSynth) AddSampledInstrument(name string, in *instr.Instrument) {
	m.samples[name] = in
}

func (m *MidiSynth) Start() error {
	started := false
	for _, input := range m.inputs {
		started = true
		if err := input.Start(m.ch); err != nil {
			return err
		}
	}
	if m.midiProc != nil {
		started = true
		if err := m.midiProc.Start(); err != nil {
			return err
		}
	}

	if !started {
		return fmt.Errorf("Either UDP or MIDI port should be specifiend")
	}
L:
	for {
		select {
		case sig := <-m.ch:
			// TODO make stops
			go m.PlayNoteSignal(sig)
		case data := <-m.midiChan:
			go m.midiProc.HandleLine(data)
		case <-m.osSignals:
			log.Printf("Interupting...")
			break L
		}
	}
	return nil
}

func (m *MidiSynth) PlayNoteSignal(s *signals.Signal) {
	var err error
	if s.Sample != "" {
		// Play sample
		dur := 15.0 * float64(s.DurationBits) / float64(m.tempo)
		err = m.PlaySample(s.Sample, dur, s.Amp)
	} else {
		// Play note
		err = m.PlayNote(s.Instrument, s.Octave, s.Note, s.DurationBits, s.Amp)
	}
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func (m *MidiSynth) PlayNote(instr int, octave int, note string, durationBits int, amp float64) error {
	freq, ok := m.scale.Note(octave, note)
	if !ok {
		return fmt.Errorf("Unknown note: %v%v", octave, note)
	}
	dur := 15.0 * float64(durationBits) / float64(m.tempo)
	m.playNote(instr, freq, dur, amp)
	return nil
}

func (m *MidiSynth) PlayNoteControlled(instr int, octave int, note string, amp float64) (stop func(), err error) {
	freq, ok := m.scale.Note(octave, note)
	if !ok {
		return nil, fmt.Errorf("Unknown note: %v%v", octave, note)
	}

	stop = m.playNoteControlled(instr, freq, amp)

	return stop, nil
}

func (m *MidiSynth) PlaySample(name string, duration float64, amp float64) error {
	in, ok := m.samples[name]
	if !ok {
		return fmt.Errorf("Unknown sample: %q", name)
	}
	data, done := m.play.PlayContext(in.Wave(), waves.NewNoteCtx(0, amp, duration))

	p := m.context.NewPlayer(data)
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
	return nil
}

func (m *MidiSynth) playNote(inst int, hz float64, dur float64, amp float64) {
	in, ok := m.instruments[inst]
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	data, done := m.play.PlayContext(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

	p := m.context.NewPlayer(data)
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
}

func (m *MidiSynth) playNoteControlled(inst int, hz float64, amp float64) (stop func()) {
	in, ok := m.instruments[inst]
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	wave := in.WaveControlled()
	data, done := m.play.PlayContext(wave, waves.NewNoteCtx(hz, amp, math.Inf(+1)))

	go func() {
		p := m.context.NewPlayer(data)
		p.Play()

		<-done
		time.Sleep(1 * time.Second)
		runtime.KeepAlive(p)

	}()
	return wave.Release
}

func (m *MidiSynth) Close() error {
	_ = m.midiProc.Close()
	for _, input := range m.inputs {
		if err := input.Close(); err != nil {
			log.Printf("Input Close() failed: %v", err)
		}
	}
	return nil
}
