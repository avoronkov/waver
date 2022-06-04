package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ValueFnParser = func(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error)

// seq [ 1 2 3 ]
func parseSequence(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", line)
	}

	values, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	key := fmt.Sprintf("seq:%v", line.Num)
	var idx *common.Index
	if i, ok := line.GlobalCtx[key].(*common.Index); ok {
		idx = i
	} else {
		idx = new(common.Index)
		line.GlobalCtx[key] = idx
	}
	return common.Sequence(idx, values), shift + 1, nil
}
