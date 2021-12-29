package wav

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
)

type WavFmt struct {
	CompressionCode          uint16
	NumberOfChannels         uint16
	SampleRate               uint32
	AvgBps                   uint32
	BlockAlign               uint16
	SignificantBitsPerSample uint16
	ExtraFormatBytes         uint16
}

func (f *WavFmt) String() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "Compression code: %v\n", f.CompressionCode)
	fmt.Fprintf(s, "Number of channels: %v\n", f.NumberOfChannels)
	fmt.Fprintf(s, "SampleRate: %v\n", f.SampleRate)
	fmt.Fprintf(s, "Average bits per second: %v\n", f.AvgBps)
	fmt.Fprintf(s, "Block align: %v\n", f.BlockAlign)
	fmt.Fprintf(s, "Significant bits per sample: %v\n", f.SignificantBitsPerSample)
	fmt.Fprintf(s, "Extra format bytes: %v\n", f.ExtraFormatBytes)
	return s.String()
}

func ParseWavFmt(buf []byte) (*WavFmt, error) {
	if len(buf) < 16 {
		return nil, errors.New(fmt.Errorf("fmt section chunch should be at least 16 bytes, actual: %v", len(buf)))
	}

	f := &WavFmt{}
	f.CompressionCode = binary.LittleEndian.Uint16(buf[0:])
	f.NumberOfChannels = binary.LittleEndian.Uint16(buf[2:])
	f.SampleRate = binary.LittleEndian.Uint32(buf[4:])
	f.AvgBps = binary.LittleEndian.Uint32(buf[8:])
	f.BlockAlign = binary.LittleEndian.Uint16(buf[12:])
	f.SignificantBitsPerSample = binary.LittleEndian.Uint16(buf[14:])
	if len(buf) > 16 {
		f.ExtraFormatBytes = binary.LittleEndian.Uint16(buf[16:])
	}

	return f, nil
}
