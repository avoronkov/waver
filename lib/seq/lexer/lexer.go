package lexer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	floatRe       = regexp.MustCompile(`^[0-9]+\.[0-9]+`)
	numberRe      = regexp.MustCompile(`^(0x)?[0-9]+`)
	hexRe         = regexp.MustCompile(`^0[xX][0-9A-Fa-f]+`)
	identRe       = regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*`)
	stringRe      = regexp.MustCompile(`^".*"`)
	singleQuoteRe = regexp.MustCompile(`^'.*'`)
)

type Lexer struct {
	in *bufio.Scanner

	line    string
	index   int
	started bool

	current Token

	literalTokens []literalToken

	bodyNextLine bool

	lineNum int
}

func NewLexer(in io.Reader) *Lexer {
	return &Lexer{
		in: bufio.NewScanner(in),
		literalTokens: []literalToken{
			{"->", ArrowRightToken{}},
			{":", ColonToken{}},
			{"+", PlusToken{}},
			{"-", MinusToken{}},
			{"*", MultiplyToken{}},
			{",", ComaToken{}},
			{"=", AssignToken{}},
			{"{", LCurlyBracket{}},
			{"}", RCurlyBracket{}},
			{"[", LSquareBracket{}},
			{"]", RSquareBracket{}},
			{"%%", DoublePercentToken{}},
			{"%", PercentToken{}},
		},
	}
}

func (l *Lexer) Top() (token Token, err error) {
	if l.current == nil {
		l.current, err = l.nextToken()
		if err != nil {
			return nil, err
		}
	}
	return l.current, nil
}

func (l *Lexer) Pop() (token Token, err error) {
	return l.nextToken()
}

func (l *Lexer) nextToken() (token Token, err error) {
	if l.current != nil {
		token = l.current
		l.current = nil
		return
	}

	// Handle new line
	if l.index >= len(l.line) {
		if l.bodyNextLine {
			l.bodyNextLine = false
			return l.scanBody()
		}
		tok, err := l.scanNextLine()
		if tok != nil || err != nil {
			return tok, err
		}
	}

	// Skip whitespaces
	for l.index < len(l.line) && l.line[l.index] == ' ' {
		l.index++
	}
	if l.index >= len(l.line) {
		// TODO get rid of recursion
		return l.nextToken()
	}

	line := l.line[l.index:]

	mp := DoublePercentToken{}
	for _, lt := range l.literalTokens {
		if strings.HasPrefix(line, lt.literal) {
			l.index += len(lt.literal)
			if lt.token == mp {
				l.bodyNextLine = true
			}
			return lt.token, nil
		}
	}

	if strings.HasPrefix(line, "#") {
		tok := CommentToken(line[1:])
		l.index += len(line)
		return tok, nil
	}

	if f := floatRe.FindString(line); f != "" {
		l.index += len(f)
		fl, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot convert to float '%v': %w", f, err)
		}
		return FloatToken(fl), nil
	}

	if h := hexRe.FindString(line); h != "" {
		l.index += len(h)
		hex, err := strconv.ParseInt(h, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot parse hex '%v': %w", h, err)
		}
		return HexToken(hex), nil

	}

	if num := numberRe.FindString(line); num != "" {
		l.index += len(num)
		i, err := strconv.ParseInt(num, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("Cannot convert to integer '%v': %v", num, err)
		}
		return NumberToken(i), nil
	}

	if ident := identRe.FindString(line); ident != "" {
		l.index += len(ident)
		return IdentToken(ident), nil
	}

	if str := stringRe.FindString(line); str != "" {
		l.index += len(str)
		return StringLiteral(strings.Trim(str, `"`)), nil
	}

	if str := singleQuoteRe.FindString(line); str != "" {
		l.index += len(str)
		return StringLiteral(strings.Trim(str, `'`)), nil
	}

	return nil, fmt.Errorf("Unexpected token near: '%v'", line)
}

func (l *Lexer) scanNextLine() (Token, error) {
	l.lineNum++
	if !l.in.Scan() {
		if err := l.in.Err(); err != nil {
			return nil, err
		}
		return EofToken{}, nil
	}

	l.line = strings.TrimSpace(l.in.Text())
	l.index = 0
	if l.started {
		return EolToken{}, nil
	} else {
		l.started = true
	}
	return nil, nil
}

func (l *Lexer) scanBody() (Token, error) {
	lines := []string{}
	for l.in.Scan() {
		l.lineNum++
		line := l.in.Text()
		if strings.TrimSpace(line) == "%%" {
			l.current = DoublePercentToken{}
			l.line = ""
			l.index = 0
			break
		}
		lines = append(lines, line)
	}
	if err := l.in.Err(); err != nil {
		return nil, err
	}
	return BodyToken(strings.Join(lines, "\n") + "\n"), nil
}

func (l *Lexer) AllTokens() (tokens []Token, err error) {
	for {
		t, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		if _, ok := t.(EofToken); ok {
			break
		}
		tokens = append(tokens, t)
	}
	return
}

func (l *Lexer) LineNum() int {
	return l.lineNum
}

func SplitLine(line string) (tokens []Token, err error) {
	lexer := NewLexer(strings.NewReader(line))
	return lexer.AllTokens()
}

type literalToken struct {
	literal string
	token   Token
}
