package surfer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/avoronkov/waver/lib/forth"
	"github.com/avoronkov/waver/lib/forth/parser"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/utils"
	"github.com/avoronkov/waver/wav"
)

type Interpreter struct {
	// Samples
	slices    [][]float64
	position  int
	slicesLen int

	// Forth
	forth *forth.Forth

	// Output
	outputLeft  []float64
	outputRight []float64

	levelLeft  float64
	levelRight float64

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
		slices:     slices,
		slicesLen:  len(slices),
		forth:      frt,
		outFile:    outFile,
		levelLeft:  1.0,
		levelRight: 1.0,
	}
	forth.WithFuncs(map[string]forth.StackFn{
		"Play":     in.Play,
		"|>":       in.Play,
		"PlayBack": in.PlayBack,
		"<|":       in.PlayBack,
		"NPlay":    in.NPlay,
		"FF":       in.FastForward,
		"Pos":      in.Position,
		"Len":      in.Length,
		"Goto":     in.Goto,
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

// "1" if Played successfully
// "0" if EOF
func (i *Interpreter) Play(f *forth.Forth) error {
	if i.position >= i.slicesLen {
		f.Stack.Push(0)
		return nil
	}
	i.writeOutput(i.slices[i.position]...)
	i.position++
	f.Stack.Push(1)
	return nil
}

func (i *Interpreter) writeOutput(values ...float64) {
	valuesLeft := utils.SliceMap(values, func(x float64) float64 {
		return x * i.levelLeft
	})
	i.outputLeft = append(i.outputLeft, valuesLeft...)

	valuesRight := utils.SliceMap(values, func(x float64) float64 {
		return x * i.levelRight
	})
	i.outputRight = append(i.outputRight, valuesRight...)
}

func (i *Interpreter) NPlay(f *forth.Forth) error {
	n, err := forth.Pop[int](f.Stack)
	if err != nil {
		return err
	}
	if n > 0 {
		cnt := 0
		for j := 0; j < n; j++ {
			if i.position >= i.slicesLen {
				break
			}
			i.writeOutput(i.slices[i.position]...)
			cnt++
			i.position++
		}
		log.Printf("NPlay returns: %v", cnt)
		f.Stack.Push(cnt)
	} else if n < 0 {
		cnt := 0
		for j := 0; j > n; j-- {
			if i.position <= 0 {
				break
			}
			i.position--
			slice := i.slices[i.position]
			for k := len(slice) - 1; k >= 0; k-- {
				i.writeOutput(-slice[k])
			}
			cnt--
		}
		f.Stack.Push(cnt)

	} else {
		f.Stack.Push(0)
	}
	return nil
}

func (i *Interpreter) PlayBack(f *forth.Forth) error {
	if i.position <= 0 {
		f.Stack.Push(0)
		return nil
	}
	i.position--
	slice := i.slices[i.position]
	for j := len(slice) - 1; j >= 0; j-- {
		i.writeOutput(-slice[j])
	}
	f.Stack.Push(-1)
	return nil
}

// Returns actual "shift"
func (i *Interpreter) FastForward(f *forth.Forth) error {
	shift, err := forth.Pop[int](f.Stack)
	if err != nil {
		return err
	}
	newPos := i.position + shift
	if newPos > i.slicesLen {
		newPos = i.slicesLen
	} else if newPos < 0 {
		newPos = 0
	}
	actualShift := newPos - i.position
	log.Printf("Actual shift (%v): %v -> %v == %v", shift, i.position, newPos, actualShift)
	i.position = newPos
	f.Stack.Push(actualShift)
	return nil
}

func (i *Interpreter) Goto(f *forth.Forth) error {
	newPos, err := forth.Pop[int](f.Stack)
	if err != nil {
		return forth.EmptyStack
	}
	if newPos < 0 {
		newPos = 0
	}
	if newPos > i.slicesLen {
		newPos = i.slicesLen
	}
	i.position = newPos
	f.Stack.Push(i.position)
	return nil
}

func (i *Interpreter) Position(f *forth.Forth) error {
	f.Stack.Push(i.position)
	return nil
}

func (i *Interpreter) Length(f *forth.Forth) error {
	f.Stack.Push(i.slicesLen)
	return nil
}

const maxInt16Value = float64((1 << 15) + 1)

func (i *Interpreter) saveWavFile() error {
	log.Printf("Saving file: %v", i.outFile)
	w := wav.CreateDefaultWav()
	buffer := new(bytes.Buffer)

	for idx, sampleLeft := range i.outputLeft {
		l := uint16(sampleLeft * maxInt16Value)
		bytesLeft := []byte{0, 0}
		binary.LittleEndian.PutUint16(bytesLeft, l)
		_, _ = buffer.Write(bytesLeft)

		sampleRight := i.outputRight[idx]
		r := uint16(sampleRight * maxInt16Value)
		bytesRight := []byte{0, 0}
		binary.LittleEndian.PutUint16(bytesRight, r)
		_, _ = buffer.Write(bytesRight)
	}
	w.Data = &wav.DataBytes{
		Samples: buffer.Bytes(),
	}

	f, err := os.Create(i.outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return w.Write(f)
}
