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

func (p *Player) PlayContext2(wave waves.Wave, ctx *waves.NoteCtx) (io.Reader, <-chan struct{}) {
	/*
		dur := 0.0
		if wd, ok := wave.(waves.WithDuration); ok {
			dur = wd.Duration(ctx)
		} else {
			dur = ctx.Dur
		}
	*/
	done := make(chan struct{})
	return &playerImpl{
		settings: p.settings,
		wave:     wave,
		ctx:      ctx,
		tm:       0.0,
		dt:       p.dt,
		dur:      math.Inf(1),
		done:     done,
	}, done
}

type playerImpl struct {
	settings *wav.Settings
	wave     waves.Wave
	ctx      *waves.NoteCtx
	tm       float64
	dt       float64
	dur      float64
	done     chan struct{}
	closed   bool
	eof      bool
}

var _ io.Reader = (*playerImpl)(nil)

func (pi *playerImpl) Read(data []byte) (n int, err error) {
	if pi.eof || pi.tm >= pi.dur {
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
	for buff.Len() < l && pi.tm < pi.dur {
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
	if (pi.eof || pi.tm >= pi.dur) && !pi.closed {
		close(pi.done)
		pi.closed = true
	}

	if buff.Len() > l {
		panic(fmt.Errorf("Buffer is to big: %v > %v", buff.Len(), l))
	}
	n = copy(data, buff.Bytes())
	return n, nil
}
