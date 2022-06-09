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

	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/static"

	yaml "gopkg.in/yaml.v3"
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
		singleLine := strings.HasPrefix(text, "% ")
		multiline := strings.HasPrefix(text, "%% ")
		if singleLine || multiline {
			fields := strings.Fields(text)
			if len(fields) < 2 {
				return fmt.Errorf("Not enough arguments for pragma ('%%'): %v", text)
			}
			var body string
			if multiline {
				b, err := p.parseMultilinePragma(sc)
				if err != nil {
					return err
				}
				body = b
			}
			log.Printf("[PRAGMA] parsing: %v", text)
			switch pragma := fields[1]; pragma {
			case "tempo":
				if err := p.parseTempo(fields); err != nil {
					return err
				}
			case "sample":
				if err := p.parseSample(fields, body); err != nil {
					return err
				}
			case "inst":
				if err := p.parseInstrument(fields, body); err != nil {
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

func (p *PragmaParser) parseMultilinePragma(sc *bufio.Scanner) (string, error) {
	var lines []string
	for sc.Scan() {
		line := sc.Text()
		if line == "%%" {
			return strings.Join(lines, "\n"), nil
		}
		lines = append(lines, line)
	}
	err := sc.Err()
	if err == nil {
		err = fmt.Errorf("Unexpected end-of-file")
	}
	return "", err
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
func (p *PragmaParser) parseSample(fields []string, body string) error {
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

// % inst 1 'sine'
func (p *PragmaParser) parseInstrument(fields []string, body string) error {
	log.Printf("parseInstrument: %v | %v", fields, body)
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'inst' pragma: %v", fields)
	}
	instIdx, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("Instrument index is not an integer: %v", fields[2])
	}
	waveName := strings.Trim(fields[3], "'")
	var options []map[string]any
	if body != "" {
		options, err = p.parsePragmaOptions(body)
		if err != nil {
			return err
		}
	}
	in, err := config.ParseInstrument(waveName, options)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(instIdx, in)
	return nil
}

func (p *PragmaParser) parsePragmaOptions(body string) ([]map[string]any, error) {
	options := []map[string]any{}
	r := strings.NewReader(body)
	err := yaml.NewDecoder(r).Decode(&options)
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (p *PragmaParser) handleWave(wave string) (waves.Wave, error) {
	if w, ok := waves.Waves[wave]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("Unknown wave: %v", wave)
}
