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
		if n, ok := p.scale.Parse(a.Value); ok {
			return common.Const(int64(n.Num)), nil
		}
		if fnp, ok := p.funcParsers[a.Value]; ok {
			return fnp(p, lx, a.Value)
		}
		return common.Var(a.Value), nil
	case lexer.LSquareBracket:
		return p.parseList(lx)
	}
	return nil, fmt.Errorf("Unexpected token while parsing atom: %v", token)
}
