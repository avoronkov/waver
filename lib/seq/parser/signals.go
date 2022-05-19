package parser

import (
	"fmt"

	"gitlab.com/avoronkov/waver/lib/notes"
	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
)

type SigParser = func(scale notes.Scale, fields []string) (types.Signaler, int, error)

// { 2 A4 }
func parseSignal(scale notes.Scale, fields []string) (types.Signaler, int, error) {
	l := len(fields)
	if l < 4 {
		return nil, 0, fmt.Errorf("Not enough arguments for signal: %v", fields)
	}

	// skip '{'
	shift := 1
	// parse instrument
	in, sh, err := parseAtom(scale, fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	// parse note
	nt, sh, err := parseAtom(scale, fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if fields[shift] == "}" {
		return common.Note(scale, in, nt), shift + 1, nil
	}

	// parse amplitude
	amp, sh, err := parseAtom(scale, fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if fields[shift] == "}" {
		return common.Note(scale, in, nt, common.NoteAmp(amp)), shift + 1, nil
	}

	// parse duration
	dur, sh, err := parseAtom(scale, fields[shift:])
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	return common.Note(scale, in, nt, common.NoteAmp(amp), common.NoteDur(dur)), shift + 1, nil
}

func parseRawSignal(scale notes.Scale, fields []string) (types.Signaler, int, error) {
	if len(fields) < 1 {
		return nil, 0, fmt.Errorf("No arguments for raw signal")
	}
	raw := fields[0]
	rawLen := len(raw)
	if rawLen > 2 && ((raw[0] == '\'' && raw[rawLen-1] == '\'') || (raw[0] == '"' && raw[rawLen-1] == '"')) {
		// TODO avoronkov fix here
		sig, err := common.Sig(raw[1 : rawLen-1])
		if err != nil {
			return nil, 0, err
		}
		return sig, 1, nil
	}
	return nil, 0, fmt.Errorf("Incorrect raw signal: %q", raw)
}
