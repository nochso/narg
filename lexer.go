package cml

import (
	"bufio"
	"io"
)

type Lexer struct {
	buf  *bufio.Reader
	line int
	col  int
}

var eof = rune(0)

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		buf:  bufio.NewReader(r),
		line: 1,
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
