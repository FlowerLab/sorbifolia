package reflectype

import (
	"net"
	"net/netip"
	"reflect"
)

var (
	IP           = reflect.TypeFor[net.IP]()
	IPNet        = reflect.TypeFor[net.IPNet]()
	HardwareAddr = reflect.TypeFor[net.HardwareAddr]()

	Addr   = reflect.TypeFor[netip.Addr]()
	Prefix = reflect.TypeFor[netip.Prefix]()
)
