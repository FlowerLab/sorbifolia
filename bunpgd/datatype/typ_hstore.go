package datatype

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

func scanHstore(dest reflect.Value, src any) error {
	dest = reflect.Indirect(dest)
	if !dest.CanSet() {
		return fmt.Errorf("bunpgd: Scan(non-settable %s)", dest.Type())
	}

	var (
		m        = make(map[string]string)
		b        byte
		pair     = [][]byte{{}, {}}
		pi       = 0
		inQuote  = false
		didQuote = false
		sawSlash = false
		idx      = 0
	)

	for idx, b = range src.([]byte) {
		if sawSlash {
			pair[pi] = append(pair[pi], b)
			sawSlash = false
			continue
		}

		switch b {
		case '\\':
			sawSlash = true
			continue
		case '"':
			if inQuote = !inQuote; !didQuote {
				didQuote = true
			}
			continue
		default:
			if !inQuote {
				switch b {
				case ' ', '\t', '\n', '\r':
					continue
				case '=':
					continue
				case '>':
					pi = 1
					didQuote = false
					continue
				case ',':
					s := string(pair[1])
					if didQuote || len(s) != 4 || strings.ToLower(s) != "null" {
						m[string(pair[0])] = string(pair[1])
					}
					pair[0] = []byte{}
					pair[1] = []byte{}
					pi = 0
					continue
				}
			}
		}
		pair[pi] = append(pair[pi], b)
	}
	if idx > 0 {
		s := string(pair[1])
		if didQuote || len(s) != 4 || strings.ToLower(s) != "null" {
			m[string(pair[0])] = string(pair[1])
		}
	}

	dest.Set(reflect.ValueOf(m))

	return nil
}

func appendHstore(_ schema.QueryGen, b []byte, v reflect.Value) []byte {
	m := v.Convert(reflectype.MapStringString).Interface().(map[string]string)
	if m == nil {
		return append(b, "NULL"...)
	}

	buf := &bytes.Buffer{}
	for key, val := range m {
		buf.WriteString(hQuote(key))
		buf.Write([]byte{'=', '>'})
		buf.WriteString(hQuote(val))
		buf.WriteByte(',')
	}
	bs := buf.Bytes()
	if len(bs) > 1 {
		bs = bs[:len(bs)-1]
	}
	b = bytes.Clone(bs)

	return append(b, bytes.Clone(bs)...)
}
