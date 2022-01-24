package player

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"gitlab.com/avoronkov/waver/lib/midisynth/wav"
	"gitlab.com/avoronkov/waver/lib/midisynth/waves"
)

type Player struct {
	settings *wav.Settings
	dt       float64
}

func New(settings *wav.Settings) *Player {
	return &Player{
		settings: settings,
		dt:       1.0 / float64(settings.SampleRate),
	}
}

const maxInt16Amp = (1 << 15) - 1

func (p *Player) PlayContext(wave waves.Wave, ctx *waves.NoteCtx) (io.Reader, <-chan struct{}) {
	done := make(chan struct{})
	return &playerImpl{
		settings: p.settings,
		wave:     wave,
		ctx:      ctx,
		tm:       0.0,
		dt:       p.dt,
		done:     done,
	}, done
}

type playerImpl struct {
	settings *wav.Settings
	wave     waves.Wave
	ctx      *waves.NoteCtx
	tm       float64
	dt       float64
	done     chan struct{}
	closed   bool
	eof      bool
}

var _ io.Reader = (*playerImpl)(nil)

func (pi *playerImpl) Read(data []byte) (n int, err error) {
	if pi.eof {
		if !pi.closed {
			close(pi.done)
			pi.closed = true
		}
		return 0, io.EOF
	}

	l := 640
	if len(data) < l {
		l = len(data)
	}
	buff := new(bytes.Buffer)
	for buff.Len() < l {
		waveValue := pi.wave.Value(pi.tm, pi.ctx)
		if math.IsNaN(waveValue) {
			pi.eof = true
			break
		}
		value := int16(maxInt16Amp * pi.wave.Value(pi.tm, pi.ctx))
		for ch := 0; ch < pi.settings.ChannelNum; ch++ {
			binary.Write(buff, binary.LittleEndian, value)
		}
		pi.tm += pi.dt
	}
	if (pi.eof) && !pi.closed {
		close(pi.done)
		pi.closed = true
	}

	if buff.Len() > l {
		panic(fmt.Errorf("Buffer is to big: %v > %v", buff.Len(), l))
	}
	n = copy(data, buff.Bytes())
	return n, nil
}
