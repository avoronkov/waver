package pragma

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/static"
)

type PragmaParser struct {
	file string

	tempoSetters []TempoSetter
	instSet      InstrumentSet
}

func New(file string, opts ...func(*PragmaParser)) *PragmaParser {
	pp := &PragmaParser{
		file: file,
	}

	for _, opt := range opts {
		opt(pp)
	}

	return pp
}

func (p *PragmaParser) Parse() error {
	f, err := os.Open(p.file)
	if err != nil {
		return err
	}
	defer f.Close()
	return p.parseReader(f)
}

func (p *PragmaParser) parseReader(reader io.Reader) error {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		text := sc.Text()
		if strings.HasPrefix(text, "% ") {
			fields := strings.Fields(text)
			if len(fields) < 2 {
				return fmt.Errorf("Not enough arguments for pragma ('%%'): %v", text)
			}
			log.Printf("[PRAGMA] parsing: %v", text)
			switch pragma := fields[1]; pragma {
			case "tempo":
				if err := p.parseTempo(fields); err != nil {
					return err
				}
			case "sample":
				if err := p.parseSample(fields); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown pragma ('%%'): %v", pragma)
			}
		} else {
			continue
		}
	}
	return sc.Err()
}

func (p *PragmaParser) parseTempo(fields []string) error {
	if len(fields) != 3 {
		return fmt.Errorf("Incorrect number of arguments for 'tempo' pragma: %v", fields)
	}
	n, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("Cannot parse tempo: %v", err)
	}
	for _, ts := range p.tempoSetters {
		ts.SetTempo(n)
	}
	return nil
}

// % sample 2k "2-2-kick.wav"
func (p *PragmaParser) parseSample(fields []string) error {
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	smp := fields[2]
	filename := filepath.Join("samples", strings.Trim(fields[3], "\""))

	return p.handleSample(smp, filename)
}

func (p *PragmaParser) handleSample(name, file string) error {
	log.Printf("Using sample '%v' from '%v'", name, file)
	data, err := static.Files.ReadFile(file)
	if err != nil {
		return err
	}
	w, err := waves.ParseSample(data)
	if err != nil {
		return err
	}
	in := instruments.NewInstrument(w)
	p.instSet.AddSampledInstrument(name, in)
	return nil
}
