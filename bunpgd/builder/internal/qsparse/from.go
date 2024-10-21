package qsparse

import (
	"encoding"
	"encoding/json"
	"errors"
	"net"
	"net/netip"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.x2ox.com/sorbifolia/bunpgd/builder/internal/query"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

func from(val []string, fd query.Field, v reflect.Value) error {
	if fd.Flag.And(reflectype.FromQS) {
		if fd.Flag.And(reflect.Pointer) && v.IsNil() {
			v.Set(reflect.New(fd.Typ))
		}
		return v.Interface().(reflectype.FromQueryParameters).FromQueryParameters(val)
	}

	if fd.Flag.Or(
		reflectype.Time, reflectype.TimeDuration,
		reflectype.IP, reflectype.IPNet,
		reflectype.HardwareAddr,
		reflectype.Addr, reflectype.Prefix) {

		if fd.Flag.And(reflect.Pointer) && v.IsNil() {
			v.Set(reflect.New(fd.Typ))
			v = reflect.Indirect(v)
		}

		return fromStruct(val, fd.Typ, v)
	}

	if fd.Flag.And(reflectype.JSONUnmarshaler) {
		if fd.Flag.And(reflect.Pointer) && v.IsNil() {
			v.Set(reflect.New(fd.Typ))
		}
		return v.Interface().(json.Unmarshaler).UnmarshalJSON([]byte(val[0]))
	}

	return fromUnknown(val, v.Type(), v)
}

func fromUnknown(val []string, rt reflect.Type, v reflect.Value) error {
	if rt.Implements(reflectype.FromQS) {
		if rt.Kind() == reflect.Pointer && v.IsNil() {
			v.Set(reflect.New(rt.Elem()))
		}
		return v.Interface().(reflectype.FromQueryParameters).FromQueryParameters(val)
	}

	switch rt {
	case reflectype.TimeDuration: // int64
		t, err := time.ParseDuration(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))
		return nil

	case reflectype.IP: // []byte
		addr, err := netip.ParseAddr(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(net.IP(addr.AsSlice())))
		return nil

	case reflectype.HardwareAddr: // []byte
		addr, err := net.ParseMAC(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(addr))
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return fromBool(val, v)
	case reflect.String:
		return fromString(val, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fromNumber(val, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fromNumber(val, v)
	case reflect.Float32, reflect.Float64:
		return fromNumber(val, v)
	case reflect.Complex64, reflect.Complex128:
		return fromNumber(val, v)

	case reflect.Array:
		return fromArray(val, rt, v)
	case reflect.Slice:
		return fromSlice(val, rt, v)
	case reflect.Pointer:
		return fromPointer(val, rt, v)
	case reflect.Struct:
		return fromStruct(val, rt, v)
	default:
	}

	return errors.New("unknown type")
}

func fromPointer(val []string, rt reflect.Type, v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(rt.Elem()))
	}

	switch rt.Elem() {
	case reflectype.Time, reflectype.IPNet, reflectype.Addr, reflectype.Prefix:
	default:
		if rt.Implements(reflectype.JSONUnmarshaler) {
			return v.Interface().(json.Unmarshaler).UnmarshalJSON([]byte(val[0]))
		}
		if rt.Implements(reflectype.TextUnmarshaler) {
			return v.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val[0]))
		}
	}

	return fromUnknown(val, rt, v.Elem())
}

func fromString(val []string, v reflect.Value) error {
	v.SetString(val[0])
	return nil
}

func fromBool(val []string, v reflect.Value) error {
	switch strings.ToLower(val[0]) {
	case "true", "t", "1":
		v.SetBool(true)
	case "false", "f", "0", "":
		v.SetBool(false)
	}
	return errors.New("invalid boolean value")
}

func fromArray(val []string, rt reflect.Type, v reflect.Value) error {
	length := len(val)
	if v.Len() < length {
		return errors.New("invalid array length")
	}
	for i := 0; i < length; i++ {
		if err := fromUnknown(val[i:], rt, v.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func fromSlice(val []string, rt reflect.Type, v reflect.Value) error {
	length := len(val)
	if v.IsNil() || v.Len() == 0 {
		v.Set(reflect.MakeSlice(v.Type(), length, length))
	}
	return fromArray(val, rt, v)
}

func fromStruct(val []string, rt reflect.Type, v reflect.Value) error {
	switch rt {
	case reflectype.TimeDuration: // int64
		t, err := time.ParseDuration(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))
		return nil

	case reflectype.IP: // []byte
		addr, err := netip.ParseAddr(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(net.IP(addr.AsSlice())))
		return nil

	case reflectype.HardwareAddr: // []byte
		addr, err := net.ParseMAC(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(addr))
		return nil

	case reflectype.Time:
		t, err := ParseTime(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(t))
		return nil

	case reflectype.IPNet:
		_, addr, err := net.ParseCIDR(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(*addr))
		return nil

	case reflectype.Addr:
		addr, err := netip.ParseAddr(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(addr))
		return nil

	case reflectype.Prefix:
		addr, err := netip.ParsePrefix(val[0])
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(addr))
		return nil
	}
	return errors.New("unknown type")
}

func fromNumber(val []string, v reflect.Value) (err error) {
	kind := v.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(val[0], 10, numberBits[kind]); err == nil {
			v.SetInt(n)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var n uint64
		if n, err = strconv.ParseUint(val[0], 10, numberBits[kind]); err == nil {
			v.SetUint(n)
		}

	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(val[0], numberBits[kind]); err == nil {
			v.SetFloat(n)
		}

	case reflect.Complex64, reflect.Complex128:
		var n complex128
		if n, err = strconv.ParseComplex(val[0], numberBits[kind]); err == nil {
			v.SetComplex(n)
		}
	default:
		panic("unknown kind")
	}
	return
}

var numberBits = []int{
	reflect.Int8: 8, reflect.Int16: 16, reflect.Int32: 32, reflect.Int64: 64,
	reflect.Uint8: 8, reflect.Uint16: 16, reflect.Uint32: 32, reflect.Uint64: 64,
	reflect.Int: 0, reflect.Uint: 0,
	reflect.Float32: 32, reflect.Float64: 64,
}

// ParseTime adapt browser format
func ParseTime(val string) (t time.Time, err error) {
	switch {
	case strings.IndexByte(val, ' ') >= 0: // toUTCString(): Mon, 02 Jan 2006 15:04:05 MST
		t, err = time.Parse(time.RFC1123, val)

	case isAllNumber(val): // getTime(): unix milli
		var num int64
		if num, err = strconv.ParseInt(val, 10, 64); err == nil {
			t = time.UnixMilli(num)
		}

	default: // toISOString(): 2006-01-02T15:04:05Z
		t, err = time.Parse(time.RFC3339Nano, val)
	}

	if err == nil {
		t = t.In(time.Local)
	}

	return
}

func isAllNumber(s string) bool {
	for i, v := range []byte(s) {
		if v >= '0' && v <= '9' {
			continue
		}
		if i == 0 && (v == '-' || v == '+') {
			continue
		}
		return false
	}
	return true
}
