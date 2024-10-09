package reflectype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"reflect"

	"github.com/uptrace/bun/schema"
)

type encoderSQL interface {
	sql.Scanner
	driver.Valuer
}
type encoderJSON interface {
	json.Marshaler
	json.Unmarshaler
}
type encoderText interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}
type encoderBinary interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

var (
	EncoderSQL    = reflect.TypeFor[encoderSQL]()
	EncoderJSON   = reflect.TypeFor[encoderJSON]()
	EncoderText   = reflect.TypeFor[encoderText]()
	EncoderBinary = reflect.TypeFor[encoderBinary]()

	Valuer  = reflect.TypeFor[driver.Valuer]()
	Scanner = reflect.TypeFor[sql.Scanner]()

	JSONUnmarshaler   = reflect.TypeFor[json.Unmarshaler]()
	TextUnmarshaler   = reflect.TypeFor[encoding.TextUnmarshaler]()
	BinaryUnmarshaler = reflect.TypeFor[encoding.BinaryUnmarshaler]()

	JSONMarshaler   = reflect.TypeFor[json.Marshaler]()
	TextMarshaler   = reflect.TypeFor[encoding.TextMarshaler]()
	BinaryMarshaler = reflect.TypeFor[encoding.BinaryMarshaler]()

	QueryAppender = reflect.TypeFor[schema.QueryAppender]()
)
