//go:generate stringer -type Type

package token

import (
	"bytes"
	"fmt"
	"strings"
)

type Type int

const (
	EOF Type = iota
	Linefeed
	Whitespace
	Comment
	BraceOpen
	BraceClose
	QuotedValue
	UnquotedValue
	InvalidValueMissingClosingQuote
	InvalidValueMissingSeparator
)

func (tt Type) IsError() bool {
	return tt == InvalidValueMissingClosingQuote ||
		tt == InvalidValueMissingSeparator
}

func (tt Type) Error(t Token) error {
	if !tt.IsError() {
		return nil
	}
	msg := "unsupported token"
	switch tt {
	case InvalidValueMissingClosingQuote:
		msg = "quoted value is missing closing quote"
	case InvalidValueMissingSeparator:
		msg = "value is missing separator from next value"
	}
	return fmt.Errorf("line %d, column %d: %s: %#v", t.Line, t.Col, msg, t.Str)
}

type Token struct {
	// Raw string including double quotes.
	// Use String() for a cleaned up version.
	Str  string
	Type Type
	Line int
	Col  int
}

func (t Token) Error() error {
	return t.Type.Error(t)
}

func (t Token) String() string {
	return Unquote(t.Str)
}

// Quote a string if needed.
//
// If needed s is surrounded with double quotes.
// Backslashes \ and double quotes " will be escaped with backslashes.
func Quote(s string) string {
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
func Unquote(s string) string {
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
