package datatype

import (
	"reflect"
	"sync"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

var typeAppenderMap sync.Map

func TypeAppender(rt reflect.Type) (sf schema.AppenderFunc) {
	if val, ok := typeAppenderMap.Load(rt); ok {
		return val.(schema.AppenderFunc)
	}

	defer func() {
		if sf != nil {
			typeAppenderMap.LoadOrStore(rt, sf)
		}
	}()

	switch {
	case rt.Implements(reflectype.QueryAppender):
		return ifQueryAppender
	case rt.Implements(reflectype.Valuer):
		return schema.Appender(nil, rt)
	case rt.Implements(reflectype.TextMarshaler):
		return ifTextMarshaler
	case rt.Implements(reflectype.JSONMarshaler):
		return ifJSONMarshaler
	case rt == reflectype.IPNet:
		return appendHardwareAddr
	case rt == reflectype.HardwareAddr:
		return appendINetIP
	}

	kind := rt.Kind()
	if kind == reflect.Ptr {
		if sf = TypeAppender(rt.Elem()); sf != nil {
			return schema.PtrAppender(sf)
		}
	}

	if kind != reflect.Ptr {
		typ := reflect.PointerTo(rt)
		switch {
		case rt.Implements(reflectype.QueryAppender):
			return addrAppender(ifQueryAppender)
		case rt.Implements(reflectype.Valuer):
			return schema.Appender(nil, rt)
		case typ.Implements(reflectype.TextMarshaler):
			return addrAppender(ifTextMarshaler)
		case typ.Implements(reflectype.JSONMarshaler):
			return addrAppender(ifJSONMarshaler)
		}
	}

	if kind == reflect.Slice {
		if sf = TypeAppender(rt.Elem()); sf != nil {
			return appendArray(sf)
		}
	}

	return schema.Appender(nil, rt)
}
