package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/midisynth/config"
	"github.com/avoronkov/waver/lib/midisynth/instruments"
	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/seq/lexer"
)

func (p *Parser) parseAssignmentStatement(lx *lexer.Lexer) error {
	first, err := lx.Pop()
	if err != nil {
		return err
	}
	name, ok := first.(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Unexpected first token in assignent statement: %v (%T)", first, first)
	}

	second, err := lx.Pop()
	if err != nil {
		return err
	}

	if _, ok := second.(lexer.AssignToken); ok {
		return p.parseAssignVar(lx, string(name))
	}

	if _, ok := second.(lexer.DoubleAssignToken); ok {
		return p.parseAssignSignaler(lx, string(name))
	}

	param, ok := second.(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Unexpected second token in assignment statement: %v (%T)", second, second)
	}

	third, err := lx.Pop()
	if err != nil {
		return err
	}
	if _, ok := third.(lexer.AssignToken); !ok {
		return fmt.Errorf("Unexpected third token in assignment statement: %v (%T)", third, third)
	}

	return p.parseUdf(lx, string(name), string(param))
}

func (p *Parser) parseAssignInstrument(lx *lexer.Lexer, name string) error {
	tok, err := lx.Pop()
	if err != nil {
		return err
	}

	options := []map[string]any{}

L:
	for {
		tok, err := lx.Pop()
		if err != nil {
			return err
		}
		switch tok.(type) {
		case lexer.EolToken, lexer.EofToken, lexer.CommentToken:
			break L
		case lexer.VerticalBar:
			opt, err := p.parseInstrumentAssignmentOption(lx)
			if err != nil {
				return err
			}
			options = append(options, opt)
		default:
			return fmt.Errorf("Unexpected token while parsing instrument filters: %v (%T)", tok, tok)
		}

	}

	var in *instruments.Instrument
	switch t := tok.(type) {
	case lexer.IdentToken:
		in, err = config.ParseInstrument(
			t.String(),
			append(options, p.globalFilters...),
			config.Param("tempo", p.tempo),
		)
	case lexer.StringLiteral:
		in, err = config.ParseSample(
			t.String(),
			append(options, p.globalFilters...),
			config.Param("tempo", p.tempo),
		)
	default:
		return fmt.Errorf("Unexpected token type: %v (%v)", tok, tok)
	}

	if err != nil {
		return err
	}
	p.instSet.AddInstrument(name, in)

	p.instrumentVariables.Add(name)
	return nil
}

func token2scalar(t lexer.Token) (any, error) {
	switch a := t.(type) {
	case lexer.NumberToken:
		return int64(a), nil
	case lexer.FloatToken:
		return float64(a), nil
	case lexer.IdentToken:
		if a.String() == "true" {
			return true, nil
		}
		if a.String() == "false" {
			return false, nil
		}
	case lexer.StringLiteral:
		return a.String(), nil
	}
	return nil, fmt.Errorf("Cannot convert token to scalar: %v (%T)", t, t)
}

// E.g. `am freq=10.0 int=0.25`, `exp=2.0`
func (p *Parser) parseInstrumentAssignmentOption(lx *lexer.Lexer) (map[string]any, error) {
	tok, err := lx.Top()
	if err != nil {
		return nil, err
	}
	code, ok := tok.(lexer.CodeLiteral)
	if ok {
		lx.Drop()
		return map[string]any{
			"code": code.String(),
		}, nil
	}
	identName, ok := tok.(lexer.IdentToken)
	if !ok {
		return nil, fmt.Errorf("Unexpected token: %v (%T)", tok, tok)
	}
	lx.Drop()
	name := identName.String()

	// Check if it's an assignment
	tok2, err := lx.Top()
	if err != nil {
		return nil, err
	}
	if _, ok := tok2.(lexer.AssignToken); ok {
		lx.Drop()
		// read next token
		tok3, err := lx.Pop()
		if err != nil {
			return nil, err
		}
		value, err := token2scalar(tok3)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			name: value,
		}, nil
	}

	opts := map[string]any{}
	for {
		// key = value
		keyTok, err := lx.Top()
		if err != nil {
			return nil, err
		}
		keyIdent, ok := keyTok.(lexer.IdentToken)
		if !ok {
			break
		}
		lx.Drop()
		assignTok, err := lx.Pop()
		if err != nil {
			return nil, err
		}
		if _, ok := assignTok.(lexer.AssignToken); !ok {
			return nil, fmt.Errorf("Unexpected token: %v (%v)", assignTok, assignTok)
		}
		valueTok, err := lx.Pop()
		if err != nil {
			return nil, err
		}
		value, err := token2scalar(valueTok)
		if err != nil {
			return nil, err
		}
		opts[keyIdent.String()] = value
	}

	if uf, ok := p.userFilters[name]; ok {
		opts["code"] = uf
		return map[string]any{
			"code": opts,
		}, nil
	}

	return map[string]any{
		name: opts,
	}, nil
}

func (p *Parser) checkInstrumentAssignment(token lexer.Token) bool {
	switch t := token.(type) {
	case lexer.StringLiteral:
		return true
	case lexer.IdentToken:
		if _, ok := waves.Waves[t.String()]; ok {
			return true
		}
		return p.instrumentVariables.Has(t.String())
	}
	return false
}

func (p *Parser) parseAssignVar(lx *lexer.Lexer, name string) error {
	// Check if it a definition of an instument.
	tok, err := lx.Top()
	if err != nil {
		return err
	}
	if p.checkInstrumentAssignment(tok) {
		return p.parseAssignInstrument(lx, name)
	}

	// Check if it is a definition of a filter
	if code, ok := tok.(lexer.CodeLiteral); ok {
		lx.Drop()
		return p.parseUserFilter(lx, name, code)
	}

	// parse atom
	atom, err := p.parseAtom(lx)
	if err != nil {
		return err
	}
	endl, err := lx.Pop()
	if err != nil {
		return err
	}
	switch endl.(type) {
	case lexer.EolToken, lexer.EofToken:
	default:
		return fmt.Errorf("Expected EOL at the end of assign statement, found: %v (%T)", endl, endl)
	}
	p.seq.Assign(name, atom)
	return nil
}

func (p *Parser) parseUserFilter(lx *lexer.Lexer, name string, code lexer.CodeLiteral) error {
	endl, err := lx.Pop()
	if err != nil {
		return err
	}
	switch endl.(type) {
	case lexer.EolToken, lexer.EofToken:
	default:
		return fmt.Errorf("Expected EOL at the end of assign statement, found: %v (%T)", endl, endl)
	}
	p.userFilters[name] = code.String()
	return nil
}

func (p *Parser) parseUdf(lx *lexer.Lexer, name string, param string) error {
	// parse atom
	atom, err := p.parseAtom(lx)
	if err != nil {
		return err
	}
	endl, err := lx.Pop()
	if err != nil {
		return err
	}
	switch endl.(type) {
	case lexer.EolToken, lexer.EofToken:
	default:
		return fmt.Errorf("Expected EOL at the end of assign statement, found: %v (%T)", endl, endl)
	}
	p.userFunctions[name] = UserFunction{
		Name: name,
		Arg:  param,
		Fn:   atom,
	}
	return nil
}

func (p *Parser) parseAssignSignaler(lx *lexer.Lexer, name string) error {
	sigs, err := p.parseSignaler(lx)
	if err != nil {
		return err
	}
	p.userSignalers[name] = sigs
	return nil
}
