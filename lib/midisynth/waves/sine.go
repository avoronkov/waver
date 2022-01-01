package waves

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"time"
)

type Sine struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int

	freq float64
	// amp [0, 1]
	amp float64
	// Duration in seconds
	duration time.Duration

	buf []byte
	pos int
}

func NewSine(
	sampleRate int,
	channelNum int,
	bitDepthInBytes int,
	freq float64,
	amp float64,
	duration time.Duration,
) *Sine {
	sine := &Sine{
		sampleRate:      sampleRate,
		channelNum:      channelNum,
		bitDepthInBytes: bitDepthInBytes,
		freq:            freq,
		amp:             amp,
		duration:        duration,
	}
	sine.init()
	return sine
}

func (s *Sine) init() {
	samples := int(float64(s.duration/time.Second) * float64(s.sampleRate))
	log.Printf("samples: %v", samples)
	buf := new(bytes.Buffer)
	for i := 0; i < samples; i++ {
		value := s.sineInt16(i)
		for ch := 0; ch < s.channelNum; ch++ {
			binary.Write(buf, binary.LittleEndian, value)
		}
	}
	s.buf = buf.Bytes()
	log.Printf("buffer size: %v", len(s.buf))
}

func (s *Sine) sineInt16(i int) int16 {
	// TODO optimization
	waveLength := float64(s.sampleRate) / s.freq
	if i == 0 {
		log.Printf("wave length: %v", waveLength)
	}
	x := 2.0 * math.Pi * float64(i) / waveLength
	amp := 32767.0 * s.amp
	return int16(amp * math.Sin(x))
}

func (s *Sine) String() string {
	return fmt.Sprintf("Sine: freq=%v, buffer length=%v", s.freq, len(s.buf))
}

func (s *Sine) Read(buf []byte) (n int, err error) {
	defer func(s *Sine) {
		l := len(s.buf)
		log.Printf("Read %v bytes: %v, %v; pos: %v; left: %v", l, n, err, s.pos, l-s.pos)
	}(s)
	if s.pos == len(s.buf) {
		return 0, io.EOF
	}
	n = copy(buf, s.buf[s.pos:])
	s.pos += n
	return n, nil
}
