package parser2

import (
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

func (p *Parser) parseList(lx *lexer.Lexer) (types.ValueFn, error) {
	atoms := []types.ValueFn{}
L:
	for {
		token, err := lx.Top()
		if err != nil {
			return nil, err
		}
		if _, ok := token.(lexer.EolToken); ok {
			_, _ = lx.Pop()
			continue L
		}
		if _, ok := token.(lexer.RSquareBracket); ok {
			_, _ = lx.Pop()
			break L
		}

		atom, err := p.parseAtom(lx)
		if err != nil {
			return nil, err
		}
		atoms = append(atoms, atom)
	}

	return common.Lst(atoms...), nil
}
