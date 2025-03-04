package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ParseModifier func(p *Parser, lx *lexer.Lexer) (types.Modifier, error)

type ModParser struct {
	Usage      string
	Parse      func(p *Parser, lx *lexer.Lexer) (types.Modifier, error)
	Deprecated bool
}

type singleArgModifier func(types.ValueFn) types.Modifier

func makeSingleArgModParser(name string, fn singleArgModifier) ParseModifier {
	return func(p *Parser, lx *lexer.Lexer) (types.Modifier, error) {
		arg, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		return fn(arg), nil
	}
}

type twoArgsModifier func(types.ValueFn, types.ValueFn) types.Modifier

func makeTwoArgsModParser(name string, fn twoArgsModifier) ParseModifier {
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

func parseTimesModifier(p *Parser, lx *lexer.Lexer) (types.Modifier, error) {
	arg, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("times:%v", lx.LineNum())
	return common.Times(arg, key), nil
}
