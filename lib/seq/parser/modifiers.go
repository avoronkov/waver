package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(scale notes.Scale, fields []string) (types.Modifier, int, error)

// : 4
func parseEvery(scale notes.Scale, fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Every (':')")
	}

	fn, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.Every(fn), shift + 1, nil
}

// "+ 2", "- 2"
func parseShift(scale notes.Scale, fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Shift ('+' / '-')")
	}

	fn, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.Shift(fn), shift + 1, nil
}

func parseBefore(scale notes.Scale, fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Before ('<')")
	}

	fn, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.Before(fn), shift + 1, nil
}

func parseAfter(scale notes.Scale, fields []string) (types.Modifier, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for After ('>')")
	}

	fn, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}

	return common.After(fn), shift + 1, nil
}
