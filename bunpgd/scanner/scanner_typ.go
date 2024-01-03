package scanner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"reflect"
	"sync"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

var scannerFuncMap sync.Map

func RegisterScanner(rt reflect.Type, sf schema.ScannerFunc) {
	if _, loaded := scannerFuncMap.LoadOrStore(rt, sf); loaded {
		panic(fmt.Sprintf("type %s is loaded", rt))
	}
}

func scanBytes(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		dest.SetBytes(nil)
	case string:
		dest.SetBytes([]byte(src))
	case []byte:
		dest.SetBytes(bytes.Clone(src))
	default:
		return errors.New("not support")
	}
	return nil
}

func scanHardwareAddr(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		setNil(dest)

	case []byte:
		hw, err := net.ParseMAC(b2s.B(src))
		if err != nil {
			return err
		}
		dest.Set(reflect.ValueOf(hw))

	case string:
		hw, err := net.ParseMAC(src)
		if err != nil {
			return err
		}
		dest.Set(reflect.ValueOf(hw))

	default:
		return errors.New("not support")
	}

	return nil
}

func scanINetIP(dest reflect.Value, src any) error {
	switch src := src.(type) {
	case nil:
		setNil(dest)
		return nil

	case []byte:
		if len(src) == 0 {
			setNil(dest)
			return nil
		}
		dest.Set(reflect.ValueOf(net.ParseIP(b2s.B(src))))

	case string:
		if len(src) == 0 {
			setNil(dest)
			return nil
		}
		dest.Set(reflect.ValueOf(net.ParseIP(src)))

	case net.IP:
		if len(src) == 0 {
			setNil(dest)
			return nil
		}
		dest.Set(reflect.ValueOf(src))

	case netip.Addr:
		if !src.IsValid() {
			setNil(dest)
			return nil
		}
		dest.Set(reflect.ValueOf(src.AsSlice()))
	}

	return ErrNotSupportValueType
}

func scanMap(dest reflect.Value, src any) error {
	if src == nil {
		setNil(dest)
		return nil
	}
	return json.Unmarshal(b2s.A(src), dest.Addr().Interface())
}

func setNil(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if rv.IsNil() {
			return
		}
	}
	rv.Set(reflect.New(rv.Type()).Elem())
}
