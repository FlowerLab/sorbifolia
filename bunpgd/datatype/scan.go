package datatype

import (
	"reflect"
	"sync"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

var typeScannerMap sync.Map

func TypeScanner(rt reflect.Type) (sf schema.ScannerFunc) {
	if val, ok := typeScannerMap.Load(rt); ok {
		return val.(schema.ScannerFunc)
	}

	defer func() {
		if sf != nil {
			typeScannerMap.LoadOrStore(rt, sf)
		}
	}()

	kind := rt.Kind()
	if kind == reflect.Ptr {
		if sf = TypeScanner(rt.Elem()); sf != nil {
			return schema.PtrScanner(sf)
		}
	}

	switch {
	case rt.Implements(reflectype.Scanner):
		return ifScanner
	case rt.Implements(reflectype.TextUnmarshaler):
		return ifTextUnmarshaler
	case rt.Implements(reflectype.JSONUnmarshaler):
		return ifJSONUnmarshaler
	case rt == reflectype.IPNet:
		return scanHardwareAddr
	case rt == reflectype.HardwareAddr:
		return scanINetIP
	}

	if kind != reflect.Ptr {
		typ := reflect.PointerTo(rt)
		switch {
		case typ.Implements(reflectype.Scanner):
			return addrScanner(ifScanner)
		case typ.Implements(reflectype.TextUnmarshaler):
			return addrScanner(ifTextUnmarshaler)
		case typ.Implements(reflectype.JSONUnmarshaler):
			return addrScanner(ifJSONUnmarshaler)
		}
	}

	if kind == reflect.Slice {
		if sf = TypeScanner(rt.Elem()); sf != nil {
			return scanArray(sf)
		}
		return scanArray(schema.Scanner(rt.Elem()))
	}

	return schema.Scanner(rt)
}
