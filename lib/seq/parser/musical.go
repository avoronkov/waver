package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func makeMusParser(name string, shifts ...int64) ValueFnParser {
	return func(fields []string) (types.ValueFn, int, error) {
		if len(fields) < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, fields)
		}
		x, shift, err := parseAtom(fields[1:])
		if err != nil {
			return nil, 0, err
		}
		return common.ChordFn(x, shifts...), shift + 1, nil
	}
}