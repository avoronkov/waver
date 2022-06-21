package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

func parseUserFunction(scale notes.Scale, line *LineCtx, name, argName string, body types.ValueFn) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
	}

	arg, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.UserFunction(argName, arg, body), shift + 1, nil
}
