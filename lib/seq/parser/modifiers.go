package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(fields []string) (types.Modifier, int, error)

// : 4
func parseEvery(fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Every (':')")
	}

	fn, shift, err := parseAtom(fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.Every(fn), shift + 1, nil
}

// "+ 2", "- 2"
func parseShift(fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Shift ('+' / '-')")
	}

	fn, shift, err := parseAtom(fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.Shift(fn), shift + 1, nil
}
