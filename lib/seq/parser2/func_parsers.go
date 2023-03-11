package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

type FunctionParser func(p *Parser, lx *lexer.Lexer, fn string) (types.ValueFn, error)

// seq [ 1 2 3 ]
func parseSequence(p *Parser, lx *lexer.Lexer, fn string) (types.ValueFn, error) {
	values, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("seq:%v", lx.LineNum())
	return common.Sequence(key, values), nil
}

// up 3 A4
// down 4 C3
func parseUpDown(p *Parser, lx *lexer.Lexer, fn string) (types.ValueFn, error) {
	arg, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	value, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	invert := fn == "down"
	return common.Up(arg, value, invert), nil
}

// repeat 4 seq [ A4 C5 ]
func parseRepeat(p *Parser, lx *lexer.Lexer, fn string) (types.ValueFn, error) {
	times, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	arg, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	keyIndex := fmt.Sprintf("repeat-idx:%v", lx.LineNum())
	keyValue := fmt.Sprintf("repeat-value:%v", lx.LineNum())
	return common.Repeat(keyIndex, keyValue, times, arg), nil
}

// loop 16 rand [ ... ]
func parseLoop(p *Parser, lx *lexer.Lexer, fn string) (types.ValueFn, error) {
	size, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	arg, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}
	return common.Loop(size, arg), nil
}

type singleArgFuncion func(types.ValueFn) types.ValueFn

func makeSingleArgValueFnParser(name string, fn singleArgFuncion) FunctionParser {
	return func(p *Parser, lx *lexer.Lexer, name string) (types.ValueFn, error) {
		value, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		return fn(value), nil
	}
}
