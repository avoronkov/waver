package envelops

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
	"waver/lib/midisynth/wav"
	"waver/lib/midisynth/waves"
)

type Duration struct {
	settings wav.Settings
	duration time.Duration
}

func NewDuration(settings wav.Settings, duration time.Duration) Envelop {
	return &Duration{
		settings: settings,
		duration: duration,
	}
}

const maxInt16Amp = (2 << 15) - 1

// TODO currently works with int16 only
func (d *Duration) Wrap(wave waves.Wave) io.Reader {
	buf := new(bytes.Buffer)
	dt := 1.0 / float64(d.settings.SampleRate)
	durSeconds := float64(d.duration) / float64(time.Second)
	for t := 0.0; t < durSeconds; t += dt {
		value := int16(maxInt16Amp * wave.Value(t))
		for ch := 0; ch < d.settings.ChannelNum; ch++ {
			binary.Write(buf, binary.LittleEndian, value)
		}
	}
	return buf
}
