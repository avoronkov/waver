package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

func makeMusParser(name string, shifts ...int64) ValueFnParser {
	return func(p *Parser, line *LineCtx) (types.ValueFn, int, error) {
		if line.Len() < 2 {
			return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
		}
		x, shift, err := p.parseAtom(line.Shift(1))
		if err != nil {
			return nil, 0, err
		}
		return common.ChordFn(x, shifts...), shift + 1, nil
	}
}
