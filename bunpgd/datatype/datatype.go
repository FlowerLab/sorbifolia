package datatype

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/sqltype"
)

func Set(field *schema.Field) {
	switch field.DiscoveredSQLType {
	case sqltype.Bytea:
		field.Append, field.Scan = appendBytes, scanBytes
	case sqltype.JSON:
		field.Append, field.Scan = appendJSON, scanJSON
	case sqltype.JSONB:
		field.Append, field.Scan = appendJSON, scanJSON
	case sqltype.HSTORE:
		field.Append, field.Scan = appendHstore, scanHstore
	default:
		field.Append, field.Scan = TypeAppender(field.StructField.Type), TypeScanner(field.StructField.Type)
	}
}

func addrScanner(fn schema.ScannerFunc) schema.ScannerFunc {
	return func(dest reflect.Value, src any) error {
		if !dest.CanAddr() {
			return fmt.Errorf("bunpgd: Scan(nonaddressable %T)", dest.Interface())
		}
		if err := fn(dest.Addr(), src); err != nil {
			return err
		}

		if dest.IsZero() {
			dest.SetZero()
		}
		return nil
	}
}

func addrAppender(fn schema.AppenderFunc) schema.AppenderFunc {
	return func(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
		if !v.CanAddr() {
			err := fmt.Errorf("bunpgd: Append(nonaddressable %T)", v.Interface())
			return dialect.AppendError(b, err)
		}
		return fn(fmter, b, v.Addr())
	}
}

var ErrNotSupportValueType = errors.New("not support value type")
