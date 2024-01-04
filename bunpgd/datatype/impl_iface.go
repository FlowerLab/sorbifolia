package datatype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

var (
	ifSQLDriver = &Adapter{Append: ifValuer, Scan: ifScanner}
	ifJSON      = &Adapter{Append: ifJSONMarshaler, Scan: ifJSONUnmarshaler}
	ifText      = &Adapter{Append: ifTextMarshaler, Scan: ifJSONTextUnmarshaler}
	ifBinary    = &Adapter{Append: ifBinaryMarshaler, Scan: ifBinaryUnmarshaler}
)

func ifScanner(dest reflect.Value, src any) error { return dest.Interface().(sql.Scanner).Scan(src) }

func ifValuer(_ schema.Formatter, b []byte, v reflect.Value) []byte {
	val, _ := v.Interface().(driver.Valuer).Value()
	switch val := val.(type) {
	case []byte:
		return append(b, val...)
	case string:
		return append(b, b2s.S(val)...)
	}
	return nil
}

func ifJSONUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(json.Unmarshaler).UnmarshalJSON(b2s.A(src))
}
func ifJSONTextUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(encoding.TextUnmarshaler).UnmarshalText(b2s.A(src))
}
func ifBinaryUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b2s.A(src))
}

func ifJSONMarshaler(_ schema.Formatter, b []byte, v reflect.Value) []byte {
	bts, _ := v.Interface().(json.Marshaler).MarshalJSON()
	return append(b, bts...)
}

func ifTextMarshaler(_ schema.Formatter, b []byte, v reflect.Value) []byte {
	bts, _ := v.Interface().(encoding.TextMarshaler).MarshalText()
	return append(b, bts...)
}

func ifBinaryMarshaler(_ schema.Formatter, b []byte, v reflect.Value) []byte {
	bts, _ := v.Interface().(encoding.BinaryMarshaler).MarshalBinary()
	return append(b, bts...)
}
