package parser

import (
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

func New(opts ...func(*Parser)) *Parser {
	p := &Parser{
		userFunctions: make(map[string]UserFunction),
		userSignalers: make(map[string][]types.Signaler),
		tempo:         120,
	}

	for _, o := range opts {
		o(p)
	}

	//  Init mod parsers
	p.ModParsers = map[lexer.Token]ModParser{
		lexer.ColonToken{}:        makeSingleArgModParser(":", common.Every),
		lexer.PlusToken{}:         makeSingleArgModParser("+", common.Shift),
		lexer.MinusToken{}:        makeSingleArgModParser("-", common.ShiftLeft),
		lexer.LessToken{}:         makeSingleArgModParser("<", common.Before),
		lexer.GreaterToken{}:      makeSingleArgModParser(">", common.After),
		lexer.IdentToken("bits"):  makeSingleArgModParser("bits", common.Bits),
		lexer.IdentToken("eucl"):  makeTwoArgsModParser("eucl", common.EuclideanFirst),
		lexer.IdentToken("eucl'"): makeTwoArgsModParser("eucl'", common.EuclideanLast),
		lexer.MultiplyToken{}:     parseTimesModifier,
	}

	// Init pragma parsers
	p.PragmaParsers = map[string]pragmaParser{
		"tempo":    parseTempo,
		"sample":   parseSample,
		"wave":     parseWave,
		"inst":     parseWave,
		"form":     parseForm,
		"lagrange": parseLagrange,
		"filter":   parseFilter,
		"stop":     parseStopPragma,
		"srand":    parseSrandPragma,
		"scale":    parseScalePragma,
	}

	p.FuncParsers = map[lexer.Token]FunctionParser{
		lexer.IdentToken("seq"):    parseSequence,
		lexer.AtToken{}:            parseSequence,
		lexer.IdentToken("rand"):   makeSingleArgValueFnParser("rand", common.Random),
		lexer.AmpersandToken{}:     makeSingleArgValueFnParser("rand", common.Random),
		lexer.IdentToken("up"):     parseUpDown,
		lexer.IdentToken("down"):   parseUpDown,
		lexer.IdentToken("repeat"): parseRepeat,
		lexer.MultiplyToken{}:      parseRepeat,
		lexer.IdentToken("concat"): makeSingleArgValueFnParser("concat", common.Concat),
		lexer.IdentToken("loop"):   parseLoop,
	}

	return p
}

func WithSeq(seq Seq) func(*Parser) {
	return func(p *Parser) {
		p.seq = seq
	}
}

func WithScale(scale notes.Scale) func(*Parser) {
	return func(p *Parser) {
		p.scale = scale
	}
}

func WithScaleSetters(s ...ScaleSetter) func(*Parser) {
	return func(p *Parser) {
		p.scaleSetters = append(p.scaleSetters, s...)
	}
}

func WithTempoSetter(setter TempoSetter) func(*Parser) {
	return func(p *Parser) {
		p.tempoSetters = append(p.tempoSetters, setter)
	}
}

func WithInstrumentSet(set InstrumentSet) func(*Parser) {
	return func(p *Parser) {
		p.instSet = set
	}
}

func WithFileInput(file string) func(p *Parser) {
	return func(p *Parser) {
		p.file = file
	}
}
