package token

import "testing"

func TestToken_String(t *testing.T) {
	tok := T{
		Line: 1,
		Col:  1,
		Type: UnquotedValue,
		Str:  "abc",
	}
	act := tok.DebugString()
	exp := `1:1 UnquotedValue "abc"`
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}

	// Invalid/unknown Type
	tok.Type = -1
	exp = `1:1 Type(-1) "abc"`
	act = tok.DebugString()
	if exp != act {
		t.Fatalf("expected %#v; got %#v", exp, act)
	}
}
