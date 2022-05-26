package parser

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

// constant: 1, 2, 3
// variable: $a, $b, $c
// function: seq [ 1 2 3 ]
// atom in braces: ( rand [ 1 2 3 ] )
func parseAtom(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	token := line.Fields[0]
	if n, err := strconv.Atoi(token); err == nil {
		return common.Const(int64(n)), 1, nil
	}
	if n, err := strconv.ParseFloat(token, 64); err == nil {
		return common.FloatConst(n), 1, nil
	}
	if n, ok := scale.Parse(token); ok {
		return common.Const(int64(n.Num)), 1, nil
	}
	if strings.HasPrefix(token, "$") {
		return common.Var(token[1:]), 1, nil
	}
	if token == "[" {
		fn, shift, err := parseList(scale, line)
		if err != nil {
			return nil, 0, err
		}
		// Include '[' and ']'
		return common.Lst(fn...), shift + 2, nil
	}
	if parser, ok := valueFnParser[token]; ok {
		fn, shift, err := parser(scale, line)
		if err != nil {
			return nil, 0, err
		}
		return fn, shift, nil
	}
	return nil, 0, fmt.Errorf("Don't know how to parse: %v", line)
}

func parseList(scale notes.Scale, line *LineCtx) ([]types.ValueFn, int, error) {
	atoms := []types.ValueFn{}
	l := len(line.Fields)
	i := 1
	for i < l {
		token := line.Fields[i]
		if token == "]" {
			return atoms, i + 1, nil
		}
		fn, shift, err := parseAtom(scale, line.Shift(i))
		if err != nil {
			return nil, 0, err
		}
		atoms = append(atoms, fn)
		i += shift
	}
	return nil, 0, fmt.Errorf("Closing ']' not found.")
}
