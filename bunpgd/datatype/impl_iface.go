package datatype

import (
	"database/sql"
	"encoding"
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

func ifScanner(dest reflect.Value, src any) error { return dest.Interface().(sql.Scanner).Scan(src) }

func ifJSONUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(json.Unmarshaler).UnmarshalJSON(b2s.A(src))
}
func ifJSONTextUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(encoding.TextUnmarshaler).UnmarshalText(b2s.A(src))
}

func ifJSONMarshaler(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
	bts, _ := v.Interface().(json.Marshaler).MarshalJSON()
	return fmter.Dialect().AppendString(b, b2s.B(bts))
}

func ifTextMarshaler(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
	bts, _ := v.Interface().(encoding.TextMarshaler).MarshalText()
	return fmter.Dialect().AppendString(b, b2s.B(bts))
}
