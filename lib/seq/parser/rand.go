package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

// rand [ 1 2 3 ]
func parseRandom(scale notes.Scale, fields []string) (types.ValueFn, int, error) {
	if len(fields) < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", fields)
	}

	values, shift, err := parseAtom(scale, fields[1:])
	if err != nil {
		return nil, 0, err
	}
	return common.Random(values), shift + 1, nil
}
