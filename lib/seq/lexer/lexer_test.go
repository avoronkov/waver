package lexer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   []Token
	}{
		{
			"Common: all tokens",
			`:4 +2 -> {kick}`,
			[]Token{
				ColonToken{},
				NumberToken(4),
				PlusToken{},
				NumberToken(2),
				ArrowRightToken{},
				LCurlyBracket{},
				IdentToken("kick"),
				RCurlyBracket{},
			},
		},
		{
			"Common: multiline",
			`a=4

c->d
`,
			[]Token{
				IdentToken("a"),
				AssignToken{},
				NumberToken(4),
				EolToken{},
				EolToken{},
				IdentToken("c"),
				ArrowRightToken{},
				IdentToken("d"),
			},
		},
		{
			"Pragma: single line",
			`%tempo 120
y = 1`,
			[]Token{
				PercentToken{},
				IdentToken("tempo"),
				NumberToken(120),
				EolToken{},
				IdentToken("y"),
				AssignToken{},
				NumberToken(1),
			},
		},
		{
			"Pragma: multiline",
			`%%wave foo bar
- one:
    two: three
%%
foo = bar
`,
			[]Token{
				DoublePercentToken{},
				IdentToken("wave"),
				IdentToken("foo"),
				IdentToken("bar"),
				BodyToken("- one:\n    two: three\n"),
				DoublePercentToken{},
				EolToken{},
				IdentToken("foo"),
				AssignToken{},
				IdentToken("bar"),
			},
		},
		{
			"Pragma: multiline, empty body",
			`%%wave foo bar
%%
foo = bar
`,
			[]Token{
				DoublePercentToken{},
				IdentToken("wave"),
				IdentToken("foo"),
				IdentToken("bar"),
				BodyToken("\n"),
				DoublePercentToken{},
				EolToken{},
				IdentToken("foo"),
				AssignToken{},
				IdentToken("bar"),
			},
		},
		{
			"String literal",
			`%sample kick "2/kick"`,
			[]Token{
				PercentToken{},
				IdentToken("sample"),
				IdentToken("kick"),
				StringLiteral("2/kick"),
			},
		},
		{
			"Numbers",
			`0 123 23.45 0xF`,
			[]Token{
				NumberToken(0),
				NumberToken(123),
				FloatToken(23.45),
				HexToken(15),
			},
		},
		{
			"Comments: whole line comment",
			`# here is a comment`,
			[]Token{CommentToken(" here is a comment")},
		},
		{
			"Comments: part line comment",
			`x = 11 # assignment`,
			[]Token{
				IdentToken("x"),
				AssignToken{},
				NumberToken(11),
				CommentToken(" assignment"),
			},
		},
		{
			"Comments: multiple lines with comment",
			`# this is an assignment
x = 11`,
			[]Token{
				CommentToken(" this is an assignment"),
				EolToken{},
				IdentToken("x"),
				AssignToken{},
				NumberToken(11),
			},
		},
		{
			"Identifiers: eucl'",
			`eucl' 3 8`,
			[]Token{
				IdentToken("eucl'"),
				NumberToken(3),
				NumberToken(8),
			},
		},
		{
			"Range: 0..31",
			`0..31`,
			[]Token{
				NumberToken(0),
				DoubleDot{},
				NumberToken(31),
			},
		},
		{
			"Single line code",
			"sine | `v * t` |",
			[]Token{
				IdentToken("sine"),
				VerticalBar{},
				CodeLiteral("v * t"),
				VerticalBar{},
			},
		},
		{
			"Multiline code",
			"foo = `let x=0;\nx * v` baz\nbar",
			[]Token{
				IdentToken("foo"),
				AssignToken{},
				CodeLiteral("let x=0;\nx * v"),
				IdentToken("baz"),
				EolToken{},
				IdentToken("bar"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			lx := NewLexer(strings.NewReader(test.input))
			act, err := lx.AllTokens()
			is.NoErr(err)

			compareTokenSlices(t, act, test.exp)
		})
	}
}

func compareTokenSlices(t *testing.T, act, exp []Token) {
	if len(act) != len(exp) {
		t.Errorf("Token slices lengths differ:\nexpected %v (%v)\nactual   %v (%v)", exp, len(exp), act, len(act))
	}
	for i, a := range act {
		if i >= len(exp) {
			break
		}
		e := exp[i]
		if !reflect.DeepEqual(a, e) {
			t.Errorf("Token slices differ at position %v:\nexpected %v\nactual   %v\nexp %v (%T), act %v (%T)", i, exp, act, e, e, a, a)
		}
	}
}
