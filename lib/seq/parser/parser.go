package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/avoronkov/waver/lib/seq/common"
	"gitlab.com/avoronkov/waver/lib/seq/types"
	"gitlab.com/avoronkov/waver/lib/watch"
)

var modParsers = map[string]ModParser{
	":": parseEvery,
	"+": parseShift,
	"-": parseShift,
}
var sigParsers = map[string]SigParser{
	"":  parseRawSignal,
	"{": parseSignal,
}

type Parser struct {
	file string
	seq  Seq

	modParsers map[string]ModParser
	sigParsers map[string]SigParser
}

func New(file string, seq Seq) *Parser {
	return &Parser{
		file:       file,
		seq:        seq,
		modParsers: modParsers,
		sigParsers: sigParsers,
	}
}

func (p *Parser) Start(wtch bool) error {
	if err := p.parse(); err != nil {
		return err
	}
	if wtch {
		err := watch.OnFileUpdate(p.file, func() {
			if err := p.parse(); err != nil {
				log.Printf("Parsing %v failed: %v", p.file, err)
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parse() error {
	f, err := os.Open(p.file)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		text := sc.Text()
		if text == "" || text[0] == '#' {
			continue
		}
		if err := p.parseLine(text); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf("Scanner failed: %v", err)
	}
	return p.seq.Commit()
}

// : 5 -> { 4 E3 1 5 }
func (p *Parser) parseLine(line string) error {
	fields := strings.Fields(line)
	if idx := stringsFind(fields, "->"); idx >= 0 {
		// modifiers -> signals
		mods, err := p.parseModifiers(fields[:idx])
		if err != nil {
			return err
		}
		signals, err := p.parseSignal(fields[idx+1:])
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
		vfn, _, err := parseAtom(fields[2:])
		if err != nil {
			return err
		}
		p.seq.Assign(fields[0], vfn)
	} else {
		log.Printf("[WARNING] Skipping line: %q", line)
	}
	return nil
}

func (p *Parser) parseModifiers(fields []string) (result []types.Modifier, err error) {
	l := len(fields)
	for i := 0; i < l; {
		if parser, ok := p.modParsers[fields[i]]; ok {
			mod, shift, err := parser(fields[i:])
			if err != nil {
				return nil, err
			}
			result = append(result, mod)
			i += shift
		} else {
			return nil, fmt.Errorf("Unknown modifier: %v", fields[i])
		}
	}
	// OK
	return
}

func (p *Parser) parseSignal(fields []string) (result []types.Signaler, err error) {
	l := len(fields)
	for i := 0; i < l; {
		parser, ok := p.sigParsers[fields[i]]
		if !ok {
			parser, ok = p.sigParsers[""]
			if !ok {
				return nil, fmt.Errorf("Don't know how to parse signal: %q", fields[i])
			}
		}
		sig, shift, err := parser(fields[i:])
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
