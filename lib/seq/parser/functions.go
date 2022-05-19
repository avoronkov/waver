package parser

import (
	"fmt"
	"log"
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
func parseAtom(scale notes.Scale, fields []string) (types.ValueFn, int, error) {
	token := fields[0]
	if n, err := strconv.Atoi(token); err == nil {
		return common.Const(int64(n)), 1, nil
	}
	if n, ok := scale.Parse(token); ok {
		log.Printf("Scale.Parse(%v) = %v", token, n)
		return common.Const(int64(n.Num)), 1, nil
	}
	if strings.HasPrefix(token, "$") {
		return common.Var(token[1:]), 1, nil
	}
	if token == "[" {
		fn, shift, err := parseList(scale, fields)
		if err != nil {
			return nil, 0, err
		}
		// Include '[' and ']'
		return common.Lst(fn...), shift + 2, nil
	}
	if parser, ok := valueFnParser[token]; ok {
		fn, shift, err := parser(scale, fields)
		if err != nil {
			return nil, 0, err
		}
		return fn, shift, nil
	}
	return nil, 0, fmt.Errorf("Don't know how to parse: %v", fields)
}

func parseList(scale notes.Scale, fields []string) ([]types.ValueFn, int, error) {
	atoms := []types.ValueFn{}
	l := len(fields)
	i := 1
	for i < l {
		token := fields[i]
		if token == "]" {
			return atoms, i + 1, nil
		}
		fn, shift, err := parseAtom(scale, fields[i:])
		if err != nil {
			return nil, 0, err
		}
		atoms = append(atoms, fn)
		i += shift
	}
	return nil, 0, fmt.Errorf("Closing ']' not found.")
}
