package jsonkit

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

//Jsonkit a simple json kit
type Jsonkit struct {
}

//Encode ecode a object to bytes
func Encode(v interface{}) ([]byte, error) {
	kit := &Jsonkit{}
	var buf bytes.Buffer
	if err := kit.encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *Jsonkit) encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		j.encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := j.encode(buf, v.Index(i)); err != nil {
				return err
			}
		}

		buf.WriteByte(')')

	case reflect.Struct:
		buf.WriteByte('(')
		nField := v.NumField()
		for i := 0; i < nField; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := j.encode(buf, v.Field(i)); err != nil {
				return err
			}

			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := j.encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := j.encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}

	return nil
}

func Decode(data []byte, out interface{}) (err error) {
	deco := &decoder{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	deco.scan.Init(bytes.NewReader(data))
	deco.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", deco.scan.Position, x)
			fmt.Printf("error at %s: %v", deco.scan.Position, x)
		}
	}()

	deco.read(reflect.ValueOf(out).Elem())
	return nil
}

type decoder struct {
	scan  scanner.Scanner
	token rune
}

func (j *decoder) next() {
	j.token = j.scan.Scan()
}

func (j *decoder) text() string {
	return j.scan.TokenText()
}

func (j *decoder) consume(want rune) {
	if j.token != want {
		panic(fmt.Sprintf("got %q, want %q", j.text(), want))
	}
	j.next()
}

func (j *decoder) read(v reflect.Value) {
	if !v.CanSet() {
		panic(fmt.Sprintf("v %v cannot assign", v))
	}

	switch j.token {
	case scanner.Ident:
		if j.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			j.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(j.text())
		v.SetString(s)
		j.next()
		return
	case scanner.Int:
		if v.Kind() == reflect.String {
			v.SetString(j.text())
		} else {
			i, _ := strconv.Atoi(j.text())
			v.SetInt(int64(i))
		}

		j.next()
		return
	case '(':
		j.next()
		j.readList(v)
		j.next()
		return
	}
	panic(fmt.Sprintf("unsupport token %q", j.text()))
}

func (j *decoder) endList() bool {
	switch j.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func (j *decoder) readList(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		for i := 0; j.endList(); i++ {
			j.read(v.Index(i))
		}

	case reflect.Slice:
		for !j.endList() {
			item := reflect.New(v.Type().Elem()).Elem()
			j.read(item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct:
		for !j.endList() {
			j.consume('(')
			if j.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q want field name", j.text()))
			}
			name := j.text()
			j.next()
			val := v.FieldByName(name)
			j.read(val)
			j.consume(')')
		}

	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for !j.endList() {
			j.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			j.read(key)
			value := reflect.New(v.Type().Elem()).Elem()
			j.read(value)
			v.SetMapIndex(key, value)
			j.consume(')')
		}
	default:
		panic(fmt.Sprintf("cannot decode list info %v", v.Type()))
	}
}
