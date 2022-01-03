package midisynth

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"waver/lib/midisynth/waves"
	"waver/lib/notes"
)

type MidiSynth struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int

	context *oto.Context

	p oto.Player

	scale notes.Scale

	port int

	tempo int
}

func NewMidiSynth(sampleRate int, channelNum int, bitDepthInBytes int, scale notes.Scale, port int) (*MidiSynth, error) {
	c, ready, err := oto.NewContext(sampleRate, channelNum, bitDepthInBytes)
	if err != nil {
		return nil, err
	}
	<-ready

	return &MidiSynth{
		sampleRate:      sampleRate,
		channelNum:      channelNum,
		bitDepthInBytes: bitDepthInBytes,
		context:         c,
		scale:           scale,
		port:            port,
		tempo:           120,
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
	inst := msg[0]
	_ = inst
	octave := int(msg[1] - '0')
	_ = octave
	note := string(msg[2])
	_ = note
	amp := 0.5
	if len(msg) >= 4 {
		amp = 0.1 * float64(m.parseValue(msg[3]))
	}
	dur := 500 * time.Millisecond
	log.Printf("len('%s') = %v", msg, len(msg))
	if len(msg) >= 5 {
		// Evaluate duration in bits (1/4 tempo)
		dur = time.Duration(15.0 * float64(time.Second) * float64(m.parseValue(msg[4])) / float64(m.tempo))
		log.Printf("dur = %v", dur)
	}

	freq, ok := m.scale.Note(octave, note)
	if !ok {
		log.Printf("Unknown note: %v%v", octave, note)
		return
	}
	m.playNote(freq, dur, amp)
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

func (m *MidiSynth) playNote(hz float64, dur time.Duration, amp float64) {
	p := m.context.NewPlayer(waves.NewSine(m.sampleRate, m.channelNum, m.bitDepthInBytes, hz, amp, dur))
	p.Play()
	time.Sleep(dur)
	runtime.KeepAlive(p)
}

func (m *MidiSynth) Close() error {
	return nil
}
