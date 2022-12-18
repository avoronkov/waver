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
				NumberToken{4},
				PlusToken{},
				NumberToken{2},
				ArrowRightToken{},
				LCurlyBracket{},
				IdentToken{"kick"},
				RCurlyBracket{},
			},
		},
		{
			"Common: multiline",
			`a=4

c->d
`,
			[]Token{
				IdentToken{"a"},
				AssignToken{},
				NumberToken{4},
				EolToken{},
				EolToken{},
				IdentToken{"c"},
				ArrowRightToken{},
				IdentToken{"d"},
			},
		},
		{
			"Pragma: single line",
			`%tempo 120
y = 1`,
			[]Token{
				Percent{},
				IdentToken{"tempo"},
				NumberToken{120},
				EolToken{},
				IdentToken{"y"},
				AssignToken{},
				NumberToken{1},
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
				DoublePercent{},
				IdentToken{"wave"},
				IdentToken{"foo"},
				IdentToken{"bar"},
				BodyToken("- one:\n    two: three\n"),
				DoublePercent{},
				EolToken{},
				IdentToken{"foo"},
				AssignToken{},
				IdentToken{"bar"},
			},
		},
		{
			"String literal",
			`%sample kick "2/kick"`,
			[]Token{
				Percent{},
				IdentToken{"sample"},
				IdentToken{"kick"},
				StringLiteral("2/kick"),
			},
		},
		{
			"Numbers",
			`0 123 23.45 0xF`,
			[]Token{
				NumberToken{0},
				NumberToken{123},
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
				IdentToken{"x"},
				AssignToken{},
				NumberToken{11},
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
				IdentToken{"x"},
				AssignToken{},
				NumberToken{11},
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
