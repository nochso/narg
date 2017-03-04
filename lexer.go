package narg

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode/utf8"
)

// Lexer reads narg input and parses it into tokens.
type Lexer struct {
	// Token that was read after the latest successful call to Scan()
	Token Token
	buf   *bufio.Reader
	line  int
	col   int
}

var eof = rune(0)

// NewLexer returns a new Lexer for parsing tokens.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		buf:  bufio.NewReader(r),
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

	return l.Token.Type != TokenEOF && !l.Token.Type.IsError()
}

func (l *Lexer) scan() Token {
	r := l.read()
	switch r {
	case eof:
		return l.newToken(string(r), TokenEOF)
	case '{':
		return l.newToken(string(r), TokenBraceOpen)
	case '}':
		return l.newToken(string(r), TokenBraceClose)
	case '\n':
		return l.newToken(string(r), TokenLinefeed)
	case '#':
		return l.scanWhile(r, TokenComment, continuesComment)
	case '"':
		return l.scanQuotedValue()
	}
	if isWhitespace(r) {
		return l.scanWhile(r, TokenWhitespace, isWhitespace)
	}
	t := l.scanWhile(r, TokenUnquotedValue, isUnquotedValue)
	return l.notFollowedBy(t, TokenInvalidValueMissingSeparator, isQuote)
}

func (l *Lexer) scanWhile(start rune, tokenType TokenType, fn func(rune) bool) Token {
	buf := &bytes.Buffer{}
	buf.WriteRune(start)
	r := l.read()
	for r != eof {
		if !fn(r) {
			l.unread()
			break
		}
		buf.WriteRune(r)
		r = l.read()
	}
	return l.newToken(buf.String(), tokenType)
}

func (l *Lexer) scanQuotedValue() Token {
	buf := bytes.NewBufferString(`"`)
	r := l.read()
	escaped := false
	for r != eof {
		buf.WriteRune(r)
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
		return l.newToken(buf.String(), TokenInvalidValueMissingClosingQuote)
	}
	t := l.newToken(buf.String(), TokenQuotedValue)
	return l.notFollowedBy(t, TokenInvalidValueMissingSeparator, invalidAfterQuotedValue)
}

func (l *Lexer) notFollowedBy(t Token, invalidType TokenType, invalidFn func(rune) bool) Token {
	r := l.read()
	if r == eof {
		return t
	}
	if !invalidFn(r) {
		l.unread()
		return t
	}
	t.Str += string(r)
	t.Type = invalidType
	return t
}

func (l *Lexer) newToken(str string, tokenType TokenType) Token {
	return Token{
		Str:  str,
		Type: tokenType,
		Line: l.line,
		Col:  l.col,
	}
}

func (l *Lexer) read() rune {
	r, _, err := l.buf.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (l *Lexer) unread() error {
	return l.buf.UnreadRune()
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
