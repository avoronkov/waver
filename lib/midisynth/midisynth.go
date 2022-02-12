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

	tempo int

	instruments map[int]*instr.Instrument
	samples     map[string]*instr.Instrument

	osSignals chan os.Signal
	ch        chan *signals.Signal
	inputs    []signals.Input
	outputs   []signals.Output

	// Octave -> Note -> Release fn()
	notesReleases map[int]map[string]func()

	edo int
}

func NewMidiSynth(opts ...func(*MidiSynth)) (*MidiSynth, error) {
	m := &MidiSynth{
		settings:      wav.Default,
		tempo:         120,
		instruments:   make(map[int]*instr.Instrument),
		samples:       make(map[string]*instr.Instrument),
		osSignals:     make(chan os.Signal),
		ch:            make(chan *signals.Signal),
		notesReleases: make(map[int]map[string]func()),
		edo:           12,
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

L:
	for {
		select {
		case sig := <-m.ch:
			for _, output := range m.outputs {
				go output.ProcessAsync(sig)
			}
			// go m.PlayNoteSignal(sig)
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
	} else if !s.Manual {
		// Play note
		err = m.PlayNote(s.Instrument, s.Octave, s.Note, s.DurationBits, s.Amp)
	} else if s.Stop {
		// Stop manual note
		m.releaseNote(s.Octave, s.Note)
	} else {
		// Play manual note
		stop, err := m.PlayNoteControlled(
			s.Instrument,
			s.Octave,
			s.Note,
			s.Amp,
		)
		if err != nil {
			log.Printf("[Manual] error: %v", err)
			return
		}
		m.storeNoteReleaseFn(s.Octave, s.Note, stop)
	}
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func (m *MidiSynth) storeNoteReleaseFn(octave int, note string, release func()) {
	if _, ok := m.notesReleases[octave]; ok {
		m.notesReleases[octave][note] = release
	} else {
		m.notesReleases[octave] = map[string]func(){
			note: release,
		}
	}
}

func (m *MidiSynth) releaseNote(octave int, note string) {
	if notes, ok := m.notesReleases[octave]; ok {
		if release, ok := notes[note]; ok {
			release()
			delete(notes, note)
		}
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
	for _, input := range m.inputs {
		if err := input.Close(); err != nil {
			log.Printf("Input Close() failed: %v", err)
		}
	}
	return nil
}
