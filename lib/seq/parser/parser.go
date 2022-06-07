package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/common"
	"github.com/avoronkov/waver/lib/seq/types"
)

var modParsers = map[string]ModParser{
	":": parseEvery,
	"+": parseShift,
	"-": parseShift,
	"<": parseBefore,
	">": parseAfter,
}

var sigParsers = map[string]SigParser{
	"":  parseRawSignal,
	"{": parseSignal,
}

var valueFnParser map[string]ValueFnParser

func init() {
	valueFnParser = map[string]ValueFnParser{
		"seq":    parseSequence,
		"rand":   parseRandom,
		"maj":    makeMusParser("maj", 0, 4, 7),
		"maj7":   makeMusParser("maj", 0, 4, 7, 11),
		"maj9":   makeMusParser("maj", 0, 4, 7, 11, 14),
		"min":    makeMusParser("min", 0, 3, 7),
		"min7":   makeMusParser("min", 0, 3, 7, 10),
		"min9":   makeMusParser("min", 0, 3, 7, 10, 14),
		"up":     parseUpDown,
		"down":   parseUpDown,
		"repeat": parseRepeat,
	}
}

type Parser struct {
	file string
	seq  Seq

	modParsers map[string]ModParser
	sigParsers map[string]SigParser

	scale notes.Scale

	globalCtx map[string]interface{}
}

func New(seq Seq, scale notes.Scale, opts ...func(*Parser)) *Parser {
	p := &Parser{
		seq:        seq,
		scale:      scale,
		modParsers: modParsers,
		sigParsers: sigParsers,
		globalCtx:  map[string]interface{}{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Parser) ParseData(data []byte) error {
	return p.parseReader(bytes.NewReader(data))
}

func (p *Parser) parseReader(reader io.Reader) error {
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanLines)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		text := sc.Text()
		if text == "" || text[0] == '#' {
			continue
		}
		if err := p.parseLine(lineNum, text); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("Scanner failed: %v", err)
	}
	return p.seq.Commit()
}

func (p *Parser) parse() error {
	f, err := os.Open(p.file)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.parseReader(f)
}

// : 5 -> { 4 E3 1 5 }
func (p *Parser) parseLine(num int, line string) error {
	fields := strings.Fields(line)
	lineCtx := &LineCtx{
		Num:       num,
		Fields:    fields,
		GlobalCtx: p.globalCtx,
	}
	if idx := stringsFind(fields, "->"); idx >= 0 {
		// modifiers -> signals
		modCtx := &LineCtx{
			Num:       num,
			Fields:    fields[:idx],
			GlobalCtx: p.globalCtx,
		}
		mods, err := p.parseModifiers(modCtx)
		if err != nil {
			return err
		}
		signals, err := p.parseSignal(lineCtx.Shift(idx + 1))
		if err != nil {
			return err
		}
		for _, sig := range signals {
			x := common.Chain(sig, mods...)
			p.seq.Add(x)
		}
	} else if len(fields) >= 2 && fields[1] == "=" {
		// var = atom
		// TODO check shift
		vfn, _, err := parseAtom(p.scale, lineCtx.Shift(2))
		if err != nil {
			return err
		}
		p.seq.Assign(fields[0], vfn)
	} else {
		log.Printf("[WARNING] Skipping line: %q", line)
	}
	return nil
}

func (p *Parser) parseModifiers(line *LineCtx) (result []types.Modifier, err error) {
	l := line.Len()
	for i := 0; i < l; {
		if parser, ok := p.modParsers[line.Fields[i]]; ok {
			mod, shift, err := parser(p.scale, line.Shift(i))
			if err != nil {
				return nil, err
			}
			result = append(result, mod)
			i += shift
		} else {
			return nil, fmt.Errorf("Unknown modifier: %v", line.Fields[i])
		}
	}
	// OK
	return
}

func (p *Parser) parseSignal(line *LineCtx) (result []types.Signaler, err error) {
	l := line.Len()
	for i := 0; i < l; {
		parser, ok := p.sigParsers[line.Fields[i]]
		if !ok {
			parser, ok = p.sigParsers[""]
			if !ok {
				return nil, fmt.Errorf("Don't know how to parse signal: %q", line.Fields[i])
			}
		}
		sig, shift, err := parser(p.scale, line.Shift(i))
		if err != nil {
			return nil, err
		}
		result = append(result, sig)
		i += shift
	}
	return
}

func stringsFind(list []string, needle string) int {
	for i, s := range list {
		if s == needle {
			return i
		}
	}
	return -1
}
