//go:generate stringer -type TokenType

package cml

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
