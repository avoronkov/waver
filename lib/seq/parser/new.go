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
		"tempo": {
			Usage: "<int>",
			Parse: parseTempo,
		},
		"sample": {
			Usage: `<Name> "<sample-file>"`,
			Parse: parseSample,
		},
		"wave": {
			Usage: `<Name> "<wave-form>"`,
			Parse: parseWave,
		},
		"inst": {
			Parse:      parseWave,
			Deprecated: true,
		},
		"form": {
			Usage: `<Name> "<file-path>"`,
			Parse: parseForm,
		},
		"lagrange": {
			Usage: `<Name> "<file-path>"`,
			Parse: parseLagrange,
		},
		"filter": {
			Usage: ``,
			Parse: parseFilter,
		},
		"stop": {
			Usage: `<frame[int]>`,
			Parse: parseStopPragma,
		},
		"srand": {
			Usage: `<int>`,
			Parse: parseSrandPragma,
		},
		"scale": {
			Usage: "edo12|edo19",
			Parse: parseScalePragma,
		},
	}

	p.FuncParsers = map[lexer.Token]FunctionParser{
		lexer.IdentToken("seq"): {
			Usage: "(@) [ 1 2 3 ]",
			Parse: parseSequence,
		},
		lexer.AtToken{}: {
			Usage: "[ 1 2 3 ]",
			Parse: parseSequence,
		},
		lexer.IdentToken("rand"): {
			Usage: "(&) [ 1 2 3 ]",
			Parse: makeSingleArgValueFnParser("rand", common.Random),
		},
		lexer.AmpersandToken{}: {
			Usage: "[ 1 2 3 ]",
			Parse: makeSingleArgValueFnParser("rand", common.Random),
		},
		lexer.IdentToken("up"): {
			Usage: "<int> <Note>",
			Parse: parseUpDown,
		},
		lexer.IdentToken("down"): {
			Usage: "<int> <Note>",
			Parse: parseUpDown,
		},
		lexer.IdentToken("repeat"): {
			Usage: "(*) <int> <sequence...>",
			Parse: parseRepeat,
		},
		lexer.MultiplyToken{}: {
			Usage: "<int> <sequence...>",
			Parse: parseRepeat,
		},
		lexer.IdentToken("concat"): {
			Usage: "[ list1, list2 ... ]",
			Parse: makeSingleArgValueFnParser("concat", common.Concat),
		},
		lexer.IdentToken("loop"): {
			Usage: "<size[int]> <sequence...>",
			Parse: parseLoop,
		},
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
