package narg

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Decode narg input r into a given struct v.
// v must be a pointer to the struct you want to decode into.
func Decode(r io.Reader, v interface{}) error {
	items, err := Parse(r)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("decode of non-pointer %s", rv.Type())
	}
	if rv.IsNil() {
		return fmt.Errorf("decode of nil %s", rv.Type())
	}
	return decStruct(rv.Elem(), &Item{Children: items})
}

func decStruct(v reflect.Value, item *Item) error {
	if v.CanAddr() && v.Addr().CanInterface() {
		if u, ok := v.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(item.Args[0]))
		}
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		key := strings.ToLower(v.Type().Field(i).Name)
		tag, ok := v.Type().Field(i).Tag.Lookup("narg")
		if ok {
			// append positional arg to named children
			pos, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				return err
			}
			if pos < int64(len(item.Args)) {
				if len(item.Children.Filter(key)) == 0 {
					item.Children = append(item.Children, Item{Name: key, Args: []string{item.Args[pos]}})
				}
			}
		}
		for _, itm := range item.Children.Filter(key) {
			err := dec(f, &itm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dec(f reflect.Value, item *Item) error {
	if item == nil {
		return nil
	}
	var err error
	switch f.Kind() {
	case reflect.String:
		f.SetString(item.Args[0])
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = decInt(f, item)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = decUint(f, item)
	case reflect.Float32, reflect.Float64:
		err = decFloat(f, item)
	case reflect.Bool:
		return decBool(f, item)
	case reflect.Slice:
		err = decSlice(f, item)
	case reflect.Struct:
		err = decStruct(f, item)
	default:
		err = fmt.Errorf("can not decode unknown type %s", f.Kind())
	}
	return err
}

func decBool(f reflect.Value, item *Item) error {
	b, err := strconv.ParseBool(item.Args[0])
	if err == nil {
		f.SetBool(b)
	}
	return err
}

func decSlice(f reflect.Value, item *Item) error {
	if f.Type().Elem().Kind() == reflect.Struct {
		ptr := reflect.New(f.Type().Elem()).Interface()
		s := reflect.ValueOf(ptr).Elem()
		err := decStruct(s, item)
		f.Set(reflect.Append(f, s))
		return err
	}
	f.Set(reflect.MakeSlice(f.Type(), len(item.Args), len(item.Args)))
	f.SetLen(len(item.Args))
	for n, arg := range item.Args {
		err := dec(f.Index(n), &Item{Args: []string{arg}})
		if err != nil {
			return err
		}
	}
	return nil
}

func decFloat(f reflect.Value, item *Item) error {
	x, err := strconv.ParseFloat(item.Args[0], f.Type().Bits())
	if err == nil {
		f.SetFloat(x)
	}
	return err
}

func decInt(f reflect.Value, item *Item) error {
	i, err := strconv.ParseInt(item.Args[0], 0, f.Type().Bits())
	if err == nil {
		f.SetInt(i)
	}
	return err
}

func decUint(f reflect.Value, item *Item) error {
	i, err := strconv.ParseUint(item.Args[0], 0, f.Type().Bits())
	if err == nil {
		f.SetUint(i)
	}
	return err
}
