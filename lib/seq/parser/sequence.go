package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type ValueFnParser = func(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error)

// seq [ 1 2 3 ]
func parseSequence(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", line)
	}

	values, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	return common.Sequence(values), shift + 1, nil
}
