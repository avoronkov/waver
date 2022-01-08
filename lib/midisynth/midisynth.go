package midisynth

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"waver/lib/midisynth/filters"
	instr "waver/lib/midisynth/instruments"
	"waver/lib/midisynth/player"
	"waver/lib/midisynth/wav"
	"waver/lib/midisynth/waves"
	waves2 "waver/lib/midisynth/waves/v2"
	"waver/lib/notes"
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
		settings: settings,
		play:     player.New(settings),
		context:  c,
		scale:    scale,
		port:     port,
		tempo:    120,
	}
	m.initInstruments()
	return m, nil
}

func (m *MidiSynth) initInstruments() {
	m.instruments = make(map[int]*instr.Instrument)
	m.instruments[1] = instr.NewInstrument(&waves2.Sine{}, filters.NewAdsrFilter())
	m.instruments[2] = instr.NewInstrument(&waves2.Square{}, filters.NewAdsrFilter())
	m.instruments[3] = instr.NewInstrument(&waves2.Triangle{}, filters.NewAdsrFilter())
	m.instruments[4] = instr.NewInstrument(&waves2.Saw{}, filters.NewAdsrFilter())
	m.instruments[5] = instr.NewInstrument(&waves2.Sine{}, filters.NewDistortionFilter(1.5), filters.NewAdsrFilter())
	m.instruments[6] = instr.NewInstrument(
		&waves2.Sine{},
		filters.NewAdsrFilter(),
		filters.NewDelayFilter(filters.DelayInterval(0.5), filters.DelayFadeOut(0.5), filters.DelayTimes(2)),
	)
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
		log.Printf("UDP: '%s'", buff)
		go m.handleMessage(buff[:n])
	}
}

// Intended to run in separate goroutine
func (m *MidiSynth) handleMessage(msg []byte) {
	if len(msg) < 3 {
		return
	}
	inst := int(msg[0] - '0')
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
	m.playNoteV2(inst, freq, dur, amp)
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

func (m *MidiSynth) playNoteV2(inst int, hz float64, dur float64, amp float64) {
	in, ok := m.instruments[inst]
	if !ok {
		log.Printf("Unknown instrument: %v", inst)
		return
	}

	data, duration := m.play.PlayContext(in.Wave(), waves2.NewNoteCtx(hz, amp, dur))

	p := m.context.NewPlayer(data)
	p.Play()

	time.Sleep(time.Duration(duration * float64(time.Second)))
	runtime.KeepAlive(p)
}

func (m *MidiSynth) playNote(inst int, hz float64, dur float64, amp float64) {
	var wave waves.Wave
	switch inst {
	case 1:
		wave = waves.Sine(hz)
	case 2:
		wave = waves.Square(hz)
	case 3:
		wave = waves.Triangle(hz)
	case 4:
		wave = waves.Saw(hz, false)
	default:
		log.Printf("Unknown instrument: %v", inst)
		return
	}
	p := m.context.NewPlayer(
		m.play.Play(
			filters.NewAdsr(
				wave,
				filters.AdsrAttackLevel(amp),
				filters.AdsrDecayLevel(amp),
				filters.AdsrReleaseLen(dur),
			),
		),
	)
	p.Play()
	time.Sleep(time.Duration(dur * float64(time.Second)))
	runtime.KeepAlive(p)
}

func (m *MidiSynth) Close() error {
	return nil
}
