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

// <
type LessToken struct{}

func (LessToken) String() string { return "<" }

// >
type GreaterToken struct{}

func (GreaterToken) String() string { return ">" }

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

// ==
type DoubleAssignToken struct{}

func (DoubleAssignToken) String() string { return "==" }

// & (random)
type AmpersandToken struct{}

func (AmpersandToken) String() string { return "&" }

// @ (sequence)
type AtToken struct{}

func (AtToken) String() string { return "@" }

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

// |
type VerticalBar struct{}

func (VerticalBar) String() string { return "|" }

// ..
type DoubleDot struct{}

func (DoubleDot) String() string { return ".." }

// [0-9]+, 0x[0-9]+
type NumberToken int64

func (n NumberToken) String() string {
	return fmt.Sprintf("%v", int64(n))
}

// 0x[0-9A-Fa-f]+
type HexToken int64

func (t HexToken) String() string {
	return fmt.Sprintf("%#x", int64(t))
}

//[0-9]+\.[0-9]+
type FloatToken float64

func (t FloatToken) String() string {
	return fmt.Sprintf("%v", float64(t))
}

// EOL
type EolToken struct{}

func (e EolToken) String() string { return "<EOL>" }

// EOF
type EofToken struct{}

func (e EofToken) String() string { return "<EOF>" }

// Identifier
type IdentToken string

func (i IdentToken) String() string {
	return fmt.Sprintf("%v", string(i))
}

// %
type PercentToken struct{}

func (PercentToken) String() string { return "%" }

// %%
type DoublePercentToken struct{}

func (DoublePercentToken) String() string { return "%%" }

// Raw pragma body
type BodyToken string

func (b BodyToken) String() string { return string(b) }

// String literal
type StringLiteral string

func (t StringLiteral) String() string { return string(t) }

// Comment
type CommentToken string

func (t CommentToken) String() string { return fmt.Sprintf("/* %v */", string(t)) }
