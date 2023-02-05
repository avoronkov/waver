package unisynth

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/avoronkov/waver/wav"
)

type WavDataSaver struct {
	orig io.Reader
	file string

	emptyFramesLeft, emptyFramesRight int

	buffer bytes.Buffer
}

func NewWavDataSaver(r io.Reader, file string, emptyFramesLeft, emptyFramesRight int) *WavDataSaver {
	return &WavDataSaver{
		orig:             r,
		file:             file,
		emptyFramesLeft:  emptyFramesLeft,
		emptyFramesRight: emptyFramesRight,
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
	if s.emptyFramesLeft > 0 {
		left := make([]byte, 4*s.emptyFramesLeft)
		data = append(left, data...)
	}
	if s.emptyFramesRight > 0 {
		right := make([]byte, 4*s.emptyFramesRight)
		data = append(data, right...)
	}
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
	j := l
	for i+3 < l {
		if !bytes.Equal(data[i:i+4], []byte{0, 0, 0, 0}) {
			break
		}
		i += 4
	}
	if l%4 != 0 {
		log.Printf("[WavDataSaver] Warning: unexpected data length: %v", l)
	} else {
		j -= 4
		for j >= 0 {
			if !bytes.Equal(data[j:j+4], []byte{0, 0, 0, 0}) {
				break
			}
			j -= 4
		}
	}
	log.Printf("[WavDataSaver] Trimming [0:%v] -> [%v:%v] (left: %v, right: %v)", l, i, j, i, l-j)
	return data[i:j]
}
