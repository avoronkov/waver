package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

func (p *Parser) parseAtom(lx *lexer.Lexer) (types.ValueFn, error) {
	var prev types.ValueFn
	elements := []types.ValueFn{}
	for {
		var elem types.ValueFn
		if prev != nil {
			// check relative element
			tok, err := lx.Top()
			if err != nil {
				return nil, err
			}
			switch tok.(type) {
			case lexer.PlusToken, lexer.MinusToken:
				lx.Drop()
				shift, err := lx.Pop()
				if err != nil {
					return nil, err
				}
				numShift, ok := shift.(lexer.NumberToken)
				if ok {
					elem = common.Up(common.Const(int64(numShift)), prev, tok == lexer.MinusToken{})
				} else {
					return nil, fmt.Errorf("Unexpected token: %v (%T)", shift, shift)
				}
			}
		}
		if elem == nil {
			var err error
			elem, err = p.parseElement(lx)
			if err != nil {
				return nil, err
			}
			prev = elem
		}
		elements = append(elements, elem)

		// Check if next token is coma (,)
		tok, err := lx.Top()
		if err != nil {
			return nil, err
		}
		if _, ok := tok.(lexer.ComaToken); !ok {
			break
		}
		_, _ = lx.Pop()
	}
	if len(elements) == 1 {
		return elements[0], nil
	}
	return common.Lst(elements...), nil
}

func (p *Parser) parseElement(lx *lexer.Lexer) (types.ValueFn, error) {
	token, err := lx.Pop()
	if err != nil {
		return nil, err
	}
	switch a := token.(type) {
	case lexer.NumberToken:
		return common.Const(int64(a)), nil
	case lexer.FloatToken:
		return common.FloatConst(float64(a)), nil
	case lexer.IdentToken:
		sa := string(a)
		if n, ok := p.scale.Parse(sa); ok {
			return common.Const(int64(n.Num)), nil
		}
		if fnp, ok := p.FuncParsers[token]; ok {
			return fnp(p, lx, sa)
		}
		if udf, ok := p.userFunctions[sa]; ok {
			return p.parseUserFunction(lx, udf.Name, udf.Arg, udf.Fn)
		}
		return common.Var(sa), nil
	case lexer.LSquareBracket:
		return p.parseList(lx)
	default:
		if fnp, ok := p.FuncParsers[token]; ok {
			return fnp(p, lx, token.String())
		}
	}
	return nil, fmt.Errorf("Unexpected token while parsing atom: %v", token)
}
