package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

func parseRepeat(scale notes.Scale, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 3 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'repeat': %v", line)
	}

	times, shift, err := parseAtom(scale, line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	fn, shift2, err := parseAtom(scale, line.Shift(1+shift))
	if err != nil {
		return nil, 0, err
	}
	key := fmt.Sprintf("repeat:%v", line.Num)
	var idx *common.Index
	if i, ok := line.GlobalCtx[key].(*common.Index); ok {
		idx = i
	} else {
		idx = new(common.Index)
		line.GlobalCtx[key] = idx
	}
	return common.Repeat(idx, times, fn), shift + shift2 + 1, nil
}
