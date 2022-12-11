package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

// : 4 -> { sig }
func (p *Parser) parseSignalStatement(lx *lexer.Lexer) error {
	mods, err := p.parseModifiers(lx)
	if err != nil {
		return err
	}
	sigs, err := p.parseSignals(lx)
	if err != nil {
		return err
	}
	for _, sig := range sigs {
		x := common.Chain(sig, mods...)
		p.seq.Add(x)
	}
	return nil
}

func (p *Parser) parseModifiers(lx *lexer.Lexer) (result []types.Modifier, err error) {
L:
	for {
		token, err := lx.Pop()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.ArrowRightToken); ok {
			break L
		}

		parser, ok := p.modParsers[token]
		if !ok {
			return nil, fmt.Errorf("Unknown modParser: %v", token)
		}

		modifier, err := parser(p, lx)
		if err != nil {
			return nil, err
		}
		result = append(result, modifier)
	}
	return
}

func (p *Parser) parseSignals(lx *lexer.Lexer) (sigs []types.Signaler, err error) {
	for {
		token, err := lx.Top()
		if err != nil {
			return nil, err
		}
		switch token.(type) {
		case lexer.LCurlyBracket:
			sig, err := p.parseSignal(lx)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, sig)
		case lexer.EolToken, lexer.EofToken:
			return sigs, nil
		default:
			return nil, fmt.Errorf("Unexpected token while parsing signals: %v", token)
		}
	}
}

func (p *Parser) parseSignal(lx *lexer.Lexer) (types.Signaler, error) {
	// Skip '{'
	_, _ = lx.Pop()
	// Parse instrument
	in, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse note
	end, err := lx.Top()
	if err != nil {
		return nil, err
	}
	if _, ok := end.(lexer.RCurlyBracket); ok {
		_, _ = lx.Pop()
		return common.Note(p.scale, in, common.StrConst("_")), nil
	}

	note, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse amplitude
	end, err = lx.Top()
	if err != nil {
		return nil, err
	}
	if _, ok := end.(lexer.RCurlyBracket); ok {
		_, _ = lx.Pop()
		return common.Note(p.scale, in, note), nil
	}

	amp, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse duration
	end, err = lx.Top()
	if err != nil {
		return nil, err
	}
	if _, ok := end.(lexer.RCurlyBracket); ok {
		_, _ = lx.Pop()
		return common.Note(p.scale, in, note, common.NoteAmp(amp)), nil
	}

	dur, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Check closing '}'
	end, err = lx.Pop()
	if err != nil {
		return nil, err
	}
	if _, ok := end.(lexer.RCurlyBracket); !ok {
		return nil, fmt.Errorf("Expected '}', found: %v", end)
	}
	return common.Note(p.scale, in, note, common.NoteAmp(amp), common.NoteDur(dur)), nil
}
