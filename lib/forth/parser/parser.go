package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/avoronkov/waver/lib/forth"
)

func ParseFile(file string) (*forth.Forth, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(r io.Reader) (*forth.Forth, error) {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)
	tokens := []string{}
	for sc.Scan() {
		line := sc.Text()
		tokens = append(tokens, strings.Fields(line)...)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return parseTokens(tokens)
}

var funcs = map[string]forth.StackFn{
	"+":    forth.Plus,
	"-":    forth.Minus,
	"dup":  forth.Dup,
	"drop": forth.Drop,
	"top":  forth.ShowTop,
}

func parseTokens(tokens []string) (*forth.Forth, error) {
	program := []forth.StackFn{}
	idx := 0
	l := len(tokens)
	for idx < l {
		fn, newIdx, err := parseAtom(tokens, idx)
		if err != nil {
			return nil, err
		}
		program = append(program, fn)
		idx = newIdx
	}

	return forth.New(
		forth.WithProgram(program),
	), nil
}

func parseAtom(tokens []string, idx int) (forth.StackFn, int, error) {
	token := tokens[idx]

	if n, err := strconv.Atoi(token); err == nil {
		return forth.Number(n), idx + 1, nil
	}

	if fn, ok := funcs[token]; ok {
		return fn, idx + 1, nil
	}

	if token == "[" {
		return parseLoop(tokens, idx+1)
	}

	return forth.Function(token), idx + 1, nil
}

func parseLoop(tokens []string, idx int) (forth.StackFn, int, error) {
	l := len(tokens)
	funcs := []forth.StackFn{}
	for idx < l {
		token := tokens[idx]
		if token == "]" {
			return forth.Loop(funcs), idx + 1, nil
		}

		fn, newIdx, err := parseAtom(tokens, idx)
		if err != nil {
			return nil, newIdx, err
		}
		funcs = append(funcs, fn)
		idx = newIdx
	}
	return nil, idx, fmt.Errorf("Closing ']' not found")
}
