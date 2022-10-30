package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/config"
	yaml "gopkg.in/yaml.v3"
)

func (p *Parser) parsePragma(text string, sc *bufio.Scanner) error {
	singleLine := strings.HasPrefix(text, "% ")
	multiLine := strings.HasPrefix(text, "%% ")
	if !singleLine && !multiLine {
		return nil
	}
	fields := strings.Fields(text)
	if len(fields) < 2 {
		return fmt.Errorf("Not enough arguments for pragma ('%%'): %v", text)
	}
	var body string
	if multiLine {
		b, err := p.parseMultilinePragma(sc)
		if err != nil {
			return err
		}
		body = b
	}
	switch pragma := fields[1]; pragma {
	case "tempo":
		if err := p.parseTempo(fields); err != nil {
			return err
		}
	case "sample":
		if err := p.parseSample(fields, body); err != nil {
			return err
		}
	case "inst", "wave":
		if err := p.parseInstrument(fields, body); err != nil {
			return err
		}
	case "filter":
		if err := p.parseFilter(fields, body); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown pragma ('%%'): %v", pragma)
	}
	return nil
}

func (p *Parser) parseMultilinePragma(sc *bufio.Scanner) (string, error) {
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

func (p *Parser) parseTempo(fields []string) error {
	if len(fields) != 3 {
		return fmt.Errorf("Incorrect number of arguments for 'tempo' pragma: %v", fields)
	}
	n, err := strconv.Atoi(fields[2])
	if err != nil {
		return fmt.Errorf("Cannot parse tempo: %v", err)
	}
	p.tempo = n
	for _, ts := range p.tempoSetters {
		ts.SetTempo(n)
	}
	return nil
}

// % sample 2k "2-2-kick.wav"
func (p *Parser) parseSample(fields []string, body string) error {
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	smp := fields[2]
	filename := strings.Trim(fields[3], "\"'")

	var options []map[string]any
	if body != "" {
		var err error
		options, err = p.parsePragmaOptions(body)
		if err != nil {
			return err
		}
	}
	in, err := config.ParseSample(
		filename,
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(smp, in)
	return nil
}

// % inst 1 'sine'
func (p *Parser) parseInstrument(fields []string, body string) (err error) {
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'inst' pragma: %v", fields)
	}
	inst := fields[2]
	waveName := strings.Trim(fields[3], "\"'")
	var options []map[string]any
	if body != "" {
		options, err = p.parsePragmaOptions(body)
		if err != nil {
			return err
		}
	}
	in, err := config.ParseInstrument(
		waveName,
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(inst, in)
	return nil
}
func (p *Parser) parsePragmaOptions(body string) ([]map[string]any, error) {
	options := []map[string]any{}
	r := strings.NewReader(body)
	err := yaml.NewDecoder(r).Decode(&options)
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (p *Parser) parseFilter(fields []string, body string) (err error) {
	var options []map[string]any
	if body != "" {
		options, err = p.parsePragmaOptions(body)
		if err != nil {
			return err
		}
	}
	p.globalFilters = options
	return nil
}
