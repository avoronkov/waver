package wav

import (
	"encoding/binary"
	"fmt"
	"io"
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

func ReadWavFmt(in io.Reader) (*WavFmt, error) {
	_, err := readBytesExpect(in, []byte("fmt "))
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt ' section id: %w", err)
	}

	// Read 'fmt ' section size - 4 bytes as INT
	fmtSize, err := readUint32(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt ' Chunk Data Size: %w", err)
	}
	fmtChunk, err := readBytes(in, int(fmtSize))
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt' chunk: %w", err)
	}
	return ParseWavFmt(fmtChunk)
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

func (f *WavFmt) Write(w io.Writer) error {
	io.WriteString(w, "fmt ")
	binary.Write(w, binary.LittleEndian, f.ChunkSize())
	binary.Write(w, binary.LittleEndian, f.CompressionCode)
	binary.Write(w, binary.LittleEndian, f.NumberOfChannels)
	binary.Write(w, binary.LittleEndian, f.SampleRate)
	binary.Write(w, binary.LittleEndian, f.AvgBps)
	binary.Write(w, binary.LittleEndian, f.BlockAlign)
	binary.Write(w, binary.LittleEndian, f.SignificantBitsPerSample)
	if f.ExtraFormatBytes != 0 {
		return errors.New("Writing Extra Format Bytes is not supported")
	}
	return nil
}

func (f *WavFmt) FullSize() uint32 {
	return uint32(len("fmt ")) + /* section length */ 4 + f.ChunkSize()
}

func (f *WavFmt) ChunkSize() uint32 {
	var size uint32 = 16
	if f.ExtraFormatBytes != 0 {
		size += 2 + uint32(f.ExtraFormatBytes)
	}
	return size
}
