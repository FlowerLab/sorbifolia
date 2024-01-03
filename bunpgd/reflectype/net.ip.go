package reflectype

import (
	"net/netip"
	"reflect"
)

var (
	Addr   = reflect.TypeOf((*netip.Addr)(nil)).Elem()
	Prefix = reflect.TypeOf((*netip.Prefix)(nil)).Elem()
)
