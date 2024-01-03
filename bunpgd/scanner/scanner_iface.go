package scanner

import (
	"database/sql"
	"encoding"
	"encoding/json"
	"fmt"
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
func ifBinaryUnmarshaler(dest reflect.Value, src any) error {
	return dest.Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary(b2s.A(src))
}

func addrScanner(fn schema.ScannerFunc) schema.ScannerFunc {
	return func(dest reflect.Value, src any) error {
		if !dest.CanAddr() {
			return fmt.Errorf("bunpgd: Scan(nonaddressable %T)", dest.Interface())
		}
		if err := fn(dest.Addr(), src); err != nil {
			return err
		}

		if dest.Elem().IsZero() {
			dest.SetZero()
		}
		return nil
	}
}
