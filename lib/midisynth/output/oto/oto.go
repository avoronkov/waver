package oto

import (
	"io"

	"github.com/hajimehoshi/oto/v2"
)

type Oto struct {
	context *oto.Context
	// player  oto.Player
}

func New(sampleRate, channels, bits int) (*Oto, error) {

	c, ready, err := oto.NewContext(
		sampleRate,
		channels,
		bits,
	)
	if err != nil {
		return nil, err
	}
	<-ready
	return &Oto{
		context: c,
	}, nil
}

func (*Oto) Close() error {
	return nil
}

func (o *Oto) Play(r io.Reader) error {
	player := o.context.NewPlayer(r)
	player.Play()
	return nil
}
