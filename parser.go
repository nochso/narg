package narg

import (
	"errors"
	"fmt"
	"io"
)

func Parse(r io.Reader) (Doc, error) {
	p := &Parser{l: NewLexer(r)}
	return p.Parse()
}

type Parser struct {
	l *Lexer
}

// Parse all items.
func (p *Parser) Parse() (Doc, error) {
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
	return i, err
}

func (p *Parser) parseName(i Item) (Item, error) {
	t := p.scanIgnore(TokenWhitespace, TokenLinefeed, TokenComment)
	if t.Type == TokenEOF {
		return i, io.EOF
	}
	if t.Error() != nil {
		return i, t.Error()
	}
	if t.Type == TokenBraceClose {
		return i, errEos
	}
	if t.Type == TokenBraceOpen {
		return i, errSos
	}
	if t.Type != TokenQuotedValue && t.Type != TokenUnquotedValue {
		err := fmt.Errorf("line %d, column %d: expected quoted or unquoted value at beginning of item, got %s %#v", t.Line, t.Col, t.Type, t.Str)
		return i, err
	}
	i.Name = t.String()
	return i, nil
}

func (p *Parser) parseArgs(i Item) (Item, error) {
	for {
		t := p.scanIgnore(TokenWhitespace, TokenComment)
		if t.Type == TokenEOF || t.Type == TokenLinefeed {
			// valid Item end without any (more) args
			return i, nil
		}
		if t.Error() != nil {
			return i, t.Error()
		}
		if t.Type == TokenBraceOpen {
			return i, errSos
		}
		if t.Type == TokenBraceClose {
			return i, errEos
		}
		if t.Type != TokenQuotedValue && t.Type != TokenUnquotedValue {
			err := fmt.Errorf("line %d, column %d: expected quoted or unquoted value as argument no. %d, got %s %#v", t.Line, t.Col, len(i.Args)+1, t.Type, t.Str)
			return i, err
		}
		i.Args = append(i.Args, t.String())
	}
}

func (p *Parser) scanIgnore(ignore ...TokenType) Token {
	for p.l.Scan() {
		t := p.l.Token
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
	}
	return p.l.Token
}
