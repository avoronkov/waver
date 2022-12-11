package lexer

import "fmt"

type Token interface {
	// Return original token
	String() string
}

// :
type ColonToken struct{}

func (ColonToken) String() string { return ":" }

// +
type PlusToken struct{}

func (PlusToken) String() string { return "+" }

// ->

type ArrowRightToken struct{}

func (ArrowRightToken) String() string { return "->" }

// -
type MinusToken struct{}

func (MinusToken) String() string { return "-" }

// *
type MultiplyToken struct{}

func (MultiplyToken) String() string { return "*" }

// ,
type ComaToken struct{}

func (ComaToken) String() string { return "," }

// =
type AssignToken struct{}

func (AssignToken) String() string { return "=" }

// {
type LCurlyBracket struct{}

func (LCurlyBracket) String() string { return "{" }

// }
type RCurlyBracket struct{}

func (RCurlyBracket) String() string { return "}" }

// [
type LSquareBracket struct{}

func (LSquareBracket) String() string { return "[" }

// ]

type RSquareBracket struct{}

func (RSquareBracket) String() string { return "]" }

// [0-9]+, 0x[0-9]+
type NumberToken struct {
	Num int64
}

func (n NumberToken) String() string {
	return fmt.Sprintf("%v", n.Num)
}

// EOL
type EolToken struct{}

func (e EolToken) String() string { return "<EOL>" }

// EOF
type EofToken struct{}

func (e EofToken) String() string { return "<EOF>" }

// Identifier
type IdentToken struct {
	Value string
}

func (i IdentToken) String() string {
	return fmt.Sprintf("%v", i.Value)
}

// %
type Percent struct{}

func (Percent) String() string { return "%" }

// %%
type DoublePercent struct{}

func (DoublePercent) String() string { return "%%" }

// Raw pragma body
type BodyToken string

func (b BodyToken) String() string { return string(b) }
