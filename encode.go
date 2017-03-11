package narg

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Encode the given value by writing its narg representation to w.
func Encode(w io.Writer, v interface{}) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	itm, err := encodeStruct("", rv)
	if err != nil {
		return err
	}
	for i, item := range itm.Children {
		if i > 0 {
			w.Write([]byte{'\n'})
		}
		item.writeString(w, 0)
	}
	return nil
}

func encodeStruct(name string, v reflect.Value) (Item, error) {
	itm := Item{}
	itm.Name = name
	var err error
	for i := 0; i < v.NumField(); i++ {
		ft := v.Type().Field(i)
		var child Item
		if ft.Type.Kind() == reflect.Struct {
			child, err = encodeStruct(strings.ToLower(ft.Name), v.Field(i))
			if err != nil {
				return itm, err
			}
		} else if ft.Type.Kind() == reflect.Slice {
			if ft.Type.Elem().Kind() == reflect.Struct {
				for x := 0; x < v.Field(i).Len(); x++ {
					foo, err := encodeStruct(strings.ToLower(ft.Name), v.Field(i).Index(x))
					if err != nil {
						return itm, err
					}
					itm.Children = append(itm.Children, foo)
				}
				continue
			} else {
				for x := 0; x < v.Field(i).Len(); x++ {
					child.Args = append(child.Args, fmt.Sprintf("%v", v.Field(i).Index(x)))
				}
			}
		} else {
			child.Args = append(child.Args, fmt.Sprintf("%v", v.Field(i)))
		}
		child.Name = strings.ToLower(ft.Name)
		itm.Children = append(itm.Children, child)
	}
	return itm, err
}
