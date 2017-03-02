package narg

import (
	"bytes"
	"strings"
)

type ItemSlice []Item

type Item struct {
	Name     string
	Args     []string
	Children ItemSlice
}

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

func (i Item) String() string {
	buf := &bytes.Buffer{}
	i.writeString(buf, 0)
	return buf.String()
}

func (i Item) writeString(buf *bytes.Buffer, indent int) {
	prefix := strings.Repeat("\t", indent)
	buf.WriteString(prefix)
	buf.WriteString(Quote(i.Name))
	for _, arg := range i.Args {
		buf.WriteByte(' ')
		buf.WriteString(Quote(arg))
	}
	if len(i.Children) > 0 {
		buf.WriteString(" {\n")
		for _, child := range i.Children {
			child.writeString(buf, indent+1)
			buf.WriteByte('\n')
		}
		buf.WriteString(prefix)
		buf.WriteString("}")
	}
}
