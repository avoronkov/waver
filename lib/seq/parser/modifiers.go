package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(scale notes.Scale, line *LineCtx) (types.Modifier, int, error)

type singleArgModifier = func(types.ValueFn) types.Modifier

func makeSingleArgModParser(name string, fn singleArgModifier) ModParser {
	return func(scale notes.Scale, line *LineCtx) (types.Modifier, int, error) {
		if line.Len() < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for %v: %v", name, line)
		}

		arg, shift, err := parseAtom(scale, line.Shift(1))
		if err != nil {
			return nil, 0, err
		}

		return fn(arg), shift + 1, nil
	}
}

// "+ 2", "- 2"
func parseShift(scale notes.Scale, line *LineCtx) (types.Modifier, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Shift ('+' / '-'): %v", line)
	}

	fn, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.Shift(fn), shift + 1, nil
}
