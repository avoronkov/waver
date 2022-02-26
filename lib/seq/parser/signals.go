package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type SigParser = func(fields []string) (types.Signaler, int, error)

// { 2 A4 }
func parseSignal(fields []string) (types.Signaler, int, error) {
	l := len(fields)
	if l < 4 {
		return nil, 0, fmt.Errorf("Not enough arguments for signal: %v", fields)
	}

	// parse instrument

	// parse note

	// parse amplitude

	// parse duration

	panic("NIY")
}

func parseRawSignal(fields []string) (types.Signaler, int, error) {
	if len(fields) < 1 {
		return nil, 0, fmt.Errorf("No arguments for raw signal")
	}
	raw := fields[0]
	rawLen := len(raw)
	if rawLen > 2 && raw[0] == '\'' && raw[rawLen-1] == '\'' {
		return common.Sig(raw[1 : rawLen-1]), 1, nil
	}
	return nil, 0, fmt.Errorf("Incorrect raw signal: %q", raw)
}
