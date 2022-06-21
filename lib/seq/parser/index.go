package parser

import "github.com/avoronkov/waver/lib/seq/common"

var modParsers = map[string]ModParser{
	":":    makeSingleArgModParser("Every", common.Every),
	"+":    parseShift,
	"-":    parseShift,
	"<":    makeSingleArgModParser("Before", common.Before),
	">":    makeSingleArgModParser("After", common.After),
	"bits": makeSingleArgModParser("Bits", common.Bits),
}

var sigParsers = map[string]SigParser{
	"":  parseRawSignal,
	"{": parseSignal,
}

var valueFnParser map[string]ValueFnParser

func init() {
	valueFnParser = map[string]ValueFnParser{
		"seq":    parseSequence,
		"rand":   MakeSingleArgValueFnParser("rand", common.Random),
		"maj":    makeMusParser("maj", 0, 4, 7),
		"maj7":   makeMusParser("maj", 0, 4, 7, 11),
		"maj9":   makeMusParser("maj", 0, 4, 7, 11, 14),
		"min":    makeMusParser("min", 0, 3, 7),
		"min7":   makeMusParser("min", 0, 3, 7, 10),
		"min9":   makeMusParser("min", 0, 3, 7, 10, 14),
		"up":     parseUpDown,
		"down":   parseUpDown,
		"repeat": parseRepeat,
		"concat": MakeSingleArgValueFnParser("concat", common.Concat),
	}
}
