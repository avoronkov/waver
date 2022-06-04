package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ModParser = func(scale notes.Scale, line *LineCtx) (types.Modifier, int, error)

// : 4
func parseEvery(scale notes.Scale, line *LineCtx) (types.Modifier, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Every (':'): %v", line)
	}

	fn, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.Every(fn), shift + 1, nil
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

func parseBefore(scale notes.Scale, line *LineCtx) (types.Modifier, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for Before ('<'): %v", line)
	}

	fn, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.Before(fn), shift + 1, nil
}

func parseAfter(scale notes.Scale, line *LineCtx) (types.Modifier, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for After ('>')")
	}

	fn, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.After(fn), shift + 1, nil
}
