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
	TokenInvalidValueMissingClosingQuote
	TokenInvalidValueMissingSeparator
)

func (tt TokenType) IsError() bool {
	return tt == TokenInvalidValueMissingClosingQuote ||
		tt == TokenInvalidValueMissingSeparator
}

func (tt TokenType) Error(t Token) error {
	if tt == TokenInvalidValueMissingClosingQuote {
		return fmt.Errorf(
			"error on line %d, column %d: quoted value is missing closing quote: %#v",
			t.Line,
			t.Col,
			t.Str,
		)
	}
	if tt == TokenInvalidValueMissingSeparator {
		return fmt.Errorf(
			"error on line %d, column %d: value is missing separator from next value: %#v",
			t.Line,
			t.Col,
			t.Str,
		)
	}
	return nil
}

type Token struct {
	Str  string
	Type TokenType
	Line int
	Col  int
}

func (t Token) Error() error {
	return t.Type.Error(t)
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%d %s %#v", t.Line, t.Col, t.Type, t.Str)
}
