package midisynth

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"waver/lib/midisynth/filters"
	"waver/lib/midisynth/player"
	"waver/lib/midisynth/wav"
	"waver/lib/midisynth/waves"
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
}

func NewMidiSynth(settings *wav.Settings, scale notes.Scale, port int) (*MidiSynth, error) {
	c, ready, err := oto.NewContext(settings.SampleRate, settings.ChannelNum, settings.BitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready

	return &MidiSynth{
		settings: settings,
		play:     player.New(settings),
		context:  c,
		scale:    scale,
		port:     port,
		tempo:    120,
	}, nil
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
	m.playNote(inst, freq, dur, amp)
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
