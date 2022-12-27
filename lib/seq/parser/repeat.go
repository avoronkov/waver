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

	keyIndex := fmt.Sprintf("repeat-idx:%v", line.Num)
	keyValue := fmt.Sprintf("repeat-value:%v", line.Num)
	return common.Repeat(keyIndex, keyValue, times, fn), shift + shift2 + 1, nil
}
