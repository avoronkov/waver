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

	// skip '{'
	shift := 1
	// parse instrument
	in, sh, err := parseAtom(fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	// parse note
	nt, sh, err := parseAtom(fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if fields[shift] == "}" {
		return common.Note(in, nt), shift + 1, nil
	}

	// parse amplitude
	amp, sh, err := parseAtom(fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if fields[shift] == "}" {
		return common.Note(in, nt, common.NoteAmp(amp)), shift + 1, nil
	}

	// parse duration
	dur, sh, err := parseAtom(fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	return common.Note(in, nt, common.NoteAmp(amp), common.NoteDur(dur)), shift + 1, nil
}

func parseRawSignal(fields []string) (types.Signaler, int, error) {
	if len(fields) < 1 {
		return nil, 0, fmt.Errorf("No arguments for raw signal")
	}
	raw := fields[0]
	rawLen := len(raw)
	if rawLen > 2 && ((raw[0] == '\'' && raw[rawLen-1] == '\'') || (raw[0] == '"' && raw[rawLen-1] == '"')) {
		sig, err := common.Sig(raw[1 : rawLen-1])
		if err != nil {
			return nil, 0, err
		}
		return sig, 1, nil
	}
	return nil, 0, fmt.Errorf("Incorrect raw signal: %q", raw)
}
