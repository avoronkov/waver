package parser2

import (
	"fmt"

	"github.com/avoronkov/waver/lib/seq/lexer"
)

func (p *Parser) parseAssignmentStatement(lx *lexer.Lexer) error {
	first, err := lx.Pop()
	if err != nil {
		return err
	}
	name, ok := first.(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Unexpected first token in assignent statement: %v (%T)", first, first)
	}

	second, err := lx.Pop()
	if err != nil {
		return err
	}

	if _, ok := second.(lexer.AssignToken); ok {
		return p.parseAssignVar(lx, string(name))
	}

	if _, ok := second.(lexer.DoubleAssignToken); ok {
		return p.parseAssignSignaler(lx, string(name))
	}

	param, ok := second.(lexer.IdentToken)
	if !ok {
		return fmt.Errorf("Unexpected second token in assignment statement: %v (%T)", second, second)
	}

	third, err := lx.Pop()
	if err != nil {
		return err
	}
	if _, ok := third.(lexer.AssignToken); !ok {
		return fmt.Errorf("Unexpected third token in assignment statement: %v (%T)", third, third)
	}

	return p.parseUdf(lx, string(name), string(param))
}

func (p *Parser) parseAssignVar(lx *lexer.Lexer, name string) error {
	// parse atom
	atom, err := p.parseAtom(lx)
	if err != nil {
		return err
	}
	endl, err := lx.Pop()
	if err != nil {
		return err
	}
	switch endl.(type) {
	case lexer.EolToken, lexer.EofToken:
	default:
		return fmt.Errorf("Expected EOL at the end of assign statement, found: %v (%T)", endl, endl)
	}
	p.seq.Assign(name, atom)
	return nil
}

func (p *Parser) parseUdf(lx *lexer.Lexer, name string, param string) error {
	// parse atom
	atom, err := p.parseAtom(lx)
	if err != nil {
		return err
	}
	endl, err := lx.Pop()
	if err != nil {
		return err
	}
	switch endl.(type) {
	case lexer.EolToken, lexer.EofToken:
	default:
		return fmt.Errorf("Expected EOL at the end of assign statement, found: %v (%T)", endl, endl)
	}
	p.userFunctions[name] = UserFunction{
		Name: name,
		Arg:  param,
		Fn:   atom,
	}
	return nil
}

func (p *Parser) parseAssignSignaler(lx *lexer.Lexer, name string) error {
	sigs, err := p.parseSignaler(lx)
	if err != nil {
		return err
	}
	p.userSignalers[name] = sigs
	return nil
}
