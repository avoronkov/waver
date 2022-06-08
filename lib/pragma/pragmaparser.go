package pragma

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PragmaParser struct {
	file string

	tempoSetters []TempoSetter
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
			switch pragma := fields[1]; pragma {
			case "tempo":
				return p.parseTempo(fields)
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
