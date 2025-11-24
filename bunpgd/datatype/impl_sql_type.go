package datatype

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/extra/bunjson"
	"github.com/uptrace/bun/schema"
)

func scanBytes(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetBytes(nil)
	case string:
		dest.SetBytes([]byte(src))
	case []byte:
		dest.SetBytes(bytes.Clone(src))
	default:
		return ErrNotSupportValueType
	}
	return nil
}

func scanJSON(dest reflect.Value, src any) error {
	if src == nil {
		return nil
	}

	if dest.IsNil() {
		var b []byte
		if str, isStr := src.(string); isStr {
			b = []byte(str)
		} else if v, isBts := src.([]byte); isBts {
			b = v
		} else {
			return ErrNotSupportValueType
		}

		return bunjson.Unmarshal(b, dest.Addr().Interface())
	}

	dest = dest.Elem()
	if fn := TypeScanner(dest.Type()); fn != nil {
		return fn(dest, src)
	}
	return ErrNotSupportValueType
}

func appendJSON(gen schema.QueryGen, b []byte, v reflect.Value) []byte {
	bb, err := json.Marshal(v.Interface())
	if err != nil {
		return dialect.AppendError(b, err)
	}

	if len(bb) > 0 && bb[len(bb)-1] == '\n' {
		bb = bb[:len(bb)-1]
	}

	return gen.Dialect().AppendJSON(b, bb)
}

func appendBytes(_ schema.QueryGen, b []byte, v reflect.Value) []byte {
	bs := v.Bytes()
	if bs == nil {
		return dialect.AppendNull(b)
	}

	b = append(b, `'\x`...)

	s := len(b)
	b = append(b, make([]byte, hex.EncodedLen(len(bs)))...)
	hex.Encode(b[s:], bs)

	b = append(b, '\'')

	return b
}
