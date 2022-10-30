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
	var parserFn func(*Parser, []string, []map[string]any) error
	switch pragma := fields[1]; pragma {
	case "tempo":
		parserFn = parseTempo
	case "sample":
		parserFn = parseSample
	case "inst", "wave":
		parserFn = parseWave
	case "filter":
		parserFn = parseFilter
	default:
		return fmt.Errorf("Unknown pragma ('%%'): %v", pragma)
	}
	options := []map[string]any{}
	if err := p.parsePragmaOptions(body, &options); err != nil {
		return err
	}
	return parserFn(p, fields, options)
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

func (p *Parser) parsePragmaOptions(body string, opts *[]map[string]any) error {
	if body == "" {
		return nil
	}
	r := strings.NewReader(body)
	err := yaml.NewDecoder(r).Decode(opts)
	if err != nil {
		return err
	}
	return nil
}

// % sample kick "7/kick"
func parseSample(p *Parser, fields []string, options []map[string]any) error {
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	smp := fields[2]
	filename := strings.Trim(fields[3], "\"'")
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

func parseWave(p *Parser, fields []string, options []map[string]any) error {
	if len(fields) != 4 {
		return fmt.Errorf("Incorrect number of arguments for 'inst' pragma: %v", fields)
	}
	inst := fields[2]
	waveName := strings.Trim(fields[3], "\"'")
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

func parseTempo(p *Parser, fields []string, options []map[string]any) error {
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

func parseFilter(p *Parser, fields []string, options []map[string]any) error {
	p.globalFilters = options
	return nil
}
