package lexer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestAllTokens(t *testing.T) {
	input := `:4 +2 -> {kick}`
	lx := NewLexer(strings.NewReader(input))

	act, err := lx.AllTokens()
	if err != nil {
		t.Fatalf("Lexer.AllTokens() failed: %v", err)
	}

	exp := []Token{
		ColonToken{},
		NumberToken{4},
		PlusToken{},
		NumberToken{2},
		ArrowRightToken{},
		LCurlyBracket{},
		IdentToken{"kick"},
		RCurlyBracket{},
	}

	compareTokenSlices(t, act, exp)
}

func TestMultiline(t *testing.T) {
	input := `a=4

c->d
`
	lx := NewLexer(strings.NewReader(input))

	act, err := lx.AllTokens()
	if err != nil {
		t.Fatalf("Lexer.AllTokens() failed: %v", err)
	}

	exp := []Token{
		IdentToken{"a"},
		AssignToken{},
		NumberToken{4},
		EolToken{},
		EolToken{},
		IdentToken{"c"},
		ArrowRightToken{},
		IdentToken{"d"},
	}

	compareTokenSlices(t, act, exp)
}

func TestPragma(t *testing.T) {
	is := is.New(t)
	input := `%tempo 120
y = 1`

	lx := NewLexer(strings.NewReader(input))
	act, err := lx.AllTokens()
	is.NoErr(err) // Lexer.AllTokens() failed

	exp := []Token{
		Percent{},
		IdentToken{"tempo"},
		NumberToken{120},
		EolToken{},
		IdentToken{"y"},
		AssignToken{},
		NumberToken{1},
	}
	compareTokenSlices(t, act, exp)
}

func TestStringLiteral(t *testing.T) {
	is := is.New(t)
	input := `%sample kick "2/kick"`

	lx := NewLexer(strings.NewReader(input))
	act, err := lx.AllTokens()
	is.NoErr(err) // Lexer.AllTokens() failed

	exp := []Token{
		Percent{},
		IdentToken{"sample"},
		IdentToken{"kick"},
		StringLiteral("2/kick"),
	}
	compareTokenSlices(t, act, exp)
}

func TestMultilinePragma(t *testing.T) {
	is := is.New(t)
	input := `%%wave foo bar
- one:
    two: three
%%
foo = bar
`

	lx := NewLexer(strings.NewReader(input))
	act, err := lx.AllTokens()
	is.NoErr(err) // LexerAllTokens() failed

	exp := []Token{
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
	}
	compareTokenSlices(t, act, exp)
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
