package parser

import (
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

func (p *Parser) parseUserFunction(lx *lexer.Lexer, name, argName string, body types.ValueFn) (types.ValueFn, error) {
	arg, err := p.parseAtom(lx)
	if err != nil {
		return nil, err
	}

	return common.UserFunction(argName, arg, body), nil
}
