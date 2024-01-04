package datatype

import (
	"reflect"
	"sync"

	"github.com/uptrace/bun/schema"
	"go.x2ox.com/sorbifolia/bunpgd/reflectype"
)

var adapterMap sync.Map

func Set(field *schema.Field) {
	if adp := sqlTypeAdapters[field.DiscoveredSQLType]; adp != nil {
		adp.Set(field)
		return
	}
	if val, ok := adapterMap.Load(field.StructField.Type); ok {
		val.(*Adapter).Set(field)
		return
	}

	if adp := FindDatatype(field.StructField.Type); adp != nil {
		adapterMap.LoadOrStore(field.StructField.Type, adp)
		adp.Set(field)
		return
	}
}

func FindDatatype(rt reflect.Type) *Adapter {
	if adp := FindWithImplement(rt); adp != nil {
		return adp
	}
	if adp := typeAdapters[rt]; adp != nil {
		return adp
	}

	kind := rt.Kind()
	if kind == reflect.Ptr {
		if adp := FindDatatype(rt.Elem()); adp != nil {
			return adp.Ptr()
		}
	}

	if kind != reflect.Ptr {
		if adp := FindWithImplement(reflect.PtrTo(rt)); adp != nil {
			return adp.Addr()
		}
	}

	if kind == reflect.Slice {
		if adp := FindWithSlice(rt); adp != nil {
			return adp
		}
	}

	return kindAdapters[rt.Kind()]
}

func FindWithImplement(rt reflect.Type) *Adapter {
	switch {
	case rt.Implements(reflectype.EncoderSQL):
		return ifSQLDriver
	case rt.Implements(reflectype.TextUnmarshaler):
		return ifText
	case rt.Implements(reflectype.JSONUnmarshaler):
		return ifJSON
	case rt.Implements(reflectype.BinaryUnmarshaler):
		return ifBinary
	}

	return nil
}

func FindWithSlice(rt reflect.Type) *Adapter {
	adp := FindDatatype(rt.Elem())
	if adp == nil {
		return nil // not support type
	}
	return arrayAdapter(adp)
}
