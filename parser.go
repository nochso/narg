package narg

import (
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
	return doc, err
}

func (p *Parser) parse() (i Item, err error) {
	i, err = p.parseName(i)
	if err != nil {
		return
	}
	i, err = p.parseArgs(i)
	return
}

func (p *Parser) parseName(i Item) (Item, error) {
	t := p.scanIgnore(TokenWhitespace, TokenLinefeed, TokenComment)
	if t.Type == TokenEOF {
		return i, io.EOF
	}
	if t.Error() != nil {
		return i, t.Error()
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
		if t.Type == TokenBraceOpen || t.Type == TokenBraceClose {
			return i, fmt.Errorf("nested items not yet implemented")
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
