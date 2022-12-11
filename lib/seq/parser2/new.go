package parser2

import (
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
)

func New(opts ...func(*Parser)) *Parser {
	p := &Parser{
		modParsers:    make(map[lexer.Token]ModParser),
		userFunctions: make(map[string]parser.UserFunction),
	}

	for _, o := range opts {
		o(p)
	}

	//  Init mod parsers
	p.modParsers[lexer.ColonToken{}] = makeSingleArgModParser(":", common.Every)
	p.modParsers[lexer.PlusToken{}] = makeSingleArgModParser("+", common.Shift)
	p.modParsers[lexer.IdentToken{Value: "bits"}] = makeSingleArgModParser("bits", common.Bits)
	p.modParsers[lexer.IdentToken{Value: "eucl"}] = makeTwoArgsModParser("eucl", common.EuclideanFirst)
	p.modParsers[lexer.IdentToken{Value: "eucz"}] = makeTwoArgsModParser("eucz", common.EuclideanLast)

	return p
}

func WithSeq(seq parser.Seq) func(*Parser) {
	return func(p *Parser) {
		p.seq = seq
	}
}
