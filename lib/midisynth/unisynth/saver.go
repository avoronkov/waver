package unisynth

import (
	"bytes"
	"io"
	"os"

	"github.com/avoronkov/waver/wav"
)

type WavDataSaver struct {
	orig io.Reader
	file string

	buffer bytes.Buffer
}

func NewWavDataSaver(r io.Reader, file string) *WavDataSaver {
	return &WavDataSaver{
		orig: r,
		file: file,
	}
}

func (s *WavDataSaver) Read(data []byte) (int, error) {
	n, err := s.orig.Read(data)
	if err == nil {
		_, _ = s.buffer.Write(data[:n])
	}
	return n, err
}

func (s *WavDataSaver) Close() error {
	w := wav.CreateDefaultWav()
	data := trimSamples(s.buffer.Bytes())
	w.Data = &wav.DataBytes{Samples: data}

	f, err := os.Create(s.file)
	if err != nil {
		return err
	}
	defer f.Close()

	return w.Write(f)
}

func trimSamples(data []byte) []byte {
	i := 0
	l := len(data)
	for i+3 < l {
		if !bytes.Equal(data[i:i+4], []byte{0, 0, 0, 0}) {
			break
		}
		i += 4
	}
	return data[i:]
}
