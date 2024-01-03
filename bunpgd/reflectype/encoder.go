package reflectype

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"reflect"
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
	EncoderSQL    = reflect.TypeOf((*encoderSQL)(nil)).Elem()
	EncoderJSON   = reflect.TypeOf((*encoderJSON)(nil)).Elem()
	EncoderText   = reflect.TypeOf((*encoderText)(nil)).Elem()
	EncoderBinary = reflect.TypeOf((*encoderBinary)(nil)).Elem()

	Driver  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	Scanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

	JSONUnmarshaler   = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	TextUnmarshaler   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	BinaryUnmarshaler = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
)
