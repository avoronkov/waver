package parser2

import (
	"fmt"
	"strings"

	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/seq/lexer"
	yaml "gopkg.in/yaml.v3"
)

type pragmaParser func(*Parser, []lexer.Token, []map[string]any) error

func (p *Parser) parsePragma(lx *lexer.Lexer) error {
	percent, err := lx.Pop()
	if err != nil {
		return err
	}
	switch percent.(type) {
	case lexer.Percent, lexer.DoublePercent:
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
	parserFn, ok := p.pragmaParsers[pragma.Value]
	if !ok {
		return fmt.Errorf("Unknown pragma: %v", pragma.Value)
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
			if _, ok := p.(lexer.DoublePercent); !ok {
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

	return parserFn(p, fields, options)
}

// % sample kick "7/kick"
func parseSample(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 2 {
		return fmt.Errorf("Incorrect number of arguments for 'sample' pragma: %v", fields)
	}
	smp := fields[0].(lexer.IdentToken).Value
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
	inst := fields[0].(lexer.IdentToken).Value
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

// %tempo 130
func parseTempo(p *Parser, fields []lexer.Token, options []map[string]any) error {
	if len(fields) != 1 {
		return fmt.Errorf("Incorrect number of arguments for 'tempo' pragma: %v", fields)
	}
	n, ok := fields[0].(lexer.NumberToken)
	if !ok {
		return fmt.Errorf("Cannot parse tempo, unexpected token: %v (%T)", fields[0], fields[0])
	}
	p.tempo = int(n.Num)
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
