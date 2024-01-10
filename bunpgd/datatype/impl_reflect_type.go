package datatype

import (
	"net"
	"net/netip"
	"reflect"

	"go.x2ox.com/sorbifolia/bunpgd/internal/b2s"
)

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
		return ErrNotSupportValueType
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
