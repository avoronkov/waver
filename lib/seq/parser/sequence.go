package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type ValueFnParser = func(p *Parser, line *LineCtx) (types.ValueFn, int, error)

// seq [ 1 2 3 ]
func parseSequence(p *Parser, line *LineCtx) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for 'seq': %v", line)
	}

	values, shift, err := p.parseAtom(line.Shift(1))
	if err != nil {
		return nil, 0, err
	}
	key := fmt.Sprintf("seq:%v", line.Num)
	var idx *common.Index
	if i, ok := p.globalCtx[key].(*common.Index); ok {
		idx = i
	} else {
		idx = new(common.Index)
		p.globalCtx[key] = idx
	}
	return common.Sequence(idx, values), shift + 1, nil
}
