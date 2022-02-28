package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type ValueFnParser = func(fields []string) (types.ValueFn, int, error)

// seq [ 1 2 3 ]
func parseSequence(fields []string) (types.ValueFn, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", fields)
	}

	values, shift, err := parseList(fields[1:])
	if err != nil {
		return nil, 0, err
	}
	return common.Sequence(values...), shift + 1, nil
}
