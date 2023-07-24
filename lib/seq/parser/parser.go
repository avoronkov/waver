package parser

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/avoronkov/waver/etc/std"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/types"
)

type Parser struct {
	file string

	seq          Seq
	tempoSetters []TempoSetter
	instSet      InstrumentSet

	ModParsers    map[lexer.Token]ModParser
	FuncParsers   map[lexer.Token]FunctionParser
	PragmaParsers map[string]pragmaParser

	userFunctions map[string]UserFunction
	userSignalers map[string][]types.Signaler

	scale notes.Scale

	globalFilters []map[string]any

	tempo int
}

func (p *Parser) parse() error {
	// WIP parse std lib
	reader := bytes.NewReader(std.StdEdo12)
	if err := p.parseReader(reader); err != nil {
		return err
	}

	// Parse the file itself
	f, err := os.Open(p.file)
	if err != nil {
		return err
	}
	defer f.Close()

	return p.parseReader(f)
}

func (p *Parser) ParseData(data []byte) error {
	// WIP parse std lib
	reader := bytes.NewReader(std.StdEdo12)
	if err := p.parseReader(reader); err != nil {
		return err
	}
	return p.parseReader(bytes.NewReader(data))
}

func (p *Parser) parseReader(reader io.Reader) error {
	lx := lexer.NewLexer(reader)
	_ = lx

	log.Printf("Parsing started.")
	for {
		// Detect type of next statement
		token, err := lx.Top()
		if err != nil {
			return err
		}
		if _, ok := token.(lexer.EolToken); ok {
			_, _ = lx.Pop()
			continue
		}
		if _, ok := token.(lexer.CommentToken); ok {
			_, _ = lx.Pop()
			continue
		}
		if _, ok := token.(lexer.EofToken); ok {
			break
		}
		if _, ok := p.ModParsers[token]; ok {
			err = p.parseSignalStatement(lx)
		} else if _, ok := token.(lexer.HexToken); ok {
			err = p.parseSignalStatement(lx)
		} else if _, ok := token.(lexer.NumberToken); ok {
			err = p.parseSignalStatement(lx)
		} else if _, ok := token.(lexer.IdentToken); ok {
			err = p.parseAssignmentStatement(lx)
		} else if _, ok := token.(lexer.PercentToken); ok {
			err = p.parsePragma(lx)
		} else if _, ok := token.(lexer.DoublePercentToken); ok {
			err = p.parsePragma(lx)
		} else {
			err = fmt.Errorf("Cannot detect type of statement: %v", token)
		}

		if err != nil {
			return err
		}
	}
	log.Printf("Parsing finished.")
	return p.seq.Commit()
}
