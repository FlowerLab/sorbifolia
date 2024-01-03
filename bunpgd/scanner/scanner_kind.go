package scanner

import (
	"reflect"
	"strconv"
	"time"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

var kindScannerFunc = []schema.ScannerFunc{
	reflect.Bool:    scanBool,
	reflect.Int:     scanInt64,
	reflect.Int8:    scanInt64,
	reflect.Int16:   scanInt64,
	reflect.Int32:   scanInt64,
	reflect.Int64:   scanInt64,
	reflect.Uint:    scanUint64,
	reflect.Uint8:   scanUint64,
	reflect.Uint16:  scanUint64,
	reflect.Uint32:  scanUint64,
	reflect.Uint64:  scanUint64,
	reflect.Uintptr: scanUint64,
	reflect.Float32: scanFloat64,
	reflect.Float64: scanFloat64,
	reflect.String:  scanString,

	reflect.Map: scanMap,

	reflect.Array: nil,
	reflect.Slice: nil,
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
