package waves

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/cryptix/wav"
)

type Sample struct {
	sampleRate float64

	data    []int16
	datalen int
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

	if nc := reader.GetNumChannels(); nc != 1 {
		return nil, fmt.Errorf("Only mono wav files supported (number of channels: %v)", nc)
	}

	if bits := reader.GetBitsPerSample(); bits != 16 {
		return nil, fmt.Errorf("Only 16bit samples are supported: %v", bits)
	}

	s := &Sample{
		sampleRate: float64(reader.GetSampleRate()),
	}

	sampleCount := int(reader.GetSampleCount())
	for i := 0; i < sampleCount; i++ {
		sample, err := reader.ReadRawSample()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ReadRawSample failed: %w", err)
		}
		if len(sample) != 2 {
			return nil, fmt.Errorf("Sample read size is not 16bit: %v", len(sample))
		}
		num := binary.LittleEndian.Uint16(sample)
		s.data = append(s.data, int16(num))
	}
	s.datalen = len(s.data)
	return s, nil
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
	reader, err := wav.NewReader(f, info.Size())
	if err != nil {
		return nil, fmt.Errorf("Failed to create wav.Reader: %w", err)
	}

	if nc := reader.GetNumChannels(); nc != 1 {
		return nil, fmt.Errorf("Only mono wav files supported (number of channels: %v)", nc)
	}

	if bits := reader.GetBitsPerSample(); bits != 16 {
		return nil, fmt.Errorf("Only 16bit samples are supported: %v", bits)
	}

	s := &Sample{
		sampleRate: float64(reader.GetSampleRate()),
	}

	sampleCount := int(reader.GetSampleCount())
	for i := 0; i < sampleCount; i++ {
		sample, err := reader.ReadRawSample()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ReadRawSample failed: %w", err)
		}
		if len(sample) != 2 {
			return nil, fmt.Errorf("Sample read size is not 16bit: %v", len(sample))
		}
		num := binary.LittleEndian.Uint16(sample)
		s.data = append(s.data, int16(num))
	}
	s.datalen = len(s.data)
	return s, nil
}

const maxInt16Value = float64((1 << 15) + 1)

func (s *Sample) Value(t float64, ctx *NoteCtx) float64 {
	n := int(t * s.sampleRate)
	if n >= s.datalen {
		return math.NaN()
	}

	res := float64(s.data[n]) / maxInt16Value
	return res
}

func (s *Sample) TimeLimited() {}

func (s *Sample) Data() []float64 {
	res := make([]float64, len(s.data))
	for i, v := range s.data {
		res[i] = float64(v) / maxInt16Value
	}
	return res
}
