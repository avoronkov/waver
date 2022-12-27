package parser2

import (
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
)

func New(opts ...func(*Parser)) *Parser {
	p := &Parser{
		userFunctions: make(map[string]parser.UserFunction),
	}

	for _, o := range opts {
		o(p)
	}

	//  Init mod parsers
	p.modParsers = map[lexer.Token]ModParser{
		lexer.ColonToken{}:        makeSingleArgModParser(":", common.Every),
		lexer.PlusToken{}:         makeSingleArgModParser("+", common.Shift),
		lexer.LessToken{}:         makeSingleArgModParser("<", common.Before),
		lexer.GreaterToken{}:      makeSingleArgModParser(">", common.After),
		lexer.IdentToken("bits"):  makeSingleArgModParser("bits", common.Bits),
		lexer.IdentToken("eucl"):  makeTwoArgsModParser("eucl", common.EuclideanFirst),
		lexer.IdentToken("eucl'"): makeTwoArgsModParser("eucl'", common.EuclideanLast),
		lexer.MultiplyToken{}:     parseTimesModifier,
	}

	// Init pragma parsers
	p.pragmaParsers = map[string]pragmaParser{
		"tempo":  parseTempo,
		"sample": parseSample,
		"wave":   parseWave,
		"inst":   parseWave,
		"filter": parseFilter,
	}

	p.funcParsers = map[string]FunctionParser{
		"seq":    parseSequence,
		"rand":   makeSingleArgValueFnParser("rand", common.Random),
		"up":     parseUpDown,
		"down":   parseUpDown,
		"repeat": parseRepeat,
		"concat": makeSingleArgValueFnParser("concat", common.Concat),
	}

	return p
}

func WithSeq(seq parser.Seq) func(*Parser) {
	return func(p *Parser) {
		p.seq = seq
	}
}

func WithScale(scale notes.Scale) func(*Parser) {
	return func(p *Parser) {
		p.scale = scale
	}
}

func WithTempoSetter(setter parser.TempoSetter) func(*Parser) {
	return func(p *Parser) {
		p.tempoSetters = append(p.tempoSetters, setter)
	}
}

func WithInstrumentSet(set parser.InstrumentSet) func(*Parser) {
	return func(p *Parser) {
		p.instSet = set
	}
}

func WithFileInput(file string) func(p *Parser) {
	return func(p *Parser) {
		p.file = file
	}
}
