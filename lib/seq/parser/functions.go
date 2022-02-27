package parser

import (
	"fmt"
	"strconv"
	"strings"

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

// constant: 1, 2, 3
// variable: $a, $b, $c
// function: seq [ 1 2 3 ]
// atom in braces: ( rand [ 1 2 3 ] )
func parseAtom(fields []string) (types.ValueFn, int, error) {
	token := fields[0]
	if n, err := strconv.Atoi(token); err == nil {
		return common.Const(int64(n)), 1, nil
	}
	if n, err := common.ParseStandardNote(token); err == nil {
		return common.Const(n.Number), 1, nil
	}
	if strings.HasPrefix(token, "$") {
		return common.Var(token[1:]), 1, nil
	}
	if token == "[" {
		fn, shift, err := parseList(fields[1:])
		if err != nil {
			return nil, 0, err
		}
		// Include '[' and ']'
		return fn, shift + 2, nil
	}
	return nil, 0, fmt.Errorf("Don't know how to parse: %v", fields)
}

func parseList(fields []string) (types.ValueFn, int, error) {
	atoms := []types.ValueFn{}
	l := len(fields)
	i := 0
	for i < l {
		token := fields[i]
		if token == "]" {
			return common.Lst(atoms...), i, nil
		}
		fn, shift, err := parseAtom(fields[i:])
		if err != nil {
			return nil, 0, err
		}
		atoms = append(atoms, fn)
		i += shift
	}
	return nil, 0, fmt.Errorf("Closing ']' not found.")
}
