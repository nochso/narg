package narg

import (
	"errors"
	"fmt"
	"io"

	"io/ioutil"

	"github.com/nochso/narg/token"
)

// Parse narg input into an ItemSlice.
func Parse(r io.Reader) (ItemSlice, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	p := &Parser{l: NewLexer(string(b))}
	return p.Parse()
}

// Parser turns tokens into items.
type Parser struct {
	l   *Lexer
	buf []token.T
}

// Parse all items.
func (p *Parser) Parse() (ItemSlice, error) {
	doc := []Item{}
	item, err := p.parse()
	for err == nil {
		doc = append(doc, item)
		item, err = p.parse()
	}
	if err == io.EOF {
		err = nil
	}
	if err == errEos || err == errSos {
		t := p.l.Token
		err = fmt.Errorf("line %d, column %d: expected quoted or unquoted value at beginning of item, got %s %#v", t.Line, t.Col, t.Type, t.Str)
	}
	return doc, err
}

// parse exactly one Item
func (p *Parser) parse() (i Item, err error) {
	i, err = p.parseName(i)
	if err != nil {
		return
	}
	i, err = p.parseArgs(i)
	if err == errSos {
		i, err = p.parseChildren(i)
		if err == errEos { // ok, it ends with a '}'
			err = nil
		}
		return
	}
	return
}

var errEos = errors.New("end of child section")
var errSos = errors.New("start of child section")

func (p *Parser) parseChildren(i Item) (Item, error) {
	child, err := p.parse()
	for err == nil {
		i.Children = append(i.Children, child)
		child, err = p.parse()
	}
	if err == errEos && child.Name != "" {
		i.Children = append(i.Children, child)
	}
	return i, err
}

func (p *Parser) parseName(i Item) (Item, error) {
	t := p.scanIgnore(token.Whitespace, token.Linefeed, token.Comment)
	if t.Type == token.EOF {
		return i, io.EOF
	}
	if p.l.Err != nil {
		return i, p.l.Err
	}
	if t.Type == token.BraceClose {
		return i, errEos
	}
	if t.Type == token.BraceOpen {
		return i, errSos
	}
	if t.Type != token.QuotedValue && t.Type != token.UnquotedValue {
		err := fmt.Errorf("line %d, column %d: expected quoted or unquoted value at beginning of item, got %s %#v", t.Line, t.Col, t.Type, t.Str)
		return i, err
	}
	i.Name = t.String()
	return i, nil
}

func (p *Parser) parseArgs(i Item) (Item, error) {
	for {
		t := p.scanIgnore(token.Whitespace, token.Comment)
		if t.Type == token.EOF {
			// valid Item end without any (more) args
			return i, nil
		}
		if t.Type == token.Linefeed {
			// valid end of args, but look ahead
			t = p.scanIgnore(token.Whitespace, token.Comment, token.Linefeed)
			if t.Type != token.BraceClose && t.Type != token.BraceOpen {
				// nah, we can't use this. put it back.
				p.unscan(t)
				return i, nil
			}
			// fall through to open/close brace
		}
		if p.l.Err != nil {
			return i, p.l.Err
		}
		if t.Type == token.BraceOpen {
			return i, errSos
		}
		if t.Type == token.BraceClose {
			return i, errEos
		}
		if t.Type != token.QuotedValue && t.Type != token.UnquotedValue {
			err := fmt.Errorf("line %d, column %d: expected quoted or unquoted value as argument no. %d, got %s %#v", t.Line, t.Col, len(i.Args)+1, t.Type, t.Str)
			return i, err
		}
		i.Args = append(i.Args, t.String())
	}
}

func (p *Parser) scan() (t token.T) {
	if len(p.buf) > 0 {
		t, p.buf = p.buf[len(p.buf)-1], p.buf[:len(p.buf)-1]
		return
	}
	p.l.Scan()
	return p.l.Token
}

func (p *Parser) unscan(t token.T) {
	p.buf = append(p.buf, t)
}

func (p *Parser) scanIgnore(ignore ...token.Type) token.T {
	t := p.scan()
	for t.Type != token.EOF && t.Type != token.Invalid {
		ignored := false
		for _, it := range ignore {
			if t.Type == it {
				ignored = true
				continue
			}
		}
		if !ignored {
			return t
		}
		t = p.scan()
	}
	return t
}
