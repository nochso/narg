package narg

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/nochso/narg/token"
)

// Lexer reads narg input and parses it into tokens.
type Lexer struct {
	// Token that was read after the latest successful call to Scan()
	Token token.T
	Err   error
	str   string
	start int
	pos   int
	line  int
	col   int
}

var eof = rune(0)

// NewLexer returns a new Lexer for parsing tokens.
func NewLexer(s string) *Lexer {
	return &Lexer{
		str:  s,
		line: 1,
		col:  1,
	}
}

// Scan attempts to read the next Token into Lexer.Token.
// Returns true when a new Token is ready.
func (l *Lexer) Scan() bool {
	l.Token = l.scan()
	l.line += strings.Count(l.Token.Str, "\n")
	lastIndex := strings.LastIndexByte(l.Token.Str, '\n')
	if lastIndex == -1 {
		// there's no linefeed, read full token str
		lastIndex = 0
	} else {
		// begin new line
		l.col = 0
	}
	l.col += utf8.RuneCountInString(l.Token.Str[lastIndex:])

	return l.Token.Type != token.EOF && l.Err == nil
}

func (l *Lexer) scan() token.T {
	r := l.read()
	switch r {
	case eof:
		return l.newToken(token.EOF)
	case '{':
		return l.newToken(token.BraceOpen)
	case '}':
		return l.newToken(token.BraceClose)
	case '\n':
		return l.newToken(token.Linefeed)
	case '#':
		return l.scanWhile(token.Comment, continuesComment)
	case '"':
		return l.scanQuotedValue()
	}
	if isWhitespace(r) {
		return l.scanWhile(token.Whitespace, isWhitespace)
	}
	t := l.scanWhile(token.UnquotedValue, isUnquotedValue)
	return l.notFollowedBy(t, "value is missing separator from next value", isQuote)
}

func (l *Lexer) scanWhile(tokenType token.Type, fn func(rune) bool) token.T {
	r := l.read()
	for r != eof {
		if !fn(r) {
			l.unread()
			break
		}
		r = l.read()
	}
	return l.newToken(tokenType)
}

func (l *Lexer) scanQuotedValue() token.T {
	r := l.read()
	escaped := false
	for r != eof {
		if !escaped && r == '"' {
			break
		}
		if !escaped && r == '\\' {
			escaped = true
			r = l.read()
			continue
		}
		r = l.read()
		escaped = false
	}
	if r == eof {
		t := l.newToken(token.Invalid)
		return l.setErr(t, "quoted value is missing closing quote")
	}
	t := l.newToken(token.QuotedValue)
	return l.notFollowedBy(t, "value is missing separator from next value", invalidAfterQuotedValue)
}

func (l *Lexer) setErr(t token.T, err string) token.T {
	l.Err = fmt.Errorf("line %d, column %d: %s: %#v", t.Line, t.Col, err, t.Str)
	t.Type = token.Invalid
	return t
}

func (l *Lexer) notFollowedBy(t token.T, err string, invalidFn func(rune) bool) token.T {
	r := l.read()
	if r == eof {
		return t
	}
	if !invalidFn(r) {
		l.unread()
		return t
	}
	t.Str += string(r)
	return l.setErr(t, err)
}

func (l *Lexer) newToken(tokenType token.Type) token.T {
	tok := token.T{
		Str:  l.str[l.start:l.pos],
		Type: tokenType,
		Line: l.line,
		Col:  l.col,
	}
	l.start = l.pos
	return tok
}

func (l *Lexer) read() rune {
	if len(l.str[l.pos:]) == 0 {
		return eof
	}
	r, s := utf8.DecodeRuneInString(l.str[l.pos:])
	l.pos += s
	return r
}

func (l *Lexer) unread() error {
	if l.pos == 0 {
		return io.EOF
	}
	_, s := utf8.DecodeLastRuneInString(l.str[:l.pos])
	l.pos -= s
	if l.pos < l.start {
		l.pos = l.start
		return io.EOF
	}
	return nil
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\r' || r == '\t'
}

func isQuote(r rune) bool {
	return r == '"'
}

func continuesComment(r rune) bool {
	return r != '\n'
}

func invalidAfterQuotedValue(r rune) bool {
	return isQuote(r) || isUnquotedValue(r)
}

// stop on anything meaningful
func isUnquotedValue(r rune) bool {
	return !isWhitespace(r) && r != '{' && r != '}' && r != '\n' && r != '#' && r != '"'
}
