package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

func makeMusParser(name string, shifts ...int64) ValueFnParser {
	return func(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
		if line.Len() < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
		}
		x, shift, err := parseAtom(scale, line.Shift(1))
		if err != nil {
			return nil, 0, err
		}
		return common.ChordFn(x, shifts...), shift + 1, nil
	}
}
