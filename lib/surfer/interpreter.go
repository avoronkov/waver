package surfer

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/avoronkov/waver/lib/forth"
	"github.com/avoronkov/waver/lib/forth/parser"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/cryptix/wav"
)

type Interpreter struct {
	// Samples
	slices    [][]float64
	position  int
	slicesLen int

	// Forth
	forth *forth.Forth

	// Output
	output  []float64
	outFile string
}

func InitInterpreter(forthFile, wavFile, outFile string) (*Interpreter, error) {
	frt, err := parser.ParseFile(forthFile)
	if err != nil {
		return nil, err
	}

	sample, err := waves.ParseSampleFile(wavFile)
	if err != nil {
		return nil, fmt.Errorf("ParseSampleFile failed: %v", err)
	}

	slices := SlicesFromSamples(sample.Data())

	in := &Interpreter{
		slices:    slices,
		slicesLen: len(slices),
		forth:     frt,
		outFile:   outFile,
	}
	forth.WithFuncs(map[string]forth.StackFn{
		"Play": in.Play,
	})(in.forth)

	return in, nil
}

func (i *Interpreter) Run() error {
	err := i.forth.Run()
	if err != nil {
		return err
	}

	return i.saveWavFile()
}

func (i *Interpreter) Play(f *forth.Forth) error {
	if i.position >= i.slicesLen {
		f.Stack.Push(0)
		return nil
	}
	i.output = append(i.output, i.slices[i.position]...)
	i.position++
	f.Stack.Push(1)
	return nil
}

const maxInt16Value = float64((1 << 15) + 1)

func (i *Interpreter) saveWavFile() error {
	meta := wav.File{
		Channels:        1,
		SampleRate:      44100,
		SignificantBits: 16,
	}
	f, err := os.Create(i.outFile)
	if err != nil {
		return nil
	}
	defer f.Close()
	writer, err := meta.NewWriter(f)
	if err != nil {
		return err
	}

	for _, sample := range i.output {
		u := uint16(sample * maxInt16Value)
		bytes := []byte{0, 0}
		binary.LittleEndian.PutUint16(bytes, u)
		writer.WriteSample(bytes)
	}

	return writer.Close()
}
