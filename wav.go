package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-errors/errors"
)

type Wav struct {
	Riff []byte
}

func (w *Wav) String() string {
	return fmt.Sprintf("``%s``", w.Riff)
}

func ReadWav(in io.Reader) (*Wav, error) {
	// Read 4 bytes
	buffer := make([]byte, 4)
	n, err := in.Read(buffer)
	if err != nil {
		return nil, errors.New(err)
	}
	if n != 4 {
		return nil, errors.New("Unexpected EOF while reading 'RIFF' section")
	}
	if bytes.Compare(buffer, []byte("RIFF")) != 0 {
		return nil, errors.New(fmt.Errorf("Incorrect RIFF chunk: expected 'RIFF', got '%s'", buffer))
	}

	wav := &Wav{}
	wav.Riff = buffer
	// TODO
	return wav, nil
}
