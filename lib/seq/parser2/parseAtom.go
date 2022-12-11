package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

func (p *Parser) parseAtom(lx *lexer.Lexer) (types.ValueFn, error) {
	token, err := lx.Pop()
	if err != nil {
		return nil, err
	}
	switch a := token.(type) {
	case lexer.NumberToken:
		return common.Const(a.Num), nil
	case lexer.IdentToken:
		return common.Var(a.Value), nil
	case lexer.LSquareBracket:
		return p.parseList(lx)
	}
	return nil, fmt.Errorf("Unexpected token while parsing atom: %v", token)
}
