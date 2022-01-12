package midisynth

import (
	"fmt"
	"log"
	"math"
	"net"
	"runtime"

	oto "github.com/hajimehoshi/oto/v2"

	instr "gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/player"
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

	port int

	tempo int

	instruments map[int]*instr.Instrument
}

func NewMidiSynth(settings *wav.Settings, scale notes.Scale, port int) (*MidiSynth, error) {
	c, ready, err := oto.NewContext(settings.SampleRate, settings.ChannelNum, settings.BitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready

	m := &MidiSynth{
		settings:    settings,
		play:        player.New(settings),
		context:     c,
		scale:       scale,
		port:        port,
		tempo:       120,
		instruments: make(map[int]*instr.Instrument),
	}
	return m, nil
}

func (m *MidiSynth) AddInstrument(n int, in *instr.Instrument) {
	m.instruments[n] = in
}

func (m *MidiSynth) Start() error {
	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%v", m.port))
	if err != nil {
		return fmt.Errorf("Starting UDP server failed: %w", err)
	}
	defer pc.Close()
	log.Printf("Listening to UDP on localhost:%v", m.port)
	for {
		buff := make([]byte, 64)
		n, _, err := pc.ReadFrom(buff)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			continue
		}
		go m.handleMessage(buff[:n])
	}
}

// Intended to run in separate goroutine
func (m *MidiSynth) handleMessage(msg []byte) {
	if len(msg) < 3 {
		return
	}
	inst := m.parseValue(msg[0])
	octave := int(msg[1] - '0')
	note := string(msg[2])
	amp := 0.5
	if len(msg) >= 4 {
		amp = 0.1 * float64(m.parseValue(msg[3]))
	}
	dur := 0.5
	if len(msg) >= 5 {
		// Evaluate duration in bits (1/4 tempo)
		dur = 15.0 * float64(m.parseValue(msg[4])) / float64(m.tempo)
		log.Printf("dur = %v", dur)
	}

	freq, ok := m.scale.Note(octave, note)
	if !ok {
		log.Printf("Unknown note: %v%v", octave, note)
		return
	}
	m.playNote(inst, freq, dur, amp)
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

func (m *MidiSynth) parseValue(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'a' && b <= 'z' {
		return 10 + int(b-'a')
	}
	if b >= 'A' && b <= 'Z' {
		return 10 + int(b-'Z')
	}
	return 0
}

func (m *MidiSynth) playNote(inst int, hz float64, dur float64, amp float64) {
	in, ok := m.instruments[inst]
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	data, done := m.play.PlayContext2(in.Wave(), waves.NewNoteCtx(hz, amp, dur))

	p := m.context.NewPlayer(data)
	p.Play()

	<-done
	runtime.KeepAlive(p)
}

func (m *MidiSynth) playNoteControlled(inst int, hz float64, amp float64) (stop func()) {
	in, ok := m.instruments[inst]
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	wave := in.WaveControlled()
	data, done := m.play.PlayContext2(wave, waves.NewNoteCtx(hz, amp, math.Inf(+1)))

	go func() {
		p := m.context.NewPlayer(data)
		p.Play()

		<-done
		runtime.KeepAlive(p)

	}()
	return wave.Release
}

func (m *MidiSynth) Close() error {
	return nil
}
