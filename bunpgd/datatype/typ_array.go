package datatype

import (
	"reflect"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

func appendArray(sf schema.AppenderFunc) func(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
	return func(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
		length := v.Len()
		if length == 0 {
			return append(b, "'{}'"...)
		}

		for i := 0; i < length; i++ {
			if i == 0 {
				b = append(b, "'{"...)
			}

			s := len(b)
			b = sf(fmter, b, v.Index(i))
			e := len(b)

			if s == e {
				continue
			}

			if b[s] == '\'' && b[e-1] == '\'' {
				b[s], b[e-1] = '"', '"' // schema.AppenderFunc will also append `''`
			}

			if i+1 == length {
				b = append(b, "}'"...)
				break
			}
			b = append(b, ',')
		}

		return b
	}
}

func scanArray(sf schema.ScannerFunc) func(dest reflect.Value, src any) error {
	return func(dest reflect.Value, src any) error {
		b := b2s.A(src)
		if len(b) < 2 || (b[0] != '{' && b[1] != '}') {
			setNil(dest)
			return nil
		}

		arr, err := scanLinearArray(b)
		if err != nil {
			return err
		}
		length := len(arr)
		if length == 0 {
			setNil(dest)
			return nil
		}

		slice := reflect.MakeSlice(dest.Type(), length, length)

		for i := 0; i < length; i++ {
			if err = sf(slice.Index(i), arr[i]); err != nil {
				return err
			}
		}
		dest.Set(slice)

		return nil
	}
}
