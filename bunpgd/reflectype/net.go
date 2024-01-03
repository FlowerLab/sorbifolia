package reflectype

import (
	"net"
	"reflect"
)

var (
	IP           = reflect.TypeOf((*net.IP)(nil)).Elem()
	IPNet        = reflect.TypeOf((*net.IPNet)(nil)).Elem()
	HardwareAddr = reflect.TypeOf((*net.HardwareAddr)(nil)).Elem()
)
