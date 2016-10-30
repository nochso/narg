//go:generate stringer -type TokenType

package narg

import "fmt"

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLinefeed
	TokenWhitespace
	TokenComment
	TokenBraceOpen
	TokenBraceClose
	TokenQuotedValue
	TokenUnquotedValue
	TokenInvalidValueMissingClosingQuote
	TokenInvalidValueMissingSeparator
)

func (tt TokenType) IsError() bool {
	return tt == TokenInvalidValueMissingClosingQuote ||
		tt == TokenInvalidValueMissingSeparator
}

func (tt TokenType) Error(t Token) error {
	if !tt.IsError() {
		return nil
	}
	msg := "unsupported token"
	switch tt {
	case TokenInvalidValueMissingClosingQuote:
		msg = "quoted value is missing closing quote"
	case TokenInvalidValueMissingSeparator:
		msg = "value is missing separator from next value"
	}
	return fmt.Errorf("error on line %d, column %d: %s: %#v", t.Line, t.Col, msg, t.Str)
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

func (t Token) DebugString() string {
	return fmt.Sprintf("%d:%d %s %#v", t.Line, t.Col, t.Type, t.Str)
}
