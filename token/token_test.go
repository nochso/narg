package token

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

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

var quoteTests = []struct {
	in  string
	exp string
}{
	{"", `""`},
	{"abc", "abc"},
	{"a b c", `"a b c"`},
	{`"some wise quote"`, `"\"some wise quote\""`},
}

func TestQuote(t *testing.T) {
	for _, qt := range quoteTests {
		act := Quote(qt.in)
		if act != qt.exp {
			t.Error(pretty.Compare(act, qt.exp))
		}
	}
}
