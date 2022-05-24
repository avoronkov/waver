package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func parseUpDown(scale notes.Scale, fields []string) (types.ValueFn, int, error) {
	if len(fields) < 3 {
		return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", fields[0], fields)
	}
	arg, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}
	value, shift2, err := parseAtom(scale, fields[shift+1:])
	if err != nil {
		return nil, 0, err
	}
	invert := fields[0] == "down"
	return common.Up(arg, value, invert), shift + shift2 + 1, nil
}
