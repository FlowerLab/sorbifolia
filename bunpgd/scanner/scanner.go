package scanner

import (
	"reflect"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
	"go.x2ox.com/sorbifolia/bunpgd/sqltype"
)

func SetScanner(field *schema.Field) error {
	switch field.DiscoveredSQLType {
	case sqltype.Bytea:
		field.Scan = scanBytes
	case sqltype.JSON, sqltype.JSONB:
		field.Scan = scanMap
	case sqltype.HSTORE:
		panic("not support")
	default:
		field.Scan = Scanner(field.StructField.Type)
	}

	if field.Scan == nil {
		return ErrUnknownTypeCannotMatchScanner
	}
	return nil
}

func Scanner(rt reflect.Type) schema.ScannerFunc {
	if val, ok := scannerFuncMap.Load(rt); ok {
		return val.(schema.ScannerFunc)
	}
	if sf := scanner(rt); sf != nil {
		scannerFuncMap.LoadOrStore(rt, sf)
		return sf
	}
	return nil
}

func scanner(rt reflect.Type) schema.ScannerFunc {
	kind := rt.Kind()
	if kind == reflect.Ptr {
		if fn := Scanner(rt.Elem()); fn != nil {
			return schema.PtrScanner(fn)
		}
	}

	switch {
	case rt.Implements(reflectype.Scanner):
		return ifScanner
	case rt.Implements(reflectype.TextUnmarshaler):
		return ifJSONTextUnmarshaler
	case rt.Implements(reflectype.JSONUnmarshaler):
		return ifJSONUnmarshaler
	case rt.Implements(reflectype.BinaryUnmarshaler):
		return ifBinaryUnmarshaler
	case rt == reflectype.IPNet:
		return scanINetIP
	case rt == reflectype.HardwareAddr:
		return scanHardwareAddr
	}

	if kind != reflect.Ptr {
		typ := reflect.PtrTo(rt)
		switch {
		case typ.Implements(reflectype.Scanner):
			return addrScanner(ifScanner)
		case typ.Implements(reflectype.TextUnmarshaler):
			return addrScanner(ifJSONTextUnmarshaler)
		case typ.Implements(reflectype.JSONUnmarshaler):
			return addrScanner(ifJSONUnmarshaler)
		case typ.Implements(reflectype.BinaryUnmarshaler):
			return addrScanner(ifBinaryUnmarshaler)
		}
	}

	return kindScannerFunc[kind]
}
