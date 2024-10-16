package waves

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/cryptix/wav"
)

type Sample struct {
	sampleRate float64

	data     [][]int32
	datalen  int
	channels int
	maxValue float64
}

func ParseSample(data []byte) (*Sample, error) {
	f := bytes.NewReader(data)
	return parseSample(f, int64(len(data)))
}

func ParseSampleFile(path string) (*Sample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Opening file '%v' failed: %w", path, err)
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("Stat file '%v' failed: %w", path, err)
	}
	return parseSample(f, stat.Size())
}

func parseSample(f io.ReadSeeker, size int64) (*Sample, error) {
	reader, err := wav.NewReader(f, size)
	if err != nil {
		return nil, fmt.Errorf("Failed to create wav.Reader: %w", err)
	}

	nc := reader.GetNumChannels()

	nBits := int(reader.GetBitsPerSample())
	if nBits > 32 {
		return nil, fmt.Errorf("%b bit samples are not supported (too large)", nBits)
	}
	nBytes := nBits / 8

	s := &Sample{
		sampleRate: float64(reader.GetSampleRate()),
		data:       make([][]int32, nc),
		channels:   int(nc),
		maxValue:   float64((int32(1) << (nBits - 1)) + 1),
	}

	sampleCount := int(reader.GetSampleCount()) / s.channels
	for i := 0; i < sampleCount; i++ {
		for ch := 0; ch < s.channels; ch++ {
			sample, err := reader.ReadRawSample()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("ReadRawSample failed: %w", err)
			}
			if len(sample) != nBytes {
				return nil, fmt.Errorf("Sample read size is not %v bit: %v", nBytes, len(sample))
			}
			num := bytesToInt32(sample, nBits)
			s.data[ch] = append(s.data[ch], num)
		}
	}
	s.datalen = len(s.data[0])
	return s, nil
}

// See also: binary/encoding.LittleEndian.Uint16
func bytesToInt32(b []byte, nBits int) int32 {
	switch nBits {
	case 16:
		return int32(int16(uint16(b[0]) | uint16(b[1])<<8))
	case 24:
		n := int32(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16)
		if n >= (1 << 23) {
			n -= (1 << 24)
		}
		return n
	case 32:
		return int32(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	default:
		panic(fmt.Errorf("bit size not supported: %v", nBits))
	}
}

func ReadSample(file string) (*Sample, error) {
	info, err := os.Stat(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to stat %v: %w", file, err)
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %v: %w", file, err)
	}
	defer f.Close()
	return parseSample(f, info.Size())
}

func (s *Sample) Value(t float64, ctx *NoteCtx) float64 {
	n := int(t * s.sampleRate)
	if n >= s.datalen {
		return math.NaN()
	}

	res := float64(s.data[0][n]) / s.maxValue
	return res
}

func (s *Sample) TimeLimited() {}

func (s *Sample) Data(ch int) []float64 {
	res := make([]float64, len(s.data[ch]))
	for i, v := range s.data[ch] {
		res[i] = float64(v) / s.maxValue
	}
	return res
}
