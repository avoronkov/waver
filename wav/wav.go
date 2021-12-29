package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/go-errors/errors"
)

type Wav struct {
	Riff          []byte
	RiffChunkSize uint32
	RiffId        []byte

	Fmt *WavFmt
}

func (w *Wav) String() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "%s (size: %v)\n", w.Riff, w.RiffChunkSize)
	fmt.Fprintf(s, "%s\n", w.RiffId)
	if w.Fmt != nil {
		fmt.Fprintf(s, "%v", w.Fmt)
	}
	return s.String()
}

func ReadWav(in io.Reader) (*Wav, error) {
	wav := &Wav{}

	// Read 'RIFF' section
	riff, err := readBytesExpect(in, []byte("RIFF"))
	if err != nil {
		return nil, fmt.Errorf("Incorrect RIFF section: %w", err)
	}

	wav.Riff = riff

	// Read 'Chunk Data Size' 4 bytes as INT
	size, err := readUint32(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading RIFF Chunk Data Size: %w", err)
	}
	wav.RiffChunkSize = size

	// Read Riff ID 'WAVE'
	wave, err := readBytesExpect(in, []byte("WAVE"))
	if err != nil {
		return nil, fmt.Errorf("Error reading Riff ID ('WAVE'): %w", err)
	}
	wav.RiffId = wave

	// Read 'fmt ' section
	_, err = readBytesExpect(in, []byte("fmt "))
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt ' section: %w", err)
	}
	// Read 'fmt ' section size - 4 bytes as INT
	fmtSize, err := readUint32(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading RIFF Chunk Data Size: %w", err)
	}

	log.Printf("fmt section size: %v", fmtSize)
	fmtChunk, err := readBytes(in, int(fmtSize))

	wavFmt, err := ParseWavFmt(fmtChunk)
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt ' section data: %w", err)
	}

	wav.Fmt = wavFmt

	return wav, nil
}

func readUint32(in io.Reader) (uint32, error) {
	buffer := make([]byte, 4)
	n, err := in.Read(buffer)
	if err != nil {
		return 0, errors.New(err)
	}
	if n != 4 {
		return 0, errors.New("Unexpected EOF reading int32")
	}
	data := binary.LittleEndian.Uint32(buffer)
	return data, nil
}

func readBytes(in io.Reader, size int) ([]byte, error) {
	buffer := make([]byte, size)
	n, err := in.Read(buffer)
	if err != nil {
		return nil, errors.New(err)
	}
	if n < size {
		return nil, errors.New("Unexpected EOF while reading bytes buffer")
	}
	return buffer, nil
}

func readBytesExpect(in io.Reader, expect []byte) ([]byte, error) {
	buffer, err := readBytes(in, len(expect))
	if err != nil {
		return nil, err
	}
	if bytes.Compare(buffer, expect) != 0 {
		return nil, errors.New(fmt.Errorf("Incorrect bytes read: expected '%s', actual '%s'", expect, buffer))
	}
	return buffer, nil
}
