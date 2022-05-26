package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func parseUpDown(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 3 {
		return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", line.Fields[0], line)
	}
	arg, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	value, shift2, err := parseAtom(scale, line.Shift(shift+1))
	if err != nil {
		return nil, 0, err
	}
	invert := line.Fields[0] == "down"
	return common.Up(arg, value, invert), shift + shift2 + 1, nil
}
