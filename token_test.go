package narg

import "testing"

func TestToken_String(t *testing.T) {
	tok := Token{
		Line: 1,
		Col:  1,
		Type: TokenUnquotedValue,
		Str:  "abc",
	}
	act := tok.DebugString()
	exp := `1:1 TokenUnquotedValue "abc"`
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}

	// Invalid/unknown TokenType
	tok.Type = -1
	exp = `1:1 TokenType(-1) "abc"`
	act = tok.DebugString()
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}
}
