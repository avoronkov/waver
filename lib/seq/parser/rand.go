package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

// rand [ 1 2 3 ]
func parseRandom(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", line)
	}

	values, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	return common.Random(values), shift + 1, nil
}