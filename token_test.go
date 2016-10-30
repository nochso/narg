package cml

import "testing"

func TestToken_String(t *testing.T) {
	tok := Token{
		Line: 1,
		Col:  1,
		Type: TokenValue,
		Str:  "abc",
	}
	act := tok.String()
	exp := `1:1 TokenValue "abc"`
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}

	// Invalid/unknown TokenType
	tok.Type = -1
	exp = `1:1 TokenType(-1) "abc"`
	act = tok.String()
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}
}
