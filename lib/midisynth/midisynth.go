package midisynth

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	instr "gitlab.com/avoronkov/waver/lib/midisynth/instruments"
	"gitlab.com/avoronkov/waver/lib/midisynth/midi"
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

	// UDP udpPort for Orca messages
	udpPort int

	// Midi port for controller client
	midiPort int
	midiProc *midi.Proc

	tempo int

	instruments map[int]*instr.Instrument
	samples     map[string]*instr.Instrument

	midiChan chan string
	udpChan  chan []byte
	signals  chan os.Signal

	udpListener net.PacketConn
}

func NewMidiSynth(opts ...func(*MidiSynth)) (*MidiSynth, error) {
	m := &MidiSynth{
		settings: wav.Default,
		// play:        player.New(settings),
		// context:     c,
		// scale:       scale,
		// udpPort:     port,
		tempo:       120,
		instruments: make(map[int]*instr.Instrument),
		samples:     make(map[string]*instr.Instrument),
		midiChan:    make(chan string),
		udpChan:     make(chan []byte),
		signals:     make(chan os.Signal),
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
	if m.udpPort > 0 {
		started = true
		if err := m.startUdpListener(m.udpPort, m.udpChan); err != nil {
			return err
		}
	}
	if m.midiPort > 0 {
		log.Printf("Starting MIDI listener on port %v", m.midiPort)
		started = true
		m.midiProc = midi.NewProc(m, m.midiPort, m.midiChan)
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
		case msg := <-m.udpChan:
			go m.handleMessage(msg)
		case data := <-m.midiChan:
			go m.midiProc.HandleLine(data)
		case <-m.signals:
			log.Printf("Interupting...")
			break L
		}
	}
	return nil
}

func (m *MidiSynth) startUdpListener(port int, ch chan []byte) (err error) {
	m.udpListener, err = net.ListenPacket("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		return fmt.Errorf("Starting UDP server failed: %w", err)
	}
	log.Printf("Listening to UDP on localhost:%v", port)
	go func(pc net.PacketConn) {
	L:
		for {
			buff := make([]byte, 64)
			n, _, err := pc.ReadFrom(buff)
			if errors.Is(err, net.ErrClosed) {
				log.Printf("[ERROR] ErrClosed")
				break L
			}
			if err != nil {
				log.Printf("[ERROR] %v (%T, %v)", err, err, err.(*net.OpError).Unwrap())
				continue
			}
			ch <- buff[:n]
		}
	}(m.udpListener)
	return nil
}

// Intended to run in separate goroutine
func (m *MidiSynth) handleMessage(msg []byte) {
	if len(msg) < 3 {
		return
	}
	log.Printf("Handle UDP message: %s", msg)
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
	}

	if inst == 35 { // 'z'
		m.PlaySample(string(msg[1:3]), dur, amp)
		return
	}

	freq, ok := m.scale.Note(octave, note)
	if !ok {
		log.Printf("Unknown note: %v %v (%s)", octave, note, msg)
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

func (m *MidiSynth) PlaySample(name string, duration float64, amp float64) error {
	in, ok := m.samples[name]
	if !ok {
		return fmt.Errorf("Unknown sample: %q", name)
	}
	data, done := m.play.PlayContext2(in.Wave(), waves.NewNoteCtx(0, amp, duration))

	p := m.context.NewPlayer(data)
	p.Play()

	<-done
	time.Sleep(1 * time.Second)
	runtime.KeepAlive(p)
	return nil
}

func (m *MidiSynth) parseValue(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'a' && b <= 'z' {
		return 10 + int(b-'a')
	}
	if b >= 'A' && b <= 'Z' {
		return 10 + int(b-'A')
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
	data, done := m.play.PlayContext2(wave, waves.NewNoteCtx(hz, amp, math.Inf(+1)))

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
	_ = m.udpListener.Close()
	_ = m.midiProc.Close()
	return nil
}
