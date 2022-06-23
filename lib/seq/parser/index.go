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
		"up":     parseUpDown,
		"down":   parseUpDown,
		"repeat": parseRepeat,
		"concat": MakeSingleArgValueFnParser("concat", common.Concat),
	}
}
