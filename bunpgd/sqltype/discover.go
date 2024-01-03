package sqltype

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

func SetType(field *schema.Field) error {
	field.DiscoveredSQLType, _ = field.Tag.Option("type")
	if field.DiscoveredSQLType == "" {
		field.DiscoveredSQLType = DiscoverType(field.IndirectType)
	}
	if field.DiscoveredSQLType == UnknownType {
		return fmt.Errorf("unknown type %s", field.IndirectType.String())
	}
	field.DiscoveredSQLType = strings.ToUpper(field.DiscoveredSQLType)

	if field.AutoIncrement && !field.Identity {
		switch field.DiscoveredSQLType {
		case SmallInt:
			field.CreateTableSQLType = SmallSerial
		case Integer:
			field.CreateTableSQLType = Serial
		case BigInt:
			field.CreateTableSQLType = BigSerial
		}
	}

	return nil
}

func DiscoverType(typ reflect.Type) (st string) {
	if st = DiscoverComplexType(typ); st != UnknownType {
		return
	}
	if typ.Kind() == reflect.Pointer {
		if st = DiscoverInterface(typ); st != UnknownType {
			return
		}
		return DiscoverType(typ.Elem())
	}

	switch typ.Kind() {
	case reflect.Bool:
		return Boolean
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
		return SmallInt
	case reflect.Int32, reflect.Uint32:
		return Integer
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		return BigInt
	case reflect.Float32:
		return Real
	case reflect.Float64:
		return DoublePrecision
	case reflect.String:
		return Text
	case reflect.Map:
		return JSONB
	case reflect.Slice:
		return DiscoverSliceType(typ)
	case reflect.Struct:
		return DiscoverInterface(typ)
	case reflect.Array:
		panic("TODO future support")
	default:
		return UnknownType
	}
}

func DiscoverSliceType(typ reflect.Type) string {
	if st := DiscoverType(typ.Elem()); st != UnknownType {
		return st + "[]"
	}
	return UnknownType
}

func DiscoverComplexType(typ reflect.Type) string {
	switch typ {
	case reflectype.NullTime, reflectype.Time:
		return TimestampTZ
	case reflectype.TimeDuration:
		panic("TODO future support")
	case reflectype.NullBool:
		return Boolean
	case reflectype.NullFloat:
		return DoublePrecision
	case reflectype.NullInt64:
		return BigInt
	case reflectype.NullInt32, reflectype.Rune:
		return Integer
	case reflectype.NullInt16, reflectype.Byte:
		return SmallInt
	case reflectype.NullString:
		return Text
	case reflectype.JSONRawMessage:
		return JSONB
	case reflectype.Addr, reflectype.IP:
		return INet
	case reflectype.Prefix, reflectype.IPNet:
		return CIDR
	case reflectype.HardwareAddr:
		return MacAddr
	default:
		return UnknownType
	}
}

func DiscoverInterface(typ reflect.Type) string {
	switch {
	case typ.Implements(reflectype.EncoderJSON):
		return JSONB
	case typ.Implements(reflectype.EncoderText):
		return Text
	case typ.Implements(reflectype.EncoderBinary):
		return Bytea
	}
	return UnknownType
}
