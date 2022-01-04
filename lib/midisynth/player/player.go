package player

import (
	"bytes"
	"encoding/binary"
	"io"
	"waver/lib/midisynth/wav"
	"waver/lib/midisynth/waves"
)

type Player struct {
	settings *wav.Settings
}

func New(settings *wav.Settings) *Player {
	return &Player{
		settings: settings,
	}
}

func (p *Player) Play(wave waves.WaveDuration) io.Reader {
	duration := wave.Duration()
	return p.PlayLimited(wave, duration)
}

const maxInt16Amp = (2 << 15) - 1

func (p *Player) PlayLimited(wave waves.Wave, duration float64) io.Reader {
	b := &bytes.Buffer{}
	dt := 1.0 / float64(p.settings.SampleRate)
	for t := 0.0; t < duration; t += dt {
		value := int16(maxInt16Amp * wave.Value(t))
		for ch := 0; ch < p.settings.ChannelNum; ch++ {
			binary.Write(b, binary.LittleEndian, value)
		}
	}
	return b
}
