package datatype

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/sqltype"
)

var sqlTypeAdapters = map[string]*Adapter{
	sqltype.Bytea: {Append: appendBytes, Scan: scanBytes},
	sqltype.JSON:  {Append: appendJSON, Scan: scanBytes},
	sqltype.JSONB: {Append: appendJSON, Scan: scanBytes},

	sqltype.HSTORE: {Append: appendHstore, Scan: scanHstore},
}

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

func appendJSON(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
	bb, err := json.Marshal(v.Interface())
	if err != nil {
		return dialect.AppendError(b, err)
	}

	if len(bb) > 0 && bb[len(bb)-1] == '\n' {
		bb = bb[:len(bb)-1]
	}

	return fmter.Dialect().AppendJSON(b, bb)
}

func appendBytes(_ schema.Formatter, b []byte, v reflect.Value) []byte {
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
