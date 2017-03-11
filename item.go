package narg

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/nochso/narg/token"
)

// ItemSlice is a list of items.
type ItemSlice []Item

// Item represents a single block consisting of at least the name.
//
// narg input showing what part of an Item it's parsed into:
//
//	Name Args[0] Args[1] {
//		Children[0].Name Children[0].Args[0]
//	}
type Item struct {
	Name     string
	Args     []string
	Children ItemSlice
}

// String returns an indented string representation.
func (s ItemSlice) String() string {
	buf := &bytes.Buffer{}
	for i, item := range s {
		if i > 0 {
			buf.WriteByte('\n')
		}
		item.writeString(buf, 0)
	}
	return buf.String()
}

// Filter returns all items filtered by name.
// The comparison is case-insensitive.
func (s ItemSlice) Filter(key string) ItemSlice {
	out := ItemSlice{}
	key = strings.ToLower(key)
	for _, itm := range s {
		if strings.ToLower(itm.Name) == key {
			out = append(out, itm)
		}
	}
	return out
}

// String returns the string representation of an Item and its Children.
func (i Item) String() string {
	buf := &bytes.Buffer{}
	i.writeString(buf, 0)
	return buf.String()
}

func (i Item) writeString(w io.Writer, indent int) {
	prefix := strings.Repeat("\t", indent)
	fmt.Fprintf(w, "%s%s", prefix, token.Quote(i.Name))
	for _, arg := range i.Args {
		fmt.Fprintf(w, " %s", token.Quote(arg))
	}
	if len(i.Children) > 0 {
		fmt.Fprint(w, " {\n")
		for _, child := range i.Children {
			child.writeString(w, indent+1)
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "%s}", prefix)
	}
}
