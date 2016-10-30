package cml

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	Token Token
	buf   *bufio.Reader
	line  int
	col   int
}

var eof = rune(0)

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		buf:  bufio.NewReader(r),
		line: 1,
		col:  1,
	}
}

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

	return l.Token.Type != TokenEOF
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
	// case '#':
	// 	return l.scanComment()
	// case '"':
	// 	return l.scanQuotedValue()
	}
	// whitespace

	// unquoted value

	return l.newToken(string(r), -1)
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
