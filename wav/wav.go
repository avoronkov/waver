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

type DataInterface interface {
	Write(io.Writer) error
	FullSize() uint32
}

type Wav struct {
	Riff          []byte
	RiffChunkSize uint32
	RiffId        []byte

	Fmt *WavFmt

	DataId []byte
	Data   DataInterface
}

func (w *Wav) String() string {
	s := &strings.Builder{}
	fmt.Fprintf(s, "%s (size: %v)\n", w.Riff, w.RiffChunkSize)
	fmt.Fprintf(s, "%s\n", w.RiffId)
	if w.Fmt != nil {
		fmt.Fprintf(s, "%v", w.Fmt)
	}
	if w.Data != nil {
		fmt.Fprintf(s, "%v\n", w.Data)
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
	if err != nil {
		return nil, fmt.Errorf("Error reading 'fmt ' section data: %w", err)
	}

	wavFmt, err := ParseWavFmt(fmtChunk)
	if err != nil {
		return nil, fmt.Errorf("Error pasing 'fmt ' section data: %w", err)
	}

	wav.Fmt = wavFmt

	// Read 'data' section
	_, err = readBytesExpect(in, []byte("data"))
	if err != nil {
		return nil, fmt.Errorf("Error reading 'data' section id: %w", err)
	}
	dataSize, err := readUint32(in)
	if err != nil {
		return nil, fmt.Errorf("Error reading 'data' section seize: %w", err)
	}

	log.Printf("Data size: %v", dataSize)
	dataChunk, err := readBytes(in, int(dataSize))
	if err != nil {
		return nil, fmt.Errorf("Error reading 'data' chunk: %w", err)
	}
	data := ParseWavData(dataChunk)

	wav.Data = data

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
	if !bytes.Equal(buffer, expect) {
		return nil, errors.New(fmt.Errorf("Incorrect bytes read: expected '%s', actual '%s'", expect, buffer))
	}
	return buffer, nil
}

func (w *Wav) Write(out io.Writer) error {
	_, _ = io.WriteString(out, "RIFF")
	var chunkSize uint32 = uint32(len("WAVE")) + w.Fmt.FullSize() + w.Data.FullSize()
	_ = binary.Write(out, binary.LittleEndian, chunkSize)
	_, _ = io.WriteString(out, "WAVE")
	if err := w.Fmt.Write(out); err != nil {
		return fmt.Errorf("Failed to write 'fmt' chunk: %w", err)
	}
	if err := w.Data.Write(out); err != nil {
		return fmt.Errorf("Failed to write 'data' chunk: %w", err)
	}
	return nil
}
