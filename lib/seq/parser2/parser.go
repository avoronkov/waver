package parser2

import (
	"fmt"
	"io"

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
	userFunctions map[string]parser.UserFunction

	scale notes.Scale
}

func (p *Parser) parseReader(reader io.Reader) error {
	lx := lexer.NewLexer(reader)
	_ = lx

	for {
		// Detect type of next statement
		token, err := lx.Top()
		if err != nil {
			return err
		}
		if _, ok := token.(lexer.EofToken); ok {
			break
		}
		if _, ok := p.modParsers[token]; ok {
			if err := p.parseSignalStatement(lx); err != nil {
				return err
			}
		} else if _, ok := token.(lexer.IdentToken); ok {
			if err := p.parseAssignmentStatement(lx); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Cannot detect type of statement: %v", token)
		}
	}
	return nil
}
