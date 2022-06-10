package parser

import (
	"fmt"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

type SigParser = func(scale notes.Scale, line *LineCtx) (types.Signaler, int, error)

// { 2 A4 }
func parseSignal(scale notes.Scale, line *LineCtx) (types.Signaler, int, error) {
	l := line.Len()
	if l < 3 {
		return nil, 0, fmt.Errorf("Not enough arguments for signal: %v", line)
	}

	// skip '{'
	shift := 1
	// parse instrument
	in, sh, err := parseAtom(scale, line.Shift(shift))
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if line.Fields[shift] == "}" {
		return common.Note(scale, in, common.StrConst("_")), shift + 1, nil
	}

	// parse note
	nt, sh, err := parseAtom(scale, line.Shift(shift))
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if line.Fields[shift] == "}" {
		return common.Note(scale, in, nt), shift + 1, nil
	}

	// parse amplitude
	amp, sh, err := parseAtom(scale, line.Shift(shift))
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	if line.Fields[shift] == "}" {
		return common.Note(scale, in, nt, common.NoteAmp(amp)), shift + 1, nil
	}

	// parse duration
	dur, sh, err := parseAtom(scale, line.Shift(shift))
	if err != nil {
		return nil, 0, err
	}
	shift += sh

	return common.Note(scale, in, nt, common.NoteAmp(amp), common.NoteDur(dur)), shift + 1, nil
}

func parseRawSignal(scale notes.Scale, line *LineCtx) (types.Signaler, int, error) {
	if line.Len() < 1 {
		return nil, 0, fmt.Errorf("No arguments for raw signal")
	}
	raw := line.Fields[0]
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
