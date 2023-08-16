package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/avoronkov/waver/lib/forth"
)

type Parser struct {
	forth *forth.Forth
}

func (p *Parser) Parse(r io.Reader) (*forth.Forth, error) {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)
	tokens := []string{}
	for sc.Scan() {
		line := sc.Text()
		log.Printf("Parsing line: %v", line)
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			log.Printf("Skipping line: %v", line)
			continue
		}
		tokens = append(tokens, strings.Fields(line)...)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	err := p.parseTokens(tokens)
	if err != nil {
		return nil, err
	}
	return p.forth, nil
}

func (p *Parser) parseTokens(tokens []string) error {
	program := []forth.StackFn{}
	idx := 0
	l := len(tokens)
	for idx < l {
		log.Printf("Parsing token %v at %v", tokens[idx], idx)
		newIdx, err := p.parseDefine(tokens, idx)
		if err != nil {
			return err
		}
		defineParsed := idx != newIdx
		idx = newIdx
		// parse possible define again
		if defineParsed {
			continue
		}
		fn, newIdx, err := p.parseAtom(tokens, idx)
		if err != nil {
			return err
		}
		program = append(program, fn)
		idx = newIdx
	}

	forth.WithProgram(program)(p.forth)

	return nil
}

func ParseFile(file string) (*forth.Forth, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (*forth.Forth, error) {
	parser := &Parser{
		forth: forth.New(),
	}
	return parser.Parse(r)
}

var Funcs = map[string]forth.StackFn{
	"+":     forth.Plus,
	"-":     forth.Minus,
	"*":     forth.Multiply,
	"dup":   forth.Dup,
	"drop":  forth.Drop,
	"swap":  forth.Swap,
	"over":  forth.Over,
	"rot":   forth.Rot,
	"top":   forth.ShowTop,
	"stack": forth.ShowStack,
	"and":   forth.And,
	"or":    forth.Or,
	"not":   forth.Not,
}

func (p *Parser) parseAtom(tokens []string, idx int) (forth.StackFn, int, error) {
	token := tokens[idx]

	if n, err := strconv.Atoi(token); err == nil {
		return forth.Number(n), idx + 1, nil
	}

	if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return forth.Message(token), idx + 1, nil
	}

	if fn, ok := Funcs[token]; ok {
		return fn, idx + 1, nil
	}

	if token == "[" {
		return p.parseLoop(tokens, idx+1)
	}

	return forth.Function(token), idx + 1, nil
}

func (p *Parser) parseLoop(tokens []string, idx int) (forth.StackFn, int, error) {
	l := len(tokens)
	funcs := []forth.StackFn{}
	for idx < l {
		token := tokens[idx]
		if token == "]" {
			return forth.Loop(funcs), idx + 1, nil
		}

		fn, newIdx, err := p.parseAtom(tokens, idx)
		if err != nil {
			return nil, newIdx, err
		}
		funcs = append(funcs, fn)
		idx = newIdx
	}
	return nil, idx, fmt.Errorf("Closing ']' not found")
}

func (p *Parser) parseDefine(tokens []string, idx int) (int, error) {
	l := len(tokens)
	token := tokens[idx]
	if token != "define" {
		return idx, nil
	}

	funcs := []forth.StackFn{}
	idx++
	if idx >= l {
		return idx, fmt.Errorf("Unexpected EOF after 'define'")
	}
	name := tokens[idx]
	idx++
	for idx < l {
		token := tokens[idx]
		if token == ";" {
			forth.WithFunc(name, forth.Sequence(funcs))(p.forth)
			return idx + 1, nil
		}

		fn, newIdx, err := p.parseAtom(tokens, idx)
		if err != nil {
			return newIdx, err
		}
		funcs = append(funcs, fn)
		idx = newIdx
	}
	return idx, fmt.Errorf("Token ';' not found")
}
