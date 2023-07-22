package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

// 0..15 : 4 -> { sig }
func (p *Parser) parseSignalStatement(lx *lexer.Lexer) error {
	startingBit, endingBit, err := p.parseRange(lx)
	if err != nil {
		return err
	}
	mods, err := p.parseModifiers(lx)
	if err != nil {
		return err
	}
	if startingBit > 0 {
		mods = append(mods, common.After(common.Const(startingBit)))
	}
	if endingBit >= 0 {
		mods = append(mods, common.Before(common.Const(endingBit)))
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

func (p *Parser) parseSignaler(lx *lexer.Lexer) (sigs []types.Signaler, err error) {
	mods, err := p.parseModifiers(lx)
	if err != nil {
		return nil, err
	}
	sigs, err = p.parseSignals(lx)
	if err != nil {
		return nil, err
	}
	for i, sig := range sigs {
		sigs[i] = common.Chain(sig, mods...)
	}
	return sigs, nil
}

func (p *Parser) parseRange(lx *lexer.Lexer) (begin int64, end int64, err error) {
	end = -1
	token, e := lx.Top()
	if e != nil {
		err = e
		return
	}
	switch t := token.(type) {
	case lexer.HexToken:
		begin = int64(t)
	case lexer.NumberToken:
		begin = int64(t)
	default:
		return
	}
	// Drop begining of the range
	lx.Drop()

	token, e = lx.Top()
	if e != nil {
		err = e
		return
	}
	if _, ok := token.(lexer.DoubleDot); !ok {
		return
	}

	// Drop ".."
	lx.Drop()

	token, e = lx.Top()
	if e != nil {
		err = e
		return
	}
	switch t := token.(type) {
	case lexer.HexToken:
		end = int64(t)
	case lexer.NumberToken:
		end = int64(t)
	default:
		return
	}

	// Drop end of the range
	lx.Drop()
	return
}

func (p *Parser) parseModifiers(lx *lexer.Lexer) (result []types.Modifier, err error) {
L:
	for {
		token, err := lx.Top()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.ArrowRightToken); ok {
			_, _ = lx.Pop()
			break L
		}

		parser, ok := p.ModParsers[token]
		if !ok {
			break
		}
		_, _ = lx.Pop()

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
		switch t := token.(type) {
		case lexer.LCurlyBracket:
			sig, err := p.parseSignal(lx)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, sig)
		case lexer.VerticalBar:
			lx.Drop()
		case lexer.EolToken, lexer.EofToken:
			return sigs, nil
		case lexer.IdentToken:
			if us, ok := p.userSignalers[t.String()]; ok {
				lx.Drop()
				sigs = append(sigs, us...)
			} else {
				// the same as default
				sig, err := p.parsePlainSignal(lx)
				if err != nil {
					return nil, err
				}
				sigs = append(sigs, sig)
			}
		default:
			sig, err := p.parsePlainSignal(lx)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, sig)
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

func (p *Parser) parsePlainSignal(lx *lexer.Lexer) (types.Signaler, error) {
	// Parse instrument
	in, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse note
	tok, err := lx.Top()
	if err != nil {
		return nil, err
	}
	switch tok.(type) {
	case lexer.EolToken, lexer.EofToken:
		return common.Note(p.scale, in, common.StrConst("_")), nil
	case lexer.VerticalBar:
		_, _ = lx.Pop()
		return common.Note(p.scale, in, common.StrConst("_")), nil
	}

	note, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse amplitude
	tok, err = lx.Top()
	if err != nil {
		return nil, err
	}
	switch tok.(type) {
	case lexer.EolToken, lexer.EofToken:
		return common.Note(p.scale, in, note), nil
	case lexer.VerticalBar:
		_, _ = lx.Pop()
		return common.Note(p.scale, in, note), nil
	}

	amp, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Parse duration
	tok, err = lx.Top()
	if err != nil {
		return nil, err
	}
	switch tok.(type) {
	case lexer.EolToken, lexer.EofToken:
		return common.Note(p.scale, in, note, common.NoteAmp(amp)), nil
	case lexer.VerticalBar:
		_, _ = lx.Pop()
		return common.Note(p.scale, in, note, common.NoteAmp(amp)), nil
	}

	dur, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	// Check closing '|'
	tok, err = lx.Top()
	if err != nil {
		return nil, err
	}
	switch tok.(type) {
	case lexer.EolToken, lexer.EofToken:
		return common.Note(p.scale, in, note, common.NoteAmp(amp), common.NoteDur(dur)), nil
	case lexer.VerticalBar:
		_, _ = lx.Pop()
		return common.Note(p.scale, in, note, common.NoteAmp(amp), common.NoteDur(dur)), nil
	}
	return nil, fmt.Errorf("Unexpected token at the end of signal: %v", tok)
}
