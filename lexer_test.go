package narg

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/nochso/golden"
)

func init() {
	golden.BasePath = "test-fixtures"
}

var update = flag.Bool("update", false, "update golden test files")

var benchTest = `unquoted "quoted value" {
	   # comment after whitespace
	   sub { "answer" 42 }
}`

func BenchmarkLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(benchTest)
		l := NewLexer(r)
		for l.Scan() {
		}
	}
}

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

func TestLexer_Scan(t *testing.T) {
	t.Parallel()
	tester := func(c golden.Case) []byte {
		r := c.In.Reader()
		defer r.Close()
		l := NewLexer(r)
		act := &bytes.Buffer{}
		for l.Scan() {
			if act.Len() > 0 {
				act.WriteByte('\n')
			}
			act.WriteString(l.Token.DebugString())
		}
		if l.Err != nil {
			c.T.Error(l.Err)
		}
		return act.Bytes()
	}
	golden.TestDir(t, "lexer/ok", func(tc golden.Case) {
		tc.T.Parallel()
		if *update {
			tc.Out.Update(tester(tc))
		}
		tc.Diff(string(tester(tc)))
	})
}

func TestLexer_Scan_error(t *testing.T) {
	t.Parallel()
	tester := func(c golden.Case) []byte {
		r := c.In.Reader()
		defer r.Close()
		l := NewLexer(r)
		act := &bytes.Buffer{}
		for l.Scan() {
			act.WriteString(l.Token.DebugString())
			act.WriteByte('\n')
		}
		if l.Err != nil {
			act.WriteString(l.Err.Error())
		}
		return act.Bytes()
	}
	golden.TestDir(t, "lexer/error", func(tc golden.Case) {
		tc.T.Parallel()
		if *update {
			tc.Out.Update(tester(tc))
		}
		tc.Diff(string(tester(tc)))
	})
}

func cmp(t *testing.T, exp, act interface{}) {
	t.Fatalf("expected %#v; got %#v", exp, act)
}
