package parser

import (
	"fmt"
	"strconv"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
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
	if token == "[" {
		fn, shift, err := parseList(scale, line)
		if err != nil {
			return nil, 0, err
		}
		return common.Lst(fn...), shift, nil
	}
	if userFunc, ok := line.UserFunctions[token]; ok {
		fn, shift, err := parseUserFunction(scale, line, userFunc.name, userFunc.arg, userFunc.fn)
		if err != nil {
			return nil, 0, err
		}
		return fn, shift, err
	}
	if parser, ok := valueFnParser[token]; ok {
		fn, shift, err := parser(scale, line)
		if err != nil {
			return nil, 0, err
		}
		return fn, shift, nil
	}
	return common.Var(token), 1, nil
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

type singleArgFuncion = func(types.ValueFn) types.ValueFn

func MakeSingleArgValueFnParser(name string, fn singleArgFuncion) ValueFnParser {
	return func(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
		if line.Len() < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
		}
		values, shift, err := parseAtom(scale, line.Shift(1))
		if err != nil {
			return nil, 0, err
		}
		return fn(values), shift + 1, nil
	}
}
