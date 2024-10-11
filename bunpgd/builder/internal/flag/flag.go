package flag

import (
	"reflect"

	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

type Flag uint16

const (
	BunQueryItr Flag = 1 << iota
	JSON
	IP
	Pointer
	Slice
	String
)

func (f *Flag) Set(flag ...Flag) {
	for _, v := range flag {
		*f |= v
	}
}

func (f *Flag) Has(flag Flag) bool {
	return flag&*f != 0
}

func (f *Flag) From(rt reflect.Type) {
	*f = 0

	if rt.Implements(reflectype.QueryBuilder) {
		*f |= BunQueryItr
	}
	if rt.Implements(reflectype.JSONMarshaler) {
		*f |= JSON
	}

	kind := rt.Kind()
	if kind == reflect.Pointer {
		*f |= Pointer
		rt = rt.Elem()
		kind = rt.Kind()
	}

	switch rt {
	case reflectype.IP, reflectype.IPNet, reflectype.Addr, reflectype.Prefix:
		*f |= IP
		return
	}

	if kind == reflect.Slice {
		*f |= Slice
		rt = rt.Elem()
		kind = rt.Kind()
	}
	if kind == reflect.String {
		*f |= String
	}
}
