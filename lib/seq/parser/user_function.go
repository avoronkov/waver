package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

func (p *Parser) parseUserFunction(line *LineCtx, name, argName string, body types.ValueFn) (types.ValueFn, int, error) {
	if line.Len() < 2 {
		return nil, 0, fmt.Errorf("Not enough arguments for '%v': %v", name, line)
	}

	arg, shift, err := p.parseAtom(line.Shift(1))
	if err != nil {
		return nil, 0, err
	}

	return common.UserFunction(argName, arg, body), shift + 1, nil
}
