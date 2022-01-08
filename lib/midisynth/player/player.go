package player

import (
	"bytes"
	"encoding/binary"
	"io"

	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Player struct {
	settings *wav.Settings
}

func New(settings *wav.Settings) *Player {
	return &Player{
		settings: settings,
	}
}

const maxInt16Amp = (1 << 15) - 1

func (p *Player) PlayContext(wave waves.Wave, ctx *waves.NoteCtx) (io.Reader, float64) {
	dur := 0.0
	if wd, ok := wave.(waves.WithDuration); ok {
		dur = wd.Duration(ctx)
	} else {
		dur = ctx.Dur
	}
	b := &bytes.Buffer{}
	dt := 1.0 / float64(p.settings.SampleRate)
	for t := 0.0; t < dur; t += dt {
		value := int16(maxInt16Amp * wave.Value(t, ctx))
		for ch := 0; ch < p.settings.ChannelNum; ch++ {
			binary.Write(b, binary.LittleEndian, value)
		}
	}
	return b, dur
}
