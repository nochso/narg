package cml

import (
	"strings"
	"testing"
)

func TestLexer_read(t *testing.T) {
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

func cmp(t *testing.T, exp, act interface{}) {
	t.Fatalf("expected %#v; got %#v", exp, act)
}
