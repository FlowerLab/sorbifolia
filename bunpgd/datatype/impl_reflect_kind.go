package datatype

import (
	"reflect"
	"strconv"
	"time"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

var kindAdapters = []*Adapter{
	reflect.Bool: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte { return dialect.AppendNull(b) },
		Scan:   scanBool, IsZero: nil,
	},
	reflect.Int: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendInt(b, v.Int(), 10)
		},
		Scan: scanInt64, IsZero: nil,
	},
	reflect.Int8: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendInt(b, v.Int(), 10)
		},
		Scan: scanInt64, IsZero: nil,
	},
	reflect.Int16: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendInt(b, v.Int(), 10)
		},
		Scan: scanInt64, IsZero: nil,
	},
	reflect.Int32: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendInt(b, v.Int(), 10)
		},
		Scan: scanInt64, IsZero: nil,
	},
	reflect.Int64: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendInt(b, v.Int(), 10)
		},
		Scan: scanInt64, IsZero: nil,
	},
	reflect.Uint: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendUint(b, v.Uint(), 10)
		},
		Scan: scanUint64, IsZero: nil,
	},
	reflect.Uint8: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendUint(b, v.Uint(), 10)
		},
		Scan: scanUint64, IsZero: nil,
	},
	reflect.Uint16: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendUint(b, v.Uint(), 10)
		},
		Scan: scanUint64, IsZero: nil,
	},
	reflect.Uint32: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendUint(b, v.Uint(), 10)
		},
		Scan: scanUint64, IsZero: nil,
	},
	reflect.Uint64: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return strconv.AppendUint(b, v.Uint(), 10)
		},
		Scan: scanUint64, IsZero: nil,
	},
	reflect.Float32: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return dialect.AppendFloat32(b, float32(v.Float()))
		},
		Scan: scanFloat64, IsZero: nil,
	},
	reflect.Float64: {
		Append: func(_ schema.Formatter, b []byte, v reflect.Value) []byte {
			return dialect.AppendFloat64(b, v.Float())
		},
		Scan: scanFloat64, IsZero: nil,
	},
	reflect.String: {
		Append: func(fmter schema.Formatter, b []byte, v reflect.Value) []byte {
			return fmter.Dialect().AppendString(b, v.String())
		},
		Scan: scanString, IsZero: nil,
	},
	// reflect.Pointer: {
	//
	// },

	reflect.Array:         nil,
	reflect.Slice:         nil,
	reflect.UnsafePointer: nil,
}

func scanBool(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetBool(false)
	case bool:
		dest.SetBool(src)
	case int64:
		dest.SetBool(src != 0)
	case []byte:
		f, err := strconv.ParseBool(b2s.B(src))
		if err != nil {
			return err
		}
		dest.SetBool(f)
	case string:
		f, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		dest.SetBool(f)
	default:
		return ErrNotSupportValueType
	}
	return nil
}

func scanInt64(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetInt(0)
	case int64:
		dest.SetInt(src)
	case uint64:
		dest.SetInt(int64(src))
	case []byte:
		n, err := strconv.ParseInt(b2s.B(src), 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(n)
	case string:
		n, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(n)
	default:
		return ErrNotSupportValueType
	}
	return nil
}

func scanUint64(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetUint(0)
	case uint64:
		dest.SetUint(src)
	case int64:
		dest.SetUint(uint64(src))
	case []byte:
		n, err := strconv.ParseUint(b2s.B(src), 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
	case string:
		n, err := strconv.ParseUint(src, 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
	default:
		return ErrNotSupportValueType
	}
	return nil
}

func scanFloat64(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetFloat(0)
	case float64:
		dest.SetFloat(src)
	case []byte:
		f, err := strconv.ParseFloat(b2s.B(src), 64)
		if err != nil {
			return err
		}
		dest.SetFloat(f)
	case string:
		f, err := strconv.ParseFloat(src, 64)
		if err != nil {
			return err
		}
		dest.SetFloat(f)
	default:
		return ErrNotSupportValueType
	}
	return nil
}

func scanString(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetString("")
	case string:
		dest.SetString(src)
	case []byte:
		dest.SetString(string(src))
	case time.Time:
		dest.SetString(src.Format(time.RFC3339Nano))
	case int64:
		dest.SetString(strconv.FormatInt(src, 10))
	case uint64:
		dest.SetString(strconv.FormatUint(src, 10))
	case float64:
		dest.SetString(strconv.FormatFloat(src, 'G', -1, 64))
	default:
		return ErrNotSupportValueType
	}
	return nil
}
