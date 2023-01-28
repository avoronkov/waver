package parser2

import (
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
	"github.com/avoronkov/waver/lib/seq/types"
)

func New(opts ...func(*Parser)) *Parser {
	p := &Parser{
		userFunctions: make(map[string]parser.UserFunction),
		userSignalers: make(map[string][]types.Signaler),
		tempo:         120,
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

	p.funcParsers = map[lexer.Token]FunctionParser{
		lexer.IdentToken("seq"):    parseSequence,
		lexer.AtToken{}:            parseSequence,
		lexer.IdentToken("rand"):   makeSingleArgValueFnParser("rand", common.Random),
		lexer.AmpersandToken{}:     makeSingleArgValueFnParser("rand", common.Random),
		lexer.IdentToken("up"):     parseUpDown,
		lexer.IdentToken("down"):   parseUpDown,
		lexer.IdentToken("repeat"): parseRepeat,
		lexer.MultiplyToken{}:      parseRepeat,
		lexer.IdentToken("concat"): makeSingleArgValueFnParser("concat", common.Concat),
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
