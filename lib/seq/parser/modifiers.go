package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(p *Parser, line *LineCtx) (types.Modifier, int, error)

type singleArgModifier = func(types.ValueFn) types.Modifier

func makeSingleArgModParser(name string, fn singleArgModifier) ModParser {
	return func(p *Parser, line *LineCtx) (types.Modifier, int, error) {
		if line.Len() < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for %v: %v", name, line)
		}

		arg, shift, err := p.parseAtom(line.Shift(1))
		if err != nil {
			return nil, 0, err
		}

		return fn(arg), shift + 1, nil
	}
}

type twoArgsModifier = func(types.ValueFn, types.ValueFn) types.Modifier

func makeTwoArgsModParser(name string, fn twoArgsModifier) ModParser {
	return func(p *Parser, line *LineCtx) (types.Modifier, int, error) {
		if line.Len() < 3 {
			return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
		}

		arg1, shift, err := p.parseAtom(line.Shift(1))
		if err != nil {
			return nil, 0, err
		}
		arg2, shift2, err := p.parseAtom(line.Shift(1 + shift))
		if err != nil {
			return nil, 0, err
		}

		return fn(arg1, arg2), shift + shift2 + 1, nil
	}
}

// "+ 2", "- 2"
func parseShift(p *Parser, line *LineCtx) (types.Modifier, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Shift ('+' / '-'): %v", line)
	}

	fn, shift, err := p.parseAtom(line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.Shift(fn), shift + 1, nil
}
