//go:generate stringer -type TokenType

package narg

import (
	"bytes"
	"fmt"
	"strings"
)

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
	// Raw string including double quotes.
	// Use String() for a cleaned up version.
	Str  string
	Type TokenType
	Line int
	Col  int
}

func (t Token) Error() error {
	return t.Type.Error(t)
}

func (t Token) String() string {
	return unquote(t.Str)
}

// Quote a string if needed.
//
// If needed s is surrounded with double quotes.
// Backslashes \ and double quotes " will be escaped with backslashes.
func quote(s string) string {
	if !strings.ContainsAny(s, "#\"{} \t\r\n") {
		return s
	}
	buf := &bytes.Buffer{}
	buf.WriteByte('"')
	for _, r := range s {
		if r == '\\' || r == '"' {
			buf.WriteByte('\\')
		}
		buf.WriteRune(r)
	}
	buf.WriteByte('"')
	return buf.String()
}

// Unquote a string if needed.
//
// If s is surrounded by double quotes, escaped backslashes and double quotes
// will be unescaped. The surrounding quotes will be removed.
func unquote(s string) string {
	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		return s
	}
	buf := &bytes.Buffer{}
	i := 0
	escaped := false
	for _, r := range s {
		i++
		if i == 1 {
			continue
		}
		if !escaped && r == '\\' {
			escaped = true
			continue
		}
		if !escaped && r == '"' {
			break
		}
		buf.WriteRune(r)
		escaped = false
	}
	return buf.String()
}

func (t Token) DebugString() string {
	return fmt.Sprintf("%d:%d %s %#v", t.Line, t.Col, t.Type, t.String())
}
