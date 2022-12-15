package parser2

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/avoronkov/waver/etc/std"
	"github.com/avoronkov/waver/lib/notes"
	"github.com/avoronkov/waver/lib/seq/lexer"
	"github.com/avoronkov/waver/lib/seq/parser"
)

type Parser struct {
	file string

	seq          parser.Seq
	tempoSetters []parser.TempoSetter
	instSet      parser.InstrumentSet

	modParsers    map[lexer.Token]ModParser
	funcParsers   map[string]FunctionParser
	userFunctions map[string]parser.UserFunction
	pragmaParsers map[string]pragmaParser

	scale notes.Scale

	globalCtx map[string]interface{}

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
		if _, ok := token.(lexer.EofToken); ok {
			break
		}
		if _, ok := p.modParsers[token]; ok {
			err = p.parseSignalStatement(lx)
		} else if _, ok := token.(lexer.IdentToken); ok {
			err = p.parseAssignmentStatement(lx)
		} else if _, ok := token.(lexer.Percent); ok {
			err = p.parsePragma(lx)
		} else if _, ok := token.(lexer.DoublePercent); ok {
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
