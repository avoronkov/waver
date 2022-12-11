package parser2

import (
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ModParser func(p *Parser, lx *lexer.Lexer) (types.Modifier, error)

type singleArgModifier func(types.ValueFn) types.Modifier

func makeSingleArgModParser(name string, fn singleArgModifier) ModParser {
	return func(p *Parser, lx *lexer.Lexer) (types.Modifier, error) {
		arg, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		return fn(arg), nil
	}
}

type twoArgsModifier func(types.ValueFn, types.ValueFn) types.Modifier

func makeTwoArgsModParser(name string, fn twoArgsModifier) ModParser {
	return func(p *Parser, lx *lexer.Lexer) (types.Modifier, error) {
		a, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		b, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		return fn(a, b), nil
	}
}
