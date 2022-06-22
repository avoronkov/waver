package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

func parseRepeat(p *Parser, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 3 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'repeat': %v", line)
	}

	times, shift, err := p.parseAtom(line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	fn, shift2, err := p.parseAtom(line.Shift(1 + shift))
	if err != nil {
		return nil, 0, err
	}

	key := fmt.Sprintf("repeat-idx:%v", line.Num)
	var idx *common.Index
	if i, ok := p.globalCtx[key].(*common.Index); ok {
		idx = i
	} else {
		idx = new(common.Index)
		p.globalCtx[key] = idx
	}

	key = fmt.Sprintf("repeat-value:%v", line.Num)
	var holder *common.ValueHolder
	if i, ok := p.globalCtx[key].(*common.ValueHolder); ok {
		holder = i
	} else {
		holder = new(common.ValueHolder)
		p.globalCtx[key] = holder
	}
	return common.Repeat(idx, holder, times, fn), shift + shift2 + 1, nil
}
