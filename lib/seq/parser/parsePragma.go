package parser

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	yaml "gopkg.in/yaml.v3"
)

type pragmaParser struct {
	Usage      string
	Desc       string
	Parse      func(*Parser, []lexer.Token, []map[string]any) error
	Deprecated bool
}

func (p *Parser) parsePragma(lx *lexer.Lexer) error {
	percent, err := lx.Pop()
	if err != nil {
		return err
	}
	switch percent.(type) {
	case lexer.PercentToken, lexer.DoublePercentToken:
	default:
		return fmt.Errorf("Unexpected token at the begining of pragma: %v (%T)", percent, percent)
	}

	pr, err := lx.Pop()
	if err != nil {
		return err
	}
	pragma, ok := pr.(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Expected pragma identifier, found: %v", pr)
	}
	ps := string(pragma)
	parserFn, ok := p.PragmaParsers[ps]
	if !ok {
		return fmt.Errorf("Unknown pragma: %v", ps)
	}

	fields := []lexer.Token{}
	body := ""
L:
	for {
		t, err := lx.Pop()
		if err != nil {
			return err
		}
		switch a := t.(type) {
		case lexer.EolToken, lexer.EofToken:
			break L
		case lexer.BodyToken:
			body = string(a)
			// read closing %% and <EOL>
			p, err := lx.Pop()
			if err != nil {
				return err
			}
			if _, ok := p.(lexer.DoublePercentToken); !ok {
				return fmt.Errorf("Expected closing %%, found: %v (%T)", p, p)
			}
			end, err := lx.Pop()
			if err != nil {
				return err
			}
			switch end.(type) {
			case lexer.EolToken, lexer.EofToken:
			default:
				return fmt.Errorf("Expected EOL, found: %v (%T)", p, p)
			}
			break L
		}
		fields = append(fields, t)
	}

	options := []map[string]any{}
	if err := p.parsePragmaOptions(body, &options); err != nil {
		return err
	}

	return parserFn.Parse(p, fields, options)
}

// % sample kick "7/kick"
func parseSample(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 2 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	smp := string(fields[0].(lexer.IdentToken))
	filename := fields[1].(lexer.StringLiteral)
	in, err := config.ParseSample(
		string(filename),
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(smp, in)
	return nil
}

// %wave string "sine"
func parseWave(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 2 {
		return fmt.Errorf("Incorrect number of arguments for 'inst' pragma: %v", fields)
	}
	inst := string(fields[0].(lexer.IdentToken))
	waveName := fields[1].(lexer.StringLiteral)
	in, err := config.ParseInstrument(
		string(waveName),
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(inst, in)
	return nil
}

// %form name "path"
func parseForm(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 2 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	name := string(fields[0].(lexer.IdentToken))
	filename := fields[1].(lexer.StringLiteral)
	in, err := config.ParseForm(
		string(filename),
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(name, in)
	return nil
}

// %lagrange name "path"
func parseLagrange(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 2 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	name := string(fields[0].(lexer.IdentToken))
	filename := fields[1].(lexer.StringLiteral)
	in, err := config.ParseLagrange(
		string(filename),
		append(options, p.globalFilters...),
		config.Param("tempo", p.tempo),
	)
	if err != nil {
		return err
	}
	p.instSet.AddInstrument(name, in)
	return nil
}

// %tempo 130
func parseTempo(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 1 {
		return fmt.Errorf("Incorrect number of arguments for 'tempo' pragma: %v", fields)
	}
	n, ok := fields[0].(lexer.NumberToken)
	if !ok {
		return fmt.Errorf("Cannot parse tempo, unexpected token: %v (%T)", fields[0], fields[0])
	}
	p.tempo = int(n)
	for _, ts := range p.tempoSetters {
		ts.SetTempo(p.tempo)
	}
	return nil
}

// %%filter
func parseFilter(p *Parser, fields []lexer.Token, options []map[string]any) error {
	p.globalFilters = options
	return nil
}

// % stop <int>
func parseStopPragma(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 1 {
		return fmt.Errorf("Incorrect number of arguments for 'stop' pragma: %v", fields)
	}
	n, ok := fields[0].(lexer.NumberToken)
	if !ok {
		return fmt.Errorf("Cannot parse 'stop', unexpected token: %v (%T)", fields[0], fields[0])
	}
	p.seq.SetStopBit(int64(n))
	return nil
}

// %srand 14254
func parseSrandPragma(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 1 {
		return fmt.Errorf("Incorrect number of arguments for 'srand' pragma: %v", fields)
	}
	n, ok := fields[0].(lexer.NumberToken)
	if !ok {
		return fmt.Errorf("Cannot parse 'srand', unexpected token: %v (%T)", fields[0], fields[0])
	}
	common.Srand(int64(n))
	return nil
}

// % scale edo12|edo19
func parseScalePragma(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 1 {
		return fmt.Errorf("Incorrect number of arguments for 'scale' pragma: %v", fields)
	}
	n, ok := fields[0].(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Cannot parse 'scale', unexpected token: %v (%T)", fields[0], fields[0])
	}

	var scale notes.Scale
	switch n {
	case "edo12":
		log.Printf("Using Standard 12 tone scale.")
		scale = notes.NewStandard()
	case "edo19":
		log.Printf("Using EDO-19 scale.")
		scale = notes.NewEdo19()
	default:
		return fmt.Errorf("Unknown scale: %v", n)
	}
	p.scale = scale
	for _, setScale := range p.scaleSetters {
		setScale(scale)
	}
	if stdFuncs, ok := scale.(notes.StdFuncsScale); ok {
		reader := bytes.NewReader(stdFuncs.Std())
		if err := p.parseReader(reader); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parsePragmaOptions(body string, opts *[]map[string]any) error {
	if strings.TrimSpace(body) == "" {
		return nil
	}
	r := strings.NewReader(body)
	err := yaml.NewDecoder(r).Decode(opts)
	if err != nil {
		return err
	}
	return nil
}
