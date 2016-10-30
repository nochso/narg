package narg

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/nochso/golden"
)

func TestLexer_read(t *testing.T) {
	t.Parallel()
	l := NewLexer(strings.NewReader("a"))
	r := l.read()
	if r != 'a' {
		cmp(t, 'a', r)
	}
	r = l.read()
	if r != eof {
		cmp(t, eof, r)
	}
	return
}

func TestLexer_unread(t *testing.T) {
	t.Parallel()
	l := NewLexer(strings.NewReader("a"))
	r := l.read()
	if r != 'a' {
		cmp(t, 'a', r)
	}
	err := l.unread()
	if err != nil {
		t.Fatal(err)
	}
	err = l.unread()
	if err == nil {
		cmp(t, nil, err)
	}
}

func init() {
	golden.BasePath = "test-fixtures"
}

var update = flag.Bool("update", false, "update golden test files")

func TestLexer_Scan(t *testing.T) {
	for tc := range golden.Dir(t, "ok") {
		tc.Test(func(c golden.Case) []byte {
			c.T.Parallel()
			r := c.In.Reader()
			defer r.Close()
			l := NewLexer(r)
			act := &bytes.Buffer{}
			for l.Scan() {
				if act.Len() > 0 {
					act.WriteByte('\n')
				}
				act.WriteString(l.Token.String())
			}
			if l.Token.Error() != nil {
				tc.T.Error(l.Token.Error())
			}
			return act.Bytes()
		}, *update)
	}
}

func cmp(t *testing.T, exp, act interface{}) {
	t.Fatalf("expected %#v; got %#v", exp, act)
}
