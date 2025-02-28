package parser

import (
	"maps"
	"slices"

	"github.com/avoronkov/waver/lib/midisynth/waves"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
	"github.com/avoronkov/waver/lib/utils"
	"github.com/avoronkov/waver/static"
)

func New(opts ...func(*Parser)) *Parser {
	p := &Parser{
		userFunctions:       make(map[string]UserFunction),
		userSignalers:       make(map[string][]types.Signaler),
		userFilters:         make(map[string]string),
		instrumentVariables: utils.NewSet[string](),
		tempo:               120,
	}

	for _, o := range opts {
		o(p)
	}

	//  Init mod parsers
	p.ModParsers = map[lexer.Token]ModParser{
		lexer.ColonToken{}: {
			Usage: "<n[int]> | <n, m, l...>",
			Parse: makeSingleArgModParser(":", common.Every),
		},
		lexer.PlusToken{}: {
			Usage: "<n[int]>",
			Parse: makeSingleArgModParser("+", common.Shift),
		},
		lexer.MinusToken{}: {
			Usage: "<n[int]>",
			Parse: makeSingleArgModParser("-", common.ShiftLeft),
		},
		lexer.LessToken{}: {
			Usage: "<frame[int]>",
			Parse: makeSingleArgModParser("<", common.Before),
		},
		lexer.GreaterToken{}: {
			Usage: "<frame[int]>",
			Parse: makeSingleArgModParser(">", common.After),
		},
		lexer.IdentToken("bits"): {
			Usage: "<a, b, c... totalBits[int]>",
			Parse: makeSingleArgModParser("bits", common.Bits),
		},
		lexer.IdentToken("eucl"): {
			Usage: "<pulses> <steps> (pulses <= steps)",
			Parse: makeTwoArgsModParser("eucl", common.EuclideanFirst),
		},
		lexer.IdentToken("eucl'"): {
			Usage: "<pulses> <steps> (pulses <= steps)",
			Parse: makeTwoArgsModParser("eucl'", common.EuclideanLast),
		},
		lexer.MultiplyToken{}: {
			Usage:      "^ <times[int]>",
			Parse:      parseTimesModifier,
			Deprecated: true,
		},
		lexer.CaretToken{}: {
			Usage: "<times[int]>",
			Parse: parseTimesModifier,
		},
	}

	// Init pragma parsers
	p.PragmaParsers = map[string]pragmaParser{
		"tempo": {
			Usage: "<int>",
			Desc: `Specify tempo in BMP (bits-per-minute).
			Each "bit" is subdivided into 4 frames.`,
			Parse: parseTempo,
		},
		"sample": {
			Usage: `<Name> "<sample-file>"`,
			Desc: `Define an instrument using a sample.
			Waver has a number of builtin percussion samples which can be used.`,
			Examples: static.ListFiles(),
			Parse:    parseSample,
		},
		"wave": {
			Usage:    `<Name> "<waveform>"`,
			Desc:     `Define an instrument using basic waveform.`,
			Parse:    parseWave,
			Examples: slices.Sorted(maps.Keys(waves.Waves)),
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
			Desc:  `Stop processing file after the specified frame.`,
			Parse: parseStopPragma,
		},
		"srand": {
			Usage: `<int>`,
			Desc:  `Specify seed for PRNG.`,
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
			Desc:  `Convert a list into a sequence generator which returns each item of the list at a time.`,
			Parse: parseSequence,
		},
		lexer.AtToken{}: {
			Usage: "[ 1 2 3 ]",
			Desc: `Convert a list into a sequence generator which returns each item of the list at a time.
			Shorthand for function "seq".`,
			Parse: parseSequence,
		},
		lexer.IdentToken("rand"): {
			Usage: "(&) [ 1 2 3 ]",
			Desc:  `Take a random element from a list.`,
			Parse: makeSingleArgValueFnParser("rand", common.Random),
		},
		lexer.AmpersandToken{}: {
			Usage: "[ 1 2 3 ]",
			Desc: `Take a random element from a list.
			Shorthand for function "rand"`,
			Parse: makeSingleArgValueFnParser("rand", common.Random),
		},
		lexer.IdentToken("up"): {
			Usage: "<int> (<Note> | <Chord>)",
			Desc:  `Increase a pitch of a note (or a chord).`,
			Parse: parseUpDown,
		},
		lexer.IdentToken("down"): {
			Usage: "<int> (<Note> | <Chord>)",
			Desc:  `Decrease a pitch of a note (or a chord).`,
			Parse: parseUpDown,
		},
		lexer.IdentToken("repeat"): {
			Usage: "(*) <n[int]> <sequence...>",
			Desc:  `Repeat each element of a sequence n times.`,
			Parse: parseRepeat,
		},
		lexer.MultiplyToken{}: {
			Usage: "<int> <sequence...>",
			Desc: `Repeat each element of a sequence n times.
			Shorthand for function "repeat"`,
			Parse: parseRepeat,
		},
		lexer.IdentToken("concat"): {
			Usage: "[ list1, list2 ... ]",
			Desc:  `Concatenate lists.`,
			Parse: makeSingleArgValueFnParser("concat", common.Concat),
		},
		lexer.IdentToken("loop"): {
			Usage: "<size[int]> <sequence...>",
			Desc:  `Take "size" elements from a sequence and repeat them over and over.`,
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
