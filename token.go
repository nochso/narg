//go:generate stringer -type TokenType

package cml

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLinefeed
	TokenWhitespace
	TokenComment
	TokenBraceOpen
	TokenBraceClose
	TokenValue
)

type Token struct {
	Str  string
	Type TokenType
	Line int
	Col  int
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%d %s %#v", t.Line, t.Col, t.Type, t.Str)
}
