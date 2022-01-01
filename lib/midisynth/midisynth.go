package midisynth

import (
	"runtime"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"waver/lib/midisynth/waves"
)

type MidiSynth struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int

	context *oto.Context

	p oto.Player
}

func NewMidiSynth(sampleRate int, channelNum int, bitDepthInBytes int) (*MidiSynth, error) {
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
	}, nil
}

func (m *MidiSynth) Start() {
	dur := 16 * time.Second
	m.p = m.context.NewPlayer(waves.NewSine(m.sampleRate, m.channelNum, m.bitDepthInBytes, 440.0, 0.2, dur))
	m.p.Play()
	players := []oto.Player{m.p}
	time.Sleep(dur)
	runtime.KeepAlive(players)
}

func (m *MidiSynth) Close() error {
	return nil
}
